package storage

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	tapi "github.com/appscode/kubed/pkg/config"
	otx "github.com/appscode/osm/context"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	gcs "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	SecretMountPath = "/etc/osm"
)

func WriteOSMConfig(client clientset.Interface, snapshot tapi.Backend, namespace string, path string) error {
	osmCtx, err := NewOSMContext(client, snapshot, namespace)
	if err != nil {
		return err
	}
	osmCfg := &otx.OSMConfig{
		CurrentContext: osmCtx.Name,
		Contexts:       []*otx.Context{osmCtx},
	}
	osmBytes, err := yaml.Marshal(osmCfg)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, osmBytes, 0644)
}

func CheckBucketAccess(client clientset.Interface, spec tapi.Backend, namespace string) error {
	cfg, err := NewOSMContext(client, spec, namespace)
	if err != nil {
		return err
	}
	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	c, err := spec.Container()
	if err != nil {
		return err
	}
	container, err := loc.Container(c)
	if err != nil {
		return err
	}
	r := bytes.NewReader([]byte("CheckBucketAccess"))
	item, err := container.Put(".kubed", r, r.Size(), nil)
	if err != nil {
		return err
	}
	if err := container.RemoveItem(item.ID()); err != nil {
		return err
	}
	return nil
}

func NewOSMContext(client clientset.Interface, spec tapi.Backend, namespace string) (*otx.Context, error) {
	config := make(map[string][]byte)

	if spec.StorageSecretName != "" {
		secret, err := client.CoreV1().Secrets(namespace).Get(spec.StorageSecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		config = secret.Data
	}

	nc := &otx.Context{
		Name:   "kubed",
		Config: stow.ConfigMap{},
	}

	if spec.S3 != nil {
		nc.Provider = s3.Kind
		nc.Config[s3.ConfigAccessKeyID] = string(config[tapi.AWS_ACCESS_KEY_ID])
		nc.Config[s3.ConfigEndpoint] = spec.S3.Endpoint
		nc.Config[s3.ConfigRegion] = "us-east-1" // only used for creating buckets
		nc.Config[s3.ConfigSecretKey] = string(config[tapi.AWS_SECRET_ACCESS_KEY])
		if u, err := url.Parse(spec.S3.Endpoint); err == nil {
			nc.Config[s3.ConfigDisableSSL] = strconv.FormatBool(u.Scheme == "http")
		}
		return nc, nil
	} else if spec.GCS != nil {
		nc.Provider = gcs.Kind
		nc.Config[gcs.ConfigProjectId] = string(config[tapi.GOOGLE_PROJECT_ID])
		nc.Config[gcs.ConfigJSON] = string(config[tapi.GOOGLE_SERVICE_ACCOUNT_JSON_KEY])
		return nc, nil
	} else if spec.Azure != nil {
		nc.Provider = azure.Kind
		nc.Config[azure.ConfigAccount] = string(config[tapi.AZURE_ACCOUNT_NAME])
		nc.Config[azure.ConfigKey] = string(config[tapi.AZURE_ACCOUNT_KEY])
		return nc, nil
	} else if spec.Local != nil {
		nc.Provider = local.Kind
		nc.Config[local.ConfigKeyPath] = spec.Local.Path
		return nc, nil
	} else if spec.Swift != nil {
		nc.Provider = swift.Kind
		// https://github.com/restic/restic/blob/master/src/restic/backend/swift/config.go
		for _, val := range []struct {
			stowKey   string
			secretKey string
		}{
			// v2/v3 specific
			{swift.ConfigUsername, tapi.OS_USERNAME},
			{swift.ConfigKey, tapi.OS_PASSWORD},
			{swift.ConfigRegion, tapi.OS_REGION_NAME},
			{swift.ConfigTenantAuthURL, tapi.OS_AUTH_URL},

			// v3 specific
			{swift.ConfigDomain, tapi.OS_USER_DOMAIN_NAME},
			{swift.ConfigTenantName, tapi.OS_PROJECT_NAME},
			{swift.ConfigTenantDomain, tapi.OS_PROJECT_DOMAIN_NAME},

			// v2 specific
			{swift.ConfigTenantId, tapi.OS_TENANT_ID},
			{swift.ConfigTenantName, tapi.OS_TENANT_NAME},

			// v1 specific
			{swift.ConfigTenantAuthURL, tapi.ST_AUTH},
			{swift.ConfigUsername, tapi.ST_USER},
			{swift.ConfigKey, tapi.ST_KEY},

			// Manual authentication
			{swift.ConfigStorageURL, tapi.OS_STORAGE_URL},
			{swift.ConfigAuthToken, tapi.OS_AUTH_TOKEN},
		} {
			if _, exists := nc.Config.Config(val.stowKey); !exists {
				nc.Config[val.stowKey] = string(config[val.secretKey])
			}
		}
		return nc, nil
	}
	return nil, errors.New("No storage provider is configured.")
}
