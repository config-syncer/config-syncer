package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/prometheus/common/log"
)

const (
	DefaultCheckInterval    time.Duration = 24 * time.Hour
	DefaultMinRemainingDays time.Duration = 24 * time.Hour * 7
	CaCertPath              string        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type CertWatcher struct {
	CheckInterval        time.Duration
	MinRemainingDuration time.Duration
	CertPath             string
}

func DefaultCertWatcher() *CertWatcher {
	return &CertWatcher{
		CheckInterval:        DefaultCheckInterval,
		MinRemainingDuration: DefaultMinRemainingDays,
		CertPath:             CaCertPath,
	}
}

func (c CertWatcher) Run() {
	for range time.NewTicker(c.CheckInterval).C {
		f, err := os.Open(c.CertPath)
		if err != nil {
			//Notify admin that certificate not found
			log.Errorln(err)
		}
		defer f.Close()
		soonExp, err := c.isSoonExpired(f, c.MinRemainingDuration)
		if err != nil {
			//Notify admin that errors while parsing certificate
			log.Errorln(err)
		}
		if soonExp {
			//Notify admin that certificate will expire soon
		}
	}
}

func (c CertWatcher) isSoonExpired(r io.Reader, minDuration time.Duration) (bool, error) {
	certData, err := ioutil.ReadAll(r)
	if err != nil {
		return false, err
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return false, errors.New("failed to parse certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, err
	}

	if cert.NotAfter.Sub(time.Now()) < minDuration {
		return true, nil
	}
	return false, nil
}
