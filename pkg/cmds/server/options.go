package server

import (
	"flag"

	"github.com/appscode/kubed/pkg/server"
	srch_cs "github.com/appscode/searchlight/client"
	scs "github.com/appscode/stash/client"
	vcs "github.com/appscode/voyager/client"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/kubedb/apimachinery/client"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type OperatorOptions struct {
	ConfigPath          string
	OpsAddress          string
	ScratchDir          string
	ControllerNamespace string

	QPS   float32
	Burst int

	PrometheusCrdGroup string
	PrometheusCrdKinds prom.CrdKinds
}

func NewOperatorOptions() *OperatorOptions {
	return &OperatorOptions{
		ConfigPath: "/srv/kubed/config.yaml",
		OpsAddress: ":56790",
		ScratchDir: "/tmp",
		// ref: https://github.com/kubernetes/ingress-nginx/blob/e4d53786e771cc6bdd55f180674b79f5b692e552/pkg/ingress/controller/launch.go#L252-L259
		// High enough QPS to fit all expected use cases. QPS=0 is not set here, because client code is overriding it.
		QPS: 1e6,
		// High enough Burst to fit all expected use cases. Burst=0 is not set here, because client code is overriding it.
		Burst:              1e6,
		PrometheusCrdGroup: prom.Group,
		PrometheusCrdKinds: prom.DefaultCrdKinds,
	}
}

func (s *OperatorOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("prometheus", flag.ExitOnError)
	pfs.StringVar(&s.PrometheusCrdGroup, "prometheus-crd-apigroup", s.PrometheusCrdGroup, "prometheus CRD  API group name")
	pfs.Var(&s.PrometheusCrdKinds, "prometheus-crd-kinds", " - EXPERIMENTAL (could be removed in future releases) - customize CRD kind names")
	fs.AddGoFlagSet(pfs)

	fs.StringVar(&s.ConfigPath, "clusterconfig", s.ConfigPath, "Path to cluster config file")
	fs.StringVar(&s.ScratchDir, "scratch-dir", s.ScratchDir, "Directory used to store temporary files. Use an `emptyDir` in Kubernetes.")
	fs.StringVar(&s.OpsAddress, "ops-address", s.OpsAddress, "Address to listen on for web interface and telemetry.")
	fs.Float32Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
}

func (s *OperatorOptions) ApplyTo(config *server.OperatorConfig) error {
	var err error

	config.ClientConfig.QPS = s.QPS
	config.ClientConfig.Burst = s.Burst

	if config.KubeClient, err = kubernetes.NewForConfig(config.ClientConfig); err != nil {
		return err
	}
	if config.VoyagerClient, err = vcs.NewForConfig(config.ClientConfig); err != nil {
		return err
	}
	if config.SearchlightClient, err = srch_cs.NewForConfig(config.ClientConfig); err != nil {
		return err
	}
	if config.StashClient, err = scs.NewForConfig(config.ClientConfig); err != nil {
		return err
	}
	if config.KubeDBClient, err = kcs.NewForConfig(config.ClientConfig); err != nil {
		return err
	}
	if config.PromClient, err = prom.NewForConfig(&s.PrometheusCrdKinds, s.PrometheusCrdGroup, config.ClientConfig); err != nil {
		return err
	}

	config.OpsAddress = s.OpsAddress
	config.ScratchDir = s.ScratchDir
	config.ConfigPath = s.ConfigPath

	return nil
}
