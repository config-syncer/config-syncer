package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/appscode/kubed/pkg/notifier"
	"github.com/appscode/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	DefaultCheckInterval    time.Duration = 24 * time.Hour
	DefaultMinRemainingDays time.Duration = 24 * time.Hour * 7
	CaCertPath              string        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type CertWatcher struct {
	CheckInterval                     time.Duration
	MinRemainingDuration              time.Duration
	CertPath                          string
	KubeClient                        clientset.Interface
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
}

func DefaultCertWatcher(c clientset.Interface, secretName string, secretNamespace string) *CertWatcher {
	return &CertWatcher{
		CheckInterval:                     DefaultCheckInterval,
		MinRemainingDuration:              DefaultMinRemainingDays,
		CertPath:                          CaCertPath,
		KubeClient:                        c,
		ClusterKubedConfigSecretName:      secretName,
		ClusterKubedConfigSecretNamespace: secretNamespace,
	}
}

func (c CertWatcher) Run() {
	f, err := os.Open(c.CertPath)
	if err != nil {
		c.notify(fmt.Sprintf("Certificate  not found in path %v", c.CertPath))
		return
	}
	defer f.Close()

	for range time.NewTicker(c.CheckInterval).C {
		soonExp, days, err := c.isSoonExpired(f, c.MinRemainingDuration)
		if err != nil {
			c.notify(fmt.Sprintf("Error while parsing certificate: %v", err))
		}
		if soonExp {
			c.notify(fmt.Sprintf("Certificate will expire within %v days, please renew.", days))
		}
	}
}

func (c CertWatcher) isSoonExpired(r io.Reader, minDuration time.Duration) (bool, int, error) {
	certData, err := ioutil.ReadAll(r)
	if err != nil {
		return false, 0, err
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return false, 0, errors.New("failed to parse certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, 0, err
	}

	if cert.NotAfter.Sub(time.Now()) < minDuration {
		return true, int(cert.NotAfter.Sub(time.Now()) / (24 * time.Hour)), nil
	}
	return false, int(cert.NotAfter.Sub(time.Now()) / (24 * time.Hour)), nil
}

func (c CertWatcher) configuration() (map[string]string, error) {
	clusterConf, err := c.KubeClient.Core().
		Secrets(c.ClusterKubedConfigSecretNamespace).
		Get(c.ClusterKubedConfigSecretName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data := make(map[string]string, len(clusterConf.Data))
	for key, val := range clusterConf.Data {
		data[key] = string(val)
	}
	return data, nil
}

func (c CertWatcher) notify(msg string) {
	conf, err := c.configuration()
	if err != nil {
		log.Errorln(err)
		return
	}
	notifyVia, err := notifier.New(conf).SendNotification(msg)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Debugf("Notification successfully sent via %v. Notification: `%v`", notifyVia, msg)
	}
}
