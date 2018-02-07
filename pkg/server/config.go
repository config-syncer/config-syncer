package server

import (
	"path/filepath"
	"time"

	"github.com/appscode/kubed/pkg/eventer"
	rbin "github.com/appscode/kubed/pkg/recyclebin"
	resource_indexer "github.com/appscode/kubed/pkg/registry/resource"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kutil/tools/fsnotify"
	srch_cs "github.com/appscode/searchlight/client"
	searchlightinformers "github.com/appscode/searchlight/informers/externalversions"
	scs "github.com/appscode/stash/client"
	stashinformers "github.com/appscode/stash/informers/externalversions"
	vcs "github.com/appscode/voyager/client"
	voyagerinformers "github.com/appscode/voyager/informers/externalversions"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/kubedb/apimachinery/client"
	kubedbinformers "github.com/kubedb/apimachinery/informers/externalversions"
	"github.com/robfig/cron"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type OperatorConfig struct {
	ClientConfig *rest.Config

	OpsAddress string
	ScratchDir string
	ConfigPath string

	KubeClient        kubernetes.Interface
	VoyagerClient     vcs.Interface
	SearchlightClient srch_cs.Interface
	StashClient       scs.Interface
	KubeDBClient      kcs.Interface
	PromClient        prom.MonitoringV1Interface
}

func NewOperatorConfig(clientConfig *rest.Config) *OperatorConfig {
	return &OperatorConfig{
		ClientConfig: clientConfig,
	}
}

func (c *OperatorConfig) New() (*Operator, error) {
	op := &Operator{
		KubeClient:        c.KubeClient,
		VoyagerClient:     c.VoyagerClient,
		SearchlightClient: c.SearchlightClient,
		StashClient:       c.StashClient,
		KubeDBClient:      c.KubeDBClient,
		PromClient:        c.PromClient,
	}

	op.recorder = eventer.NewEventRecorder(op.KubeClient, "kubed")
	op.trashCan = &rbin.RecycleBin{}
	op.eventProcessor = &eventer.EventForwarder{Client: op.KubeClient.Discovery()}
	op.configSyncer = syncer.New(op.KubeClient, op.recorder)

	op.cron = cron.New()
	op.cron.Start()

	// Enable full text indexing to have search feature
	indexDir := filepath.Join(c.ScratchDir, "indices")
	op.searchIndexer = resource_indexer.NewIndexer(indexDir)

	op.watcher = &fsnotify.Watcher{
		WatchDir: filepath.Dir(c.ConfigPath),
		Reload:   op.Configure,
	}

	// ---------------------------
	op.kubeInformerFactory = informers.NewSharedInformerFactory(op.KubeClient, 10*time.Minute)
	op.voyagerInformerFactory = voyagerinformers.NewSharedInformerFactory(op.VoyagerClient, 10*time.Minute)
	op.stashInformerFactory = stashinformers.NewSharedInformerFactory(op.StashClient, 10*time.Minute)
	op.searchlightInformerFactory = searchlightinformers.NewSharedInformerFactory(op.SearchlightClient, 10*time.Minute)
	op.kubedbInformerFactory = kubedbinformers.NewSharedInformerFactory(op.KubeDBClient, 10*time.Minute)
	// ---------------------------
	op.setupWorkloadInformers()
	op.setupNetworkInformers()
	op.setupConfigInformers()
	op.setupRBACInformers()
	op.setupNodeInformers()
	op.setupEventInformers()
	op.setupCertificateInformers()
	// ---------------------------
	op.setupVoyagerInformers()
	op.setupStashInformers()
	op.setupSearchlightInformers()
	op.setupKubeDBInformers()
	op.setupPrometheusInformers()
	// ---------------------------

	if err := op.Configure(); err != nil {
		return nil, err
	}
	return op, nil
}
