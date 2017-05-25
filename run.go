package main

import (
	"github.com/appscode/go/hold"
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/kubed/pkg"
	"github.com/appscode/log"
	"github.com/spf13/cobra"
)

func NewCmdRun() *cobra.Command {
	opt := pkg.Options{
		APITokenPath:          "/var/run/secrets/appscode/api-token",
		APIEndpoint:           "api.appscode.com:50077",
		InfluxSecretName:      "appscode-influx",
		InfluxSecretNamespace: "kube-system",
		EnablePromMonitoring:  false,
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if opt.APIEndpoint == "" ||
				opt.ClusterName == "" ||
				opt.APITokenPath == "" {
				log.Fatalln("required flag not provided.")
			}
			log.Infoln("Starting Kubed Process...")
			go pkg.Run(opt)

			hold.Hold()
		},
	}

	cmd.Flags().StringVar(&opt.APITokenPath, "api-token", opt.APITokenPath, "Endpoint of elasticsearch")
	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.APIEndpoint, "api-endpoint", opt.APIEndpoint, "appscode api server host:port")
	cmd.Flags().StringVar(&opt.ClusterName, "cluster-name", opt.ClusterName, "Name of Kubernetes cluster")
	cmd.Flags().StringVar(&opt.ESEndpoint, "es-endpoint", opt.ESEndpoint, "Endpoint of elasticsearch")
	cmd.Flags().StringVar(&opt.InfluxSecretName, "influx-secret", opt.InfluxSecretName, "Influxdb secret name")
	cmd.Flags().StringVar(&opt.InfluxSecretNamespace, "influx-secret-namespace", opt.InfluxSecretNamespace, "Influxdb secret namespace")
	cmd.Flags().BoolVar(&opt.EnablePromMonitoring, "enable-prometheus-monitoring", opt.EnablePromMonitoring, "Enable Prometheus monitoring")

	return cmd
}
