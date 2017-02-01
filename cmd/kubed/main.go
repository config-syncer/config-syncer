package main

import (
	_ "appscode/pkg/errorhandlers"

	"appscode.com/kubed/pkg"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/hold"
	v "github.com/appscode/go/version"
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"github.com/spf13/pflag"
)

var (
	Version         string
	VersionStrategy string
	Os              string
	Arch            string
	CommitHash      string
	GitBranch       string
	GitTag          string
	CommitTimestamp string
	BuildTimestamp  string
	BuildHost       string
	BuildHostOs     string
	BuildHostArch   string
)

func init() {
	v.Version.Version = Version
	v.Version.VersionStrategy = VersionStrategy
	v.Version.Os = Os
	v.Version.Arch = Arch
	v.Version.CommitHash = CommitHash
	v.Version.GitBranch = GitBranch
	v.Version.GitTag = GitTag
	v.Version.CommitTimestamp = CommitTimestamp
	v.Version.BuildTimestamp = BuildTimestamp
	v.Version.BuildHost = BuildHost
	v.Version.BuildHostOs = BuildHostOs
	v.Version.BuildHostArch = BuildHostArch
}

func main() {
	config := &pkg.Config{}
	pflag.StringVarP(&config.APITokenPath, "api-token", "t", "/var/run/secrets/appscode/api-token", "Endpoint of elasticsearch")
	pflag.StringVar(&config.Master, "master", "", "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	pflag.StringVar(&config.KubeConfig, "kubeconfig", "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	pflag.StringVarP(&config.APIEndpoint, "api-endpoint", "e", "api.appscode.com:50077", "appscode api server host:port")
	pflag.StringVarP(&config.ProviderName, "cloud-provider", "c", "", "Name of cloud provider")
	pflag.StringVarP(&config.ClusterName, "cluster-name", "k", "", "Name of Kubernetes cluster")
	pflag.StringVarP(&config.LoadbalancerImageName, "haproxy-image", "h", "appscode/haproxy:1.7.0-k8s", "haproxy image name to be run")
	pflag.StringVarP(&config.ESEndpoint, "es-endpoint", "", "", "Endpoint of elasticsearch")
	pflag.StringVarP(&config.InfluxSecretName, "influx-secret", "", "", "Influxdb secret name")
	pflag.StringVarP(&config.IcingaSecretName, "icinga-secret", "", "", "Icinga secret name")
	pflag.StringVarP(&config.IngressClass, "ingress-class", "I", "", "Ingress class name to use with")

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
