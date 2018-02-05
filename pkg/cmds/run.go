package cmds

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/signals"
	"github.com/appscode/kubed/pkg/operator"
	"github.com/appscode/kutil/meta"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/clientcmd"
)

// runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml"
func NewCmdRun() *cobra.Command {
	opt := operator.Options{
		ConfigPath:        "/srv/kubed/config.yaml",
		APIAddress:        ":8080",
		WebAddress:        ":56790",
		ScratchDir:        "/tmp",
		OperatorNamespace: meta.Namespace(),
		ResyncPeriod:      10 * time.Minute,
		// ref: https://github.com/kubernetes/ingress-nginx/blob/e4d53786e771cc6bdd55f180674b79f5b692e552/pkg/ingress/controller/launch.go#L252-L259
		// High enough QPS to fit all expected use cases. QPS=0 is not set here, because client code is overriding it.
		QPS: 1e6,
		// High enough Burst to fit all expected use cases. Burst=0 is not set here, because client code is overriding it.
		Burst:              1e6,
		PrometheusCrdGroup: prom.Group,
		PrometheusCrdKinds: prom.DefaultCrdKinds,
	}
	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run daemon",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			log.Infoln("Starting kubed...")

			Run(opt)
		},
	}

	fs := flag.NewFlagSet("prometheus", flag.ExitOnError)
	fs.StringVar(&opt.PrometheusCrdGroup, "prometheus-crd-apigroup", opt.PrometheusCrdGroup, "prometheus CRD  API group name")
	fs.Var(&opt.PrometheusCrdKinds, "prometheus-crd-kinds", " - EXPERIMENTAL (could be removed in future releases) - customize CRD kind names")
	cmd.Flags().AddGoFlagSet(fs)

	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&opt.ConfigPath, "clusterconfig", opt.ConfigPath, "Path to cluster config file")
	cmd.Flags().StringVar(&opt.ScratchDir, "scratch-dir", opt.ScratchDir, "Directory used to store temporary files. Use an `emptyDir` in Kubernetes.")
	cmd.Flags().StringVar(&opt.APIAddress, "api.address", opt.APIAddress, "The address of the Kubed API Server (overrides any value in clusterconfig)")
	cmd.Flags().StringVar(&opt.WebAddress, "web.address", opt.WebAddress, "Address to listen on for web interface and telemetry.")
	cmd.Flags().DurationVar(&opt.ResyncPeriod, "resync-period", opt.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	cmd.Flags().Float32Var(&opt.QPS, "qps", opt.QPS, "The maximum QPS to the master from this client")
	cmd.Flags().IntVar(&opt.Burst, "burst", opt.Burst, "The maximum burst for throttle")

	return cmd
}

func Run(opt operator.Options) {
	log.Infoln("configurations provided for kubed", opt)
	defer runtime.HandleCrash()

	config, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.KubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	config.Burst = opt.Burst
	config.QPS = opt.QPS

	op, err := operator.New(config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	err = op.Configure()
	if err != nil {
		log.Fatalln(err)
	}

	stopCh := signals.SetupSignalHandler()

	log.Infoln("Running kubed watcher")
	op.RunAndHold(stopCh)
}
