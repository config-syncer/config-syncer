package system

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_env "github.com/appscode/go/env"
	"github.com/appscode/log"
)

var publicMatrix = map[string]string{
	"aws-ap-northeast-1":       "https://s3-ap-northeast-1.amazonaws.com/appscode-tokyo/binaries",
	"aws-ap-northeast-2":       "https://s3-ap-northeast-2.amazonaws.com/appscode-seoul/binaries",
	"aws-ap-south-1":           "https://s3-ap-south-1.amazonaws.com/appscode-mumbai/binaries",
	"aws-ap-southeast-1":       "https://s3-ap-southeast-1.amazonaws.com/appscode-singapore/binaries",
	"aws-ap-southeast-2":       "https://s3-ap-southeast-2.amazonaws.com/appscode-sydney/binaries",
	"aws-ca-central-1":         "https://s3-ca-central-1.amazonaws.com/appscode-montreal/binaries",
	"aws-eu-central-1":         "https://s3-eu-central-1.amazonaws.com/appscode-frankfurt/binaries",
	"aws-eu-west-1":            "https://s3-eu-west-1.amazonaws.com/appscode-ireland/binaries",
	"aws-eu-west-2":            "https://s3-eu-west-2.amazonaws.com/appscode-london/binaries",
	"aws-sa-east-1":            "https://s3-sa-east-1.amazonaws.com/appscode-saopaulo/binaries",
	"aws-us-east-1":            "https://s3.amazonaws.com/appscode-virginia/binaries",
	"aws-us-east-2":            "https://s3-us-east-2.amazonaws.com/appscode-ohio/binaries",
	"aws-us-west-1":            "https://s3-us-west-1.amazonaws.com/appscode-norcal/binaries",
	"aws-us-west-2":            "https://s3-us-west-2.amazonaws.com/appscode-oregon/binaries",
	"azure-eastus":             "https://storage.googleapis.com/appscode-us/binaries",
	"azure-eastus2":            "https://storage.googleapis.com/appscode-us/binaries",
	"azure-centralus":          "https://storage.googleapis.com/appscode-us/binaries",
	"azure-northcentralus":     "https://storage.googleapis.com/appscode-us/binaries",
	"azure-southcentralus":     "https://storage.googleapis.com/appscode-us/binaries",
	"azure-westcentralus":      "https://storage.googleapis.com/appscode-us/binaries",
	"azure-westus":             "https://storage.googleapis.com/appscode-us/binaries",
	"azure-westus2":            "https://storage.googleapis.com/appscode-us/binaries",
	"azure-canadaeast":         "https://storage.googleapis.com/appscode-us/binaries",
	"azure-canadacentral":      "https://storage.googleapis.com/appscode-us/binaries",
	"azure-brazilsouth":        "https://storage.googleapis.com/appscode-us/binaries",
	"azure-northeurope":        "https://storage.googleapis.com/appscode-eu/binaries",
	"azure-westeurope":         "https://storage.googleapis.com/appscode-eu/binaries",
	"azure-ukwest":             "https://storage.googleapis.com/appscode-eu/binaries",
	"azure-uksouth":            "https://storage.googleapis.com/appscode-eu/binaries",
	"azure-southeastasia":      "https://storage.googleapis.com/appscode-asia/binaries",
	"azure-eastasia":           "https://storage.googleapis.com/appscode-asia/binaries",
	"azure-australiaeast":      "https://storage.googleapis.com/appscode-asia/binaries",
	"azure-australiasoutheast": "https://storage.googleapis.com/appscode-asia/binaries",
	"azure-japaneast":          "https://storage.googleapis.com/appscode-asia/binaries",
	"azure-japanwest":          "https://storage.googleapis.com/appscode-asia/binaries",
	"digitalocean-ams2":        "https://storage.googleapis.com/appscode-eu/binaries",
	"digitalocean-ams3":        "https://storage.googleapis.com/appscode-eu/binaries",
	"digitalocean-blr1":        "https://storage.googleapis.com/appscode-asia/binaries",
	"digitalocean-fra1":        "https://storage.googleapis.com/appscode-eu/binaries",
	"digitalocean-lon1":        "https://storage.googleapis.com/appscode-eu/binaries",
	"digitalocean-nyc1":        "https://storage.googleapis.com/appscode-us/binaries",
	"digitalocean-nyc2":        "https://storage.googleapis.com/appscode-us/binaries",
	"digitalocean-nyc3":        "https://storage.googleapis.com/appscode-us/binaries",
	"digitalocean-sfo1":        "https://storage.googleapis.com/appscode-us/binaries",
	"digitalocean-sfo2":        "https://storage.googleapis.com/appscode-us/binaries",
	"digitalocean-sgp1":        "https://storage.googleapis.com/appscode-asia/binaries",
	"digitalocean-tor1":        "https://storage.googleapis.com/appscode-us/binaries",
	"gce-asia":                 "https://storage.googleapis.com/appscode-asia/binaries",
	"gce-eu":                   "https://storage.googleapis.com/appscode-eu/binaries",
	"gce-us":                   "https://storage.googleapis.com/appscode-us/binaries",
	"hetzner":                  "https://storage.googleapis.com/appscode-eu/binaries",
	"linode-2":                 "https://storage.googleapis.com/appscode-us/binaries",   // Dallas, TX, USA
	"linode-3":                 "https://storage.googleapis.com/appscode-us/binaries",   // Fremont, CA, USA
	"linode-4":                 "https://storage.googleapis.com/appscode-us/binaries",   // Atlanta, GA, USA
	"linode-6":                 "https://storage.googleapis.com/appscode-us/binaries",   // Newark, NJ, USA
	"linode-7":                 "https://storage.googleapis.com/appscode-eu/binaries",   // London, England, UK
	"linode-8":                 "https://storage.googleapis.com/appscode-asia/binaries", // Tokyo, JP
	"linode-9":                 "https://storage.googleapis.com/appscode-asia/binaries", // Singapore, SG
	"linode-10":                "https://storage.googleapis.com/appscode-eu/binaries",   // Frankfurt, DE
	"ovh-bhs1":                 "https://storage.googleapis.com/appscode-us/binaries",
	"ovh-gra1":                 "https://storage.googleapis.com/appscode-eu/binaries",
	"ovh-sbg1":                 "https://storage.googleapis.com/appscode-eu/binaries",
	"packet-ams1":              "https://storage.googleapis.com/appscode-eu/binaries",
	"packet-ewr1":              "https://storage.googleapis.com/appscode-us/binaries",
	"packet-sjc1":              "https://storage.googleapis.com/appscode-us/binaries",
	"rackspace-dfw":            "https://storage.googleapis.com/appscode-us/binaries",
	"rackspace-iad":            "https://storage.googleapis.com/appscode-us/binaries",
	"rackspace-ord":            "https://storage.googleapis.com/appscode-us/binaries",
	"rackspace-hkg":            "https://storage.googleapis.com/appscode-asia/binaries",
	"rackspace-syd":            "https://storage.googleapis.com/appscode-asia/binaries",
	"scaleway-ams1":            "https://storage.googleapis.com/appscode-eu/binaries",
	"scaleway-par1":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-ams":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-che":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-dal":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-fra":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-hkg":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-hou":            "https://storage.googleapis.com/appscode-usa/binaries",
	"softlayer-lon":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-mel":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-mex":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-mil":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-mon":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-osl":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-par":            "https://storage.googleapis.com/appscode-eu/binaries",
	"softlayer-sjc":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-sao":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-seo":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-sea":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-sng":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-syd":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-tok":            "https://storage.googleapis.com/appscode-asia/binaries",
	"softlayer-tor":            "https://storage.googleapis.com/appscode-us/binaries",
	"softlayer-wdc":            "https://storage.googleapis.com/appscode-us/binaries",
	"vultr-1":                  "https://storage.googleapis.com/appscode-us/binaries",   // New Jersey
	"vultr-2":                  "https://storage.googleapis.com/appscode-us/binaries",   // Chicago
	"vultr-3":                  "https://storage.googleapis.com/appscode-us/binaries",   // Dallas
	"vultr-4":                  "https://storage.googleapis.com/appscode-us/binaries",   // Seattle
	"vultr-5":                  "https://storage.googleapis.com/appscode-us/binaries",   // Los Angeles
	"vultr-6":                  "https://storage.googleapis.com/appscode-us/binaries",   // Atlanta
	"vultr-7":                  "https://storage.googleapis.com/appscode-eu/binaries",   // Amsterdam
	"vultr-8":                  "https://storage.googleapis.com/appscode-eu/binaries",   // London
	"vultr-9":                  "https://storage.googleapis.com/appscode-eu/binaries",   // Frankfurt
	"vultr-12":                 "https://storage.googleapis.com/appscode-us/binaries",   // Silicon Valley
	"vultr-19":                 "https://storage.googleapis.com/appscode-asia/binaries", // Sydney
	"vultr-24":                 "https://storage.googleapis.com/appscode-eu/binaries",   // Paris
	"vultr-25":                 "https://storage.googleapis.com/appscode-asia/binaries", // Tokyo
	"vultr-39":                 "https://storage.googleapis.com/appscode-us/binaries",   // Miami
	"vultr-40":                 "https://storage.googleapis.com/appscode-asia/binaries", // Singapore
}

