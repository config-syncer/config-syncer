package storage

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	stringz "github.com/appscode/go/strings"
	"github.com/appscode/go/types"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	otx "github.com/appscode/osm/context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	_s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	gcs "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	store "kmodules.xyz/objectstore-api/api/v1"
)

const (
	SecretMountPath = "/etc/osm"
)

func WriteOSMConfig(client kubernetes.Interface, snapshot store.Backend, namespace string, path string) error {
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

func CheckBucketAccess(client kubernetes.Interface, spec store.Backend, namespace string) error {
	cfg, err := NewOSMContext(client, spec, namespace)
	if err != nil {
		return err
	}
	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	c, err := api.Container(spec)
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

func NewOSMContext(client kubernetes.Interface, spec store.Backend, namespace string) (*otx.Context, error) {
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

		keyID, foundKeyID := config[store.AWS_ACCESS_KEY_ID]
		key, foundKey := config[store.AWS_SECRET_ACCESS_KEY]
		if foundKey && foundKeyID {
			nc.Config[s3.ConfigAccessKeyID] = string(keyID)
			nc.Config[s3.ConfigSecretKey] = string(key)
			nc.Config[s3.ConfigAuthType] = "accesskey"
		} else {
			nc.Config[s3.ConfigAuthType] = "iam"
		}
		if strings.HasSuffix(spec.S3.Endpoint, ".amazonaws.com") {
			// find region
			var sess *session.Session
			var err error
			if nc.Config[s3.ConfigAuthType] == "iam" {
				sess, err = session.NewSessionWithOptions(session.Options{
					Config: *aws.NewConfig(),
					// Support MFA when authing using assumed roles.
					SharedConfigState:       session.SharedConfigEnable,
					AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
				})
			} else {
				config := &aws.Config{
					Credentials: credentials.NewStaticCredentials(string(keyID), string(key), ""),
					Region:      aws.String("us-east-1"),
				}
				sess, err = session.NewSessionWithOptions(session.Options{
					Config: *config,
					// Support MFA when authing using assumed roles.
					SharedConfigState:       session.SharedConfigEnable,
					AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
				})
			}
			if err != nil {
				return nil, err
			}
			svc := _s3.New(sess)
			out, err := svc.GetBucketLocation(&_s3.GetBucketLocationInput{
				Bucket: types.StringP(spec.S3.Bucket),
			})
			nc.Config[s3.ConfigRegion] = stringz.Val(types.String(out.LocationConstraint), "us-east-1")
		} else {
			nc.Config[s3.ConfigEndpoint] = spec.S3.Endpoint
			u, err := url.Parse(spec.S3.Endpoint)
			if err != nil {
				return nil, err
			}
			nc.Config[s3.ConfigDisableSSL] = strconv.FormatBool(u.Scheme == "http")

			cacertData, ok := config[store.CA_CERT_DATA]
			if ok && u.Scheme == "https" {
				certFileName := filepath.Join(SecretMountPath, "ca.crt")
				err = os.MkdirAll(filepath.Dir(certFileName), 0755)
				if err != nil {
					return nil, err
				}
				err = ioutil.WriteFile(certFileName, cacertData, 0755)
				if err != nil {
					return nil, err
				}
				nc.Config[s3.ConfigCACertFile] = certFileName
			}
		}
		return nc, nil
	} else if spec.GCS != nil {
		nc.Provider = gcs.Kind
		nc.Config[gcs.ConfigProjectId] = string(config[store.GOOGLE_PROJECT_ID])
		nc.Config[gcs.ConfigJSON] = string(config[store.GOOGLE_SERVICE_ACCOUNT_JSON_KEY])
		return nc, nil
	} else if spec.Azure != nil {
		nc.Provider = azure.Kind
		nc.Config[azure.ConfigAccount] = string(config[store.AZURE_ACCOUNT_NAME])
		nc.Config[azure.ConfigKey] = string(config[store.AZURE_ACCOUNT_KEY])
		return nc, nil
	} else if spec.Local != nil {
		nc.Provider = local.Kind
		nc.Config[local.ConfigKeyPath] = spec.Local.MountPath
		return nc, nil
	} else if spec.Swift != nil {
		nc.Provider = swift.Kind
		// https://github.com/restic/restic/blob/master/src/restic/backend/swift/config.go
		for _, val := range []struct {
			stowKey   string
			secretKey string
		}{
			// v2/v3 specific
			{swift.ConfigUsername, store.OS_USERNAME},
			{swift.ConfigKey, store.OS_PASSWORD},
			{swift.ConfigRegion, store.OS_REGION_NAME},
			{swift.ConfigTenantAuthURL, store.OS_AUTH_URL},

			// v3 specific
			{swift.ConfigDomain, store.OS_USER_DOMAIN_NAME},
			{swift.ConfigTenantName, store.OS_PROJECT_NAME},
			{swift.ConfigTenantDomain, store.OS_PROJECT_DOMAIN_NAME},

			// v2 specific
			{swift.ConfigTenantId, store.OS_TENANT_ID},
			{swift.ConfigTenantName, store.OS_TENANT_NAME},

			// v1 specific
			{swift.ConfigTenantAuthURL, store.ST_AUTH},
			{swift.ConfigUsername, store.ST_USER},
			{swift.ConfigKey, store.ST_KEY},

			// Manual authentication
			{swift.ConfigStorageURL, store.OS_STORAGE_URL},
			{swift.ConfigAuthToken, store.OS_AUTH_TOKEN},
		} {
			if _, exists := nc.Config.Config(val.stowKey); !exists {
				nc.Config[val.stowKey] = string(config[val.secretKey])
			}
		}
		return nc, nil
	}
	return nil, errors.New("No storage provider is configured.")
}
