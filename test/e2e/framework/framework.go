package framework

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"
	kcs "github.com/appscode/kubed/client/clientset/versioned"
	sls "github.com/appscode/searchlight/client/clientset/versioned"
	srch_cs "github.com/appscode/searchlight/client/clientset/versioned"
	scs "github.com/appscode/stash/client/clientset/versioned"
	vcs "github.com/appscode/voyager/client/clientset/versioned"
	prom "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	pcs "github.com/coreos/prometheus-operator/pkg/client/versioned"
	kdbcs "github.com/kubedb/apimachinery/client/clientset/versioned"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	ecs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ka "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"kmodules.xyz/client-go/tools/certstore"
)

const (
	MaxRetry = 200
	NoRetry  = 1

	DefaultEventuallyTimeout         = 5 * time.Minute
	DefaultEventuallyPollingInterval = 2 * time.Second
)

type Framework struct {
	KubeClient         clientset.Interface
	KubedClient        kcs.Interface
	KAClient           ka.Interface
	VoyagerClient      vcs.Interface
	SearchlightClient  srch_cs.Interface
	StashClient        scs.Interface
	KubeDBClient       kdbcs.Interface
	PromClient         pcs.Interface
	crdClient          ecs.ApiextensionsV1beta1Interface
	namespace          string
	Mutex              sync.Mutex
	CertStore          *certstore.CertStore
	KubeConfigPath     string
	SelfHostedOperator bool
	ClientConfig       *rest.Config
}

func New(config *rest.Config) *Framework {
	promClient, err := pcs.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	store, err := certstore.NewCertStore(afero.NewMemMapFs(), filepath.Join("", "pki"))
	Expect(err).NotTo(HaveOccurred())

	err = store.InitCA()
	Expect(err).NotTo(HaveOccurred())

	return &Framework{
		namespace: rand.WithUniqSuffix("test-kubed"),

		ClientConfig:      config,
		KubeClient:        clientset.NewForConfigOrDie(config),
		KAClient:          ka.NewForConfigOrDie(config),
		KubedClient:       kcs.NewForConfigOrDie(config),
		crdClient:         ecs.NewForConfigOrDie(config),
		CertStore:         store,
		StashClient:       scs.NewForConfigOrDie(config),
		VoyagerClient:     vcs.NewForConfigOrDie(config),
		SearchlightClient: sls.NewForConfigOrDie(config),
		KubeDBClient:      kdbcs.NewForConfigOrDie(config),
		PromClient:        promClient,
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("kubed-e2e"),
	}
}

func (f *Invocation) App() string {
	return f.app
}

type Invocation struct {
	*Framework
	app string
}

func (f *Framework) EnsureCreatedCRDs() error {
	_, pErr := f.PromClient.MonitoringV1().Prometheuses(f.namespace).List(metav1.ListOptions{})
	_, sErr := f.PromClient.MonitoringV1().ServiceMonitors(f.namespace).List(metav1.ListOptions{})
	if pErr == nil && sErr == nil {
		return nil
	}
	promCrd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.PrometheusName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.PrometheusName,
				Kind:   prom.PrometheusesKind,
			},
		},
	}
	_, err := f.crdClient.CustomResourceDefinitions().Create(promCrd)
	if err != nil {
		return err
	}
	svcMonitorCrd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.ServiceMonitorName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.ServiceMonitorName,
				Kind:   prom.ServiceMonitorsKind,
			},
		},
	}
	_, err = f.crdClient.CustomResourceDefinitions().Create(svcMonitorCrd)
	if err != nil {
		return err
	}

	alertMgr := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.AlertmanagerName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.AlertmanagerName,
				Kind:   prom.AlertmanagersKind,
			},
		},
	}
	_, err = f.crdClient.CustomResourceDefinitions().Create(alertMgr)
	if err != nil {
		return err
	}
	return nil
}
