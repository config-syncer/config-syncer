package cert

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	DefaultCheckInterval time.Duration = 24 * time.Hour
	DefaultMinTTL        time.Duration = 7 * 24 * time.Hour
	CaCertPath           string        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type CertWatcher struct {
	CheckInterval   time.Duration
	MinTTL          time.Duration
	CertPath        string
	KubeClient      clientset.Interface
	SecretName      string
	SecretNamespace string
}

func DefaultCertWatcher(c clientset.Interface, secretName string, secretNamespace string) *CertWatcher {
	return &CertWatcher{
		CheckInterval:   DefaultCheckInterval,
		MinTTL:          DefaultMinTTL,
		CertPath:        CaCertPath,
		KubeClient:      c,
		SecretName:      secretName,
		SecretNamespace: secretNamespace,
	}
}

//This will block executions
func (c CertWatcher) RunAndHold() {
	for range time.NewTicker(c.CheckInterval).C {
		crt, err := c.loadCACert()
		if err != nil {
			c.notify("Failed to load CA cert")
			continue
		}

		var msg string
		if crt.NotAfter.Before(time.Now()) {
			msg = fmt.Sprintf("CA certificate expired at %v, please renew", crt.NotAfter)
		} else if crt.NotAfter.Sub(time.Now()) < c.MinTTL {
			msg = fmt.Sprintf("Certificate will expire within %v days, please renew.", crt.NotAfter.Sub(time.Now()).Hours()/24)
		} else {
			continue
		}
		if uid, err := c.notify(msg); err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v.", uid)
		}
	}
}

func (c CertWatcher) loadCACert() (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(c.CertPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Failed to parse certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (c CertWatcher) getLoader() (func(string) (string, bool), error) {
	cfg, err := c.KubeClient.CoreV1().
		Secrets(c.SecretNamespace).
		Get(c.SecretName, meta_v1.GetOptions{})
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
