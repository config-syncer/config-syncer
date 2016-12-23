package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/appscode/go/encoding/yaml"

	"github.com/appscode/go-dns/aws"
	"github.com/appscode/go-dns/azure"
	"github.com/appscode/go-dns/cloudflare"
	"github.com/appscode/go-dns/digitalocean"
	"github.com/appscode/go-dns/googlecloud"
	"github.com/appscode/go-dns/linode"
	"github.com/appscode/go-dns/vultr"
	_env "github.com/appscode/go/env"
	"github.com/appscode/log"
)

type URLBase struct {
	Scheme     string `json:"scheme,omitempty"`
	BaseDomain string `json:"base_domain,omitempty"`
}

type config struct {
	Network struct {
		PublicUrls       URLBase `json:"public_urls,omitempty"`
		TeamUrls         URLBase `json:"team_urls,omitempty"`
		ClusterUrls      URLBase `json:"cluster_urls,omitempty"`
		InClusterUrls    URLBase `json:"in_cluster_urls,omitempty"`
		URLShortenerUrls URLBase `json:"URL_shortener_urls,omitempty"`
		FileUrls         URLBase `json:"file_urls,omitempty"`
	} `json:"network,omitempty"`
	SkipStartupConfigAPI bool `json:"skip_startup_config_api,omitempty"`
	Phabricator          struct {
		Glusterfs struct {
			Endpoint string `json:"endpoint,omitempty"`
			Path     string `json:"path,omitempty"`
		} `json:"glusterfs"`
		DaemonImage            string `json:"daemon_image,omitempty"`
		PhabricatorDataProject string `json:"phabricator_data_project, omitempty"`
	} `json:"phabricator, omitempty"`
	Artifactory struct {
		ElasticSearchEndpoint string `json:"elasticsearch_endpoint,omitempty"`
	} `json:"artifactory,omitempty"`
	DefaultAppVersion struct {
		KubeServer   string `json:"kube_server, omitempty"`
		KubeStarter  string `json:"kube_starter, omitempty"`
		KubeSaltbase string `json:"kube_saltbase, omitempty"`
		CIStarter    string `json:"ci_starter, omitempty"`
		Jenkins      string `json:"jenkins, omitempty"`
		CISaltbase   string `json:"ci_saltbase, omitempty"`
		Hostfacts    string `json:"hostfacts, omitempty"`
	} `json:"default_app_version, omitempty"`
	Jenkins struct {
		Master struct {
			Image           string `json:"image,omitempty"`
			SSHCredentialID string `json:"ssh_credential_id,omitempty"`
			UserName        string `json:"username, omitempty"`
		}
		Proxy struct {
			Image string `json:"image,omitempty"`
		}
		//grsecSnapShot  = 15920953 // grsec enabled image id in DO.
	} `json:"jenkins, omitempty"`
	GlusterFs struct {
		ServerImage     string `json:"glusterd_image,omitempty"`
		ControllerImage string `json:"glusterc_image,omitempty"`
	} `json:"glusterfs, omitempty"`
	Bacula struct {
		ServerImage   string `json:"server_image,omitempty"`
		ClientImage   string `json:"client_image,omitempty"`
		DatabaseImage string `json:"database_image,omitempty"`
	} `json:"bacula,omitempty"`
	HAProxy struct {
		Image string `json:"image,omitempty"`
	} `json:"haproxy,omitempty"`
	GoogleAnalytics struct {
		PublicTracker string `json:"public_tracker, omitempty"`
		TeamTracker   string `json:"team_tracker, omitempty"`
	} `json:"google_analytics, omitempty"`
	LetsEncrypt struct {
		CADirURL string `json:"ca_dir_url, omitempty"`
		KeyType  string `json:"key_type, omitempty"`
	} `json:"lets_encrypt, omitempty"`
	GCE struct {
		AppscodeCustomerKeys string `json:"appscode_customer_keys, omitempty"`
		AppscodeCIData       string `json:"appscode_ci_data, omitempty"`
	} `json:"gce, omitempty"`
}

