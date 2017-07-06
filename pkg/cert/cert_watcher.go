package cert

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	DefaultCheckInterval    time.Duration = 24 * time.Hour
	DefaultMinRemainingDays int           = 7
	CaCertPath              string        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type ExpiratinStatus int

const (
	ExpiratinStatus_AleadyExpired ExpiratinStatus = 0
	ExpiratinStatus_ExpiredSoon   ExpiratinStatus = 1
	ExpiratinStatus_Fresh         ExpiratinStatus = 2
)

type CertWatcher struct {
	CheckInterval                     time.Duration
	MinRemainingDurationInDays        int
	CertPath                          string
	KubeClient                        clientset.Interface
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
}

func DefaultCertWatcher(c clientset.Interface, secretName string, secretNamespace string) *CertWatcher {
	return &CertWatcher{
		CheckInterval:                     DefaultCheckInterval,
		MinRemainingDurationInDays:        DefaultMinRemainingDays,
		CertPath:                          CaCertPath,
		KubeClient:                        c,
		ClusterKubedConfigSecretName:      secretName,
		ClusterKubedConfigSecretNamespace: secretNamespace,
	}
}

//This will block executions
func (c CertWatcher) RunAndHold() {
	for range time.NewTicker(c.CheckInterval).C {
		f, err := os.Open(c.CertPath)
		if err != nil {
			c.notify(fmt.Sprintf("Certificate  not found in path %v", c.CertPath))
			f.Close()
			continue
		}
		d, err := c.certLifetimeInDays(f)
		f.Close()
		if err != nil {
			c.notify(fmt.Sprintf("Error while parsing certificate: %v", err))
			continue
		}

		var msg string
		switch c.expirationStatus(d, c.MinRemainingDurationInDays) {
		case ExpiratinStatus_AleadyExpired:
			msg = "certificate already expired, please renew"
		case ExpiratinStatus_ExpiredSoon:
			msg = fmt.Sprintf("Certificate will expire within %v days, please renew.", d)
		}
		if uid, err := c.notify(msg); err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v.", uid)
		}
	}
}

func (c CertWatcher) expirationStatus(remD, minD int) ExpiratinStatus {
	if remD <= 0 {
		return ExpiratinStatus_AleadyExpired
	} else if remD <= minD {
		return ExpiratinStatus_ExpiredSoon
	} else {
		return ExpiratinStatus_Fresh
	}
}

func (c CertWatcher) certLifetimeInDays(r io.Reader) (int, error) {
	certData, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return 0, errors.New("failed to parse certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(cert.NotAfter.Sub(time.Now()).Hours() / 24)), nil
}

func (c CertWatcher) getLoader() (func(string) (string, bool), error) {
	cfg, err := c.KubeClient.CoreV1().
		Secrets(c.ClusterKubedConfigSecretNamespace).
		Get(c.ClusterKubedConfigSecretName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = cfg.Data[key]
		value = string(bytes)
		return
	}, nil
}

func (c CertWatcher) notify(msg string) (string, error) {
	loader, err := c.getLoader()
	if err != nil {
		return "", err
	}
	notifier, err := unified.Load(loader)
	if err != nil {
		return "", err
	}
	switch n := notifier.(type) {
	case notify.ByEmail:
		receivers := getArray(loader, "CLUSTER_ADMIN_EMAIL")
		if len(receivers) == 0 {
			return n.UID(), errors.New("Missing / invalid cluster admin email(s)")
		}
		n = n.To(receivers[0], receivers[1:]...)
		return n.UID(), n.WithSubject("Cluster CA Certificate").WithBody(msg).Send()
	case notify.BySMS:
		receivers := getArray(loader, "CLUSTER_ADMIN_PHONE")
		if len(receivers) == 0 {
			return n.UID(), errors.New("Missing / invalid cluster admin phone number(s)")
		}
		n = n.To(receivers[0], receivers[1:]...)
		return n.UID(), n.WithBody(msg).Send()
	case notify.ByChat:
		return n.UID(), n.WithBody(msg).Send()
	}
	return "", errors.New("Unknown notifier")
}

func getArray(loader func(string) (string, bool), key string) []string {
	if v, ok := loader(key); ok {
		vals := make([]string, 0)
		err := json.Unmarshal([]byte(v), &vals)
		if err != nil {
			return []string{}
		}
		return vals
	}
	return []string{}
}
