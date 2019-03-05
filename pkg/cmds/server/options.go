package server

import (
	"time"

	"github.com/appscode/kubed/pkg/operator"
	srch_cs "github.com/appscode/searchlight/client/clientset/versioned"
	scs "github.com/appscode/stash/client/clientset/versioned"
	vcs "github.com/appscode/voyager/client/clientset/versioned"
	pcm "github.com/coreos/prometheus-operator/pkg/client/versioned"
	kcs "github.com/kubedb/apimachinery/client/clientset/versioned"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/client-go/meta"
)

type OperatorOptions struct {
	ConfigPath string
	ScratchDir string

	QPS          float32
	Burst        int
	ResyncPeriod time.Duration
}

func NewOperatorOptions() *OperatorOptions {
	return &OperatorOptions{
		ConfigPath: "/srv/kubed/config.yaml",
		ScratchDir: "/tmp",
		// ref: https://github.com/kubernetes/ingress-nginx/blob/e4d53786e771cc6bdd55f180674b79f5b692e552/pkg/ingress/controller/launch.go#L252-L259
		// High enough QPS to fit all expected use cases. QPS=0 is not set here, because client code is overriding it.
		QPS: 1e6,
		// High enough Burst to fit all expected use cases. Burst=0 is not set here, because client code is overriding it.
		Burst:        1e6,
		ResyncPeriod: 10 * time.Minute,
	}
}

func (s *OperatorOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ConfigPath, "clusterconfig", s.ConfigPath, "Path to cluster config file")
	fs.StringVar(&s.ScratchDir, "scratch-dir", s.ScratchDir, "Directory used to store temporary files. Use an `emptyDir` in Kubernetes.")

	fs.Float32Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
	fs.DurationVar(&s.ResyncPeriod, "resync-period", s.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
}

func (s *OperatorOptions) ApplyTo(cfg *operator.OperatorConfig) error {
	var err error

	cfg.OperatorNamespace = meta.Namespace()
	cfg.ClientConfig.QPS = s.QPS
	cfg.ClientConfig.Burst = s.Burst
	cfg.ResyncPeriod = s.ResyncPeriod
	cfg.Test = false

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.VoyagerClient, err = vcs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.SearchlightClient, err = srch_cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.StashClient, err = scs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.KubeDBClient, err = kcs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.PromClient, err = pcm.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}

	cfg.ScratchDir = s.ScratchDir
	cfg.ConfigPath = s.ConfigPath

	return nil
}