var devMatrix = map[string]string{
	"aws": "https://s3.amazonaws.com/appscode-dev/binaries",
	"gce": "https://storage.googleapis.com/appscode-dev/binaries",
}

const cdnPrefix = "https://cdn.appscode.com/binaries"

const (
	AppKubeServer     = "kubernetes-server"
	AppKubeStarter    = "start-kubernetes"
	AppKubeSaltbase   = "kubernetes-salt"
	AppOneboxApps     = "onebox-apps"
	AppOneboxSaltbase = "onebox-salt"
	AppCIStarter      = "start-ci"
	AppCISaltbase     = "ci-salt"
	AppHostfacts      = "hostfacts"
)

func NewAppStartCI(provider, region, version string) *Application {
	return &Application{
		Name: AppCIStarter,
		URL:  cloudURL(provider, region, AppCIStarter, version, "start-ci-linux-amd64"),
	}
}

func NewAppCISalt(provider, region, version string) *Application {
	return &Application{
		Name: AppCISaltbase,
		URL:  cloudURL(provider, region, AppCISaltbase, version, "ci-salt.tar.gz"),
	}
}

func NewAppHostfacts(provider, region, version string) *Application {
	return &Application{
		Name: AppHostfacts,
		URL:  cdnURL(provider, region, AppHostfacts, version, "hostfacts-linux-amd64"),
	}
}