type secureConfig struct {
	Secret          string `json:"secret,omitempty"`
	MagicCodeSecret string `json:"magic_code_secret,omitempty"`
	Database        struct {
		MetaNamespace string   `json:"meta_ns"`
		MetaHost      string   `json:"meta_host,omitempty"`
		Port          int      `json:"port,omitempty"`
		User          string   `json:"user,omitempty"`
		Password      string   `json:"password,omitempty"`
		Hosts         []string `json:"hosts,omitempty"`
	} `json:"database,omitempty"`
	Twilio struct {
		Token       string `json:"token,omitempty"`
		AccountSid  string `json:"account_sid,omitempty"`
		PhoneNumber string `json:"phone_number,omitempty"`
	} `json:"twilio,omitempty"`
	Mailgun struct {
		ApiKey       string `json:"key,omitempty"`
		PublicDomain string `json:"public_domain,omitempty"`
		TeamDomain   string `json:"team_domain,omitempty"`
	} `json:"mailgun"`
	DNS struct {
		// Deprecated
		Credential     map[string]string `json:"credential,omitempty"`
		CredentialFile string            `json:"credential_file,omitempty"`

		// Generic DNS Providers
		Provider     string               `json:"provider,omitempty"`
		AWS          aws.Options          `json:"aws,omitempty"`
		Azure        azure.Options        `json:"azure,omitempty"`
		Cloudflare   cloudflare.Options   `json:"cloudflare,omitempty"`
		Digitalocean digitalocean.Options `json:"digitalocean,omitempty"`
		Gcloud       googlecloud.Options  `json:"gcloud,omitempty"`
		Linode       linode.Options       `json:"linode,omitempty"`
		Vultr        vultr.Options        `json:"vultr,omitempty"`
	} `json:"dns"`
	DigitalOcean struct {
		Token string `json:"token"`
	} `json:"digitalocean"`
	GCE struct {
		CredentialFile string `json:"credential_file,omitempty"`
	} `json:"gce"`
	S3 struct {
		AccessKey string `json:"access_key,omitempty"`
		SecretKey string `json:"secret_key,omitempty"`
		Region    string `json:"regopm,omitempty"`
		Endpoint  string `json:"endpoint,omitempty"`
	} `json:"s3"`
	Braintree struct {
		MerchantID string `json:"merchant_id, omitempty"`
		PublicKey  string `json:"public_key, omitempty"`
		PrivateKey string `json:"private_key, omitempty"`
	} `json:"braintree, omitempty"`
	Icinga2 struct {
		Host        string `json:"host"`
		APIUser     string `json:"api_user"`
		APIPassword string `json:"api_password"`
	} `json:"icinga2, omitempty"`
	MetricsSink struct {
		InfluxDB struct {
			URL      string `json:"url, omitempty"`
			Database string `json:"database, omitempty"`
			Username string `json:"username, omitempty"`
			Password string `json:"password, omitempty"`
		} `json:"influxdb, omitempty"`
		Statsd struct {
			URL string `json:"url, omitempty"`
		} `json:"statsd, omitempty"`
	} `json:"metrics_sink, omitempty"`
}

var Secrets secureConfig
var Config config
var cOnce sync.Once

func Init() {
	cOnce.Do(func() {
		fmt.Println("[system] Reading system config file")
		env := _env.FromHost()
		candidates := []string{
			"/srv/appscode/secrets." + env.String() + ".yaml", // ~/.appscode/secrets.env.yaml
			"/srv/appscode/secrets." + env.String() + ".yml",  // ~/.appscode/secrets.env.yml
			"/srv/appscode/secrets." + env.String() + ".json", // ~/.appscode/secrets.env.json
		}
		for _, cfgFile := range candidates {
			fmt.Printf("Searching %v\n", cfgFile)
			if _, err := os.Stat(cfgFile); err == nil {
				data, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log.Fatal(err)
				}
				jsonData, err := yaml.ToJSON(data)
				if err != nil {
					log.Fatal(err)
				}

				err = json.Unmarshal(jsonData, &Secrets)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][][][][]")
				fmt.Printf("Using system secret file %v\n", cfgFile)
				fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][][][][]")
				break
			}
		}
		fmt.Println("******************************************************")
		fmt.Println("[system] Reading config file")
		configFiles := []string{
			"/srv/appscode/config." + env.String() + ".yaml", // ~/.appscode/config.env.yaml
			"/srv/appscode/config." + env.String() + ".yml",  // ~/.appscode/config.env.yml
			"/srv/appscode/config." + env.String() + ".json", // ~/.appscode/config.env.json
		}
		for _, cfgFile := range configFiles {
			fmt.Printf("Searching %v\n", cfgFile)
			if _, err := os.Stat(cfgFile); err == nil {
				data, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					log.Fatal(err)
				}

				jsonData, err := yaml.ToJSON(data)
				if err != nil {
					log.Fatal(err)
				}

				err = json.Unmarshal(jsonData, &Config)
				if err != nil {
					log.Fatal(err)
				}
				applyDefaults(env)

				fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][][][][]")
				fmt.Printf("Using system configuration file %v\n", cfgFile)
				fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][][][][]")
				fmt.Println("******************************************************")
				return
			}
		}

		log.Fatalln("Missing system configuration file.")
	})
}

func applyDefaults(env _env.Environment) {
	if env == _env.Prod {
		Config.Network.PublicUrls.BaseDomain = "appscode.com"
		Config.Network.TeamUrls.BaseDomain = "appscode.io"
		Config.Network.ClusterUrls.BaseDomain = "appscode.net"
		Config.Network.URLShortenerUrls.BaseDomain = "appsco.de"
		Config.Network.FileUrls.BaseDomain = "appscode.space"
	} else if env == _env.QA || env == _env.Dev {
		Config.Network.PublicUrls.BaseDomain = "appscode.info"
		Config.Network.TeamUrls.BaseDomain = "appscode.ninja"
		Config.Network.ClusterUrls.BaseDomain = "appscode.xyz"
		Config.Network.URLShortenerUrls.BaseDomain = "appscode.co"
		Config.Network.FileUrls.BaseDomain = "appscode.org"
	} else if env == _env.Onebox || env == _env.BoxDev {
		// TeamUrls.BaseDomain must be provided by user
		// Config.Network.TeamUrls.BaseDomain
		Config.Network.PublicUrls.BaseDomain = Config.Network.TeamUrls.BaseDomain
		Config.Network.ClusterUrls.BaseDomain = "kubernetes." + Config.Network.TeamUrls.BaseDomain
		Config.Network.URLShortenerUrls.BaseDomain = "x." + Config.Network.TeamUrls.BaseDomain
		Config.Network.FileUrls.BaseDomain = Config.Network.TeamUrls.BaseDomain
	}
	Config.Network.InClusterUrls.Scheme = "http"
	Config.Network.InClusterUrls.BaseDomain = "default"
}
