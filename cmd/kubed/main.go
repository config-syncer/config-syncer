package main

import (
	_ "appscode/pkg/errorhandlers"

	"appscode.com/kubed/pkg"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/hold"
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"github.com/spf13/pflag"
)

func main() {
	config := &pkg.Config{
		APITokenPath:          "/var/run/secrets/appscode/api-token",
		APIEndpoint:           "api.appscode.com:50077",
		LoadbalancerImageName: "appscode/haproxy:1.7.0-1.5.0",
		InfluxSecretName:      "appscode-influx",
		InfluxSecretNamespace: "kube-system",
		IcingaSecretName:      "appscode-icinga",
		IcingaSecretNamespace: "kube-system",
		EnablePromMonitoring:  false,
	}
	pflag.StringVar(&config.APITokenPath, "api-token", config.APITokenPath, "Endpoint of elasticsearch")
	pflag.StringVar(&config.Master, "master", config.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	pflag.StringVar(&config.KubeConfig, "kubeconfig", config.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	pflag.StringVar(&config.APIEndpoint, "api-endpoint", config.APIEndpoint, "appscode api server host:port")
	pflag.StringVar(&config.ProviderName, "cloud-provider", config.ProviderName, "Name of cloud provider")
	pflag.StringVar(&config.ClusterName, "cluster-name", config.ClusterName, "Name of Kubernetes cluster")
	pflag.StringVar(&config.LoadbalancerImageName, "haproxy-image", config.LoadbalancerImageName, "haproxy image name to be run")
	pflag.StringVar(&config.ESEndpoint, "es-endpoint", config.ESEndpoint, "Endpoint of elasticsearch")
	pflag.StringVar(&config.InfluxSecretName, "influx-secret", config.InfluxSecretName, "Influxdb secret name")
	pflag.StringVar(&config.InfluxSecretNamespace, "influx-secret-namespace", config.InfluxSecretNamespace, "Influxdb secret namespace")
	pflag.StringVar(&config.IcingaSecretName, "icinga-secret", config.IcingaSecretName, "Icinga secret name")
	pflag.StringVar(&config.IcingaSecretNamespace, "icinga-secret-namespace", config.IcingaSecretNamespace, "Icinga secret namespace")
	pflag.StringVar(&config.IngressClass, "ingress-class", config.IngressClass, "Ingress class name to use with")
	pflag.BoolVar(&config.EnablePromMonitoring, "enable-prometheus-monitoring", config.EnablePromMonitoring, "Enable Prometheus monitoring")

	flags.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if config.APIEndpoint == "" ||
		config.ProviderName == "" ||
		config.ClusterName == "" ||
		config.LoadbalancerImageName == "" ||
		config.APITokenPath == "" {
		log.Fatalln("required flag not provided.")
	}

	log.Infoln("Starting Kubed Process...")
	go pkg.Run(config)

	hold.Hold()
}