func NewAppStartKubernetes(provider, region, version string) *Application {
	return &Application{
		Name: AppKubeStarter,
		URL:  cloudURL(provider, region, AppKubeStarter, version, "start-kubernetes-linux-amd64"),
	}
}

func NewAppKubernetesSalt(provider, region, version string) *Application {
	return &Application{
		Name: AppKubeSaltbase,
		URL:  cloudURL(provider, region, AppKubeSaltbase, version, "kubernetes-salt.tar.gz"),
	}
}

func NewAppOneboxApps(provider, region, version string) *Application {
	return &Application{
		Name: AppOneboxApps,
		URL:  cdnURL(provider, region, AppOneboxApps, version, "onebox-apps.tar.gz"),
	}
}

func NewAppOneboxSalt(provider, region, version string) *Application {
	return &Application{
		Name: AppOneboxSaltbase,
		URL:  cdnURL(provider, region, AppOneboxSaltbase, version, "onebox-salt.tar.gz"),
	}
}

func NewAppKubernetesServer(provider, region, version string) *Application {
	return &Application{
		Name: AppKubeServer,
		URL:  cloudURL(provider, region, AppKubeServer, version, "kubernetes-server-linux-amd64.tar.gz"),
	}
}

func cloudURL(provider, region, name, version, file string) string {
	var prefix string
	if _env.FromHost().IsPublic() {
		matrix := publicMatrix
		if provider == "aws" {
			prefix = matrix[provider+"-"+region]
		} else if provider == "azure" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else if provider == "gce" {
			prefix = matrix[provider+"-"+region[0:strings.Index(region, "-")]]
		} else if provider == "digitalocean" {
			prefix = matrix[provider+"-"+region]
		} else if provider == "linode" {
			prefix = matrix[provider+"-"+region]
		} else if provider == "vultr" {
			prefix = matrix[provider+"-"+region]
		} else if provider == "scaleway" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else if provider == "softlayer" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else if provider == "packet" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else if provider == "hetzner" {
			prefix = matrix[provider]
		} else if provider == "rackspace" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else if provider == "ovh" {
			prefix = matrix[provider+"-"+strings.ToLower(region)]
		} else {
			prefix = matrix["gce-us"]
		}
	} else {
		matrix := devMatrix
		if provider == "aws" {
			prefix = matrix["aws"]
		} else {
			prefix = matrix["gce"]
		}
	}
	return fmt.Sprintf("%s/%s/%s/%s", prefix, name, version, file)
}

func cdnURL(provider, region, name, version, file string) string {
	var prefix string
	if _env.FromHost().IsPublic() {
		prefix = cdnPrefix
	} else {
		prefix = devMatrix["gce"]
	}
	return fmt.Sprintf("%s/%s/%s/%s", prefix, name, version, file)
}

type Application struct {
	Name     string `json:"NAME"`
	URL      string `json:"URL"`
	Checksum string `json:"CHECKSUM"`
}

func (app *Application) String() string {
	return app.URL
}

func (app *Application) Download(dir string) {
	path := dir + "/" + app.Name
	if _, err := os.Stat(path); err == nil && checksum(path) == app.Checksum {
		return
	}
	fmt.Println("Downloading", app.URL, " to", path)
	o, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()
	log.Infof("Downloading %v app from %v", app.Name, app.URL)
	res, err := http.Get(app.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	n, err := io.Copy(o, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n, "bytes downloaded.")
}

func checksum(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
