package cmds

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/operator"
	srch_cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	scs "github.com/appscode/stash/client/typed/stash/v1alpha1"
	vcs "github.com/appscode/voyager/client/typed/voyager/v1beta1"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	prometheusCrdGroup = prom.Group
	prometheusCrdKinds = prom.DefaultCrdKinds
)

// runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml"
func NewCmdRun() *cobra.Command {
	opt := operator.Options{
		ConfigPath:        "/srv/kubed/config.yaml",
		APIAddress:        ":8080",
		WebAddress:        ":56790",
		ScratchDir:        "/tmp",
		OperatorNamespace: namespace(),
		ResyncPeriod:      5 * time.Minute,
		// ref: https://github.com/kubernetes/ingress-nginx/blob/e4d53786e771cc6bdd55f180674b79f5b692e552/pkg/ingress/controller/launch.go#L252-L259
		// High enough QPS to fit all expected use cases. QPS=0 is not set here, because client code is overriding it.
		QPS: 1e6,
		// High enough Burst to fit all expected use cases. Burst=0 is not set here, because client code is overriding it.
		Burst: 1e6,
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
	fs.StringVar(&prometheusCrdGroup, "prometheus-crd-apigroup", prometheusCrdGroup, "prometheus CRD  API group name")
	fs.Var(&prometheusCrdKinds, "prometheus-crd-kinds", " - EXPERIMENTAL (could be removed in future releases) - customize CRD kind names")
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

	op := &operator.Operator{
		KubeClient:        kubernetes.NewForConfigOrDie(config),
		VoyagerClient:     vcs.NewForConfigOrDie(config),
		SearchlightClient: srch_cs.NewForConfigOrDie(config),
		StashClient:       scs.NewForConfigOrDie(config),
		KubeDBClient:      kcs.NewForConfigOrDie(config),
		Opt:               opt,
	}
	op.PromClient, err = prom.NewForConfig(&prometheusCrdKinds, prometheusCrdGroup, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = op.Setup()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Running kubed watcher")
	op.RunAndHold()
}

func namespace() string {
	if ns := os.Getenv("OPERATOR_NAMESPACE"); ns != "" {
		return ns
	}
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}
	return core.NamespaceDefault
}
