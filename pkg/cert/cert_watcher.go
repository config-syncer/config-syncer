package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appscode/go-notify/hipchat"
	"github.com/appscode/go-notify/mailgun"
	"github.com/appscode/go-notify/plivo"
	"github.com/appscode/go-notify/slack"
	"github.com/appscode/go-notify/smtp"
	"github.com/appscode/go-notify/twilio"
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
		switch c.expirationStatus(d, c.MinRemainingDurationInDays) {
		case ExpiratinStatus_AleadyExpired:
			c.notify("certificate already expired, please renew")
		case ExpiratinStatus_ExpiredSoon:
			c.notify(fmt.Sprintf("Certificate will expire within %v days, please renew.", d))
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
	via, ok := conf["notify_via"]
	if !ok {
		log.Errorln("No notifier set")
		return
	}
	switch via {
	case plivo.UID:
		opts := plivo.Options{
			AuthID:    conf["plivo_auth_id"],
			AuthToken: conf["plivo_auth_token"],
			To:        strings.Split(conf["cluster_admin_phone"], ","),
			From:      conf["plivo_from"],
		}
		err := plivo.New(opts).WithBody(msg).Send()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", plivo.UID, msg)
		}
	case twilio.UID:
		opts := twilio.Options{
			AccountSid: conf["twilio_account_sid"],
			AuthToken:  conf["twilio_auth_token"],
			From:       conf["twilio_from"],
			To:         strings.Split(conf["cluster_admin_phone"], ","),
		}
		err := twilio.New(opts).WithBody(msg).Send()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", twilio.UID, msg)
		}
	case smtp.UID:
		port, err := strconv.ParseInt(conf["smtp_port"], 10, 64)
		if err != nil {
			log.Errorln(err)
			return
		}
		opts := smtp.Options{
			Host:     conf["smtp_host"],
			Port:     int(port),
			Username: conf["smtp_username"],
			Password: conf["smtp_password"],
			From:     conf["smtp_from"],
			To:       strings.Split(conf["cluster_admin_email"], ","),
		}
		err = smtp.New(opts).WithSubject("kubed notification").WithBody(msg).SendHtml()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", smtp.UID, msg)
		}
	case mailgun.UID:
		opts := mailgun.Options{
			Domain:       conf["mailgun_domain"],
			ApiKey:       conf["mailgun_api_key"],
			PublicApiKey: conf["mailgun_public_api_key"],
			From:         conf["mailgun_from"],
			To:           strings.Split(conf["cluster_admin_email"], ","),
		}
		err := mailgun.New(opts).WithSubject("kubed notification").WithBody(msg).SendHtml()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", mailgun.UID, msg)
		}

	case hipchat.UID:
		opts := hipchat.Options{
			AuthToken: conf["hipchat_auth_token"],
			To:        strings.Split(conf["hipchat_room"], ","),
		}
		err := hipchat.New(opts).WithBody(msg).Send()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", hipchat.UID, msg)
		}
	case slack.UID:
		opts := slack.Options{
			AuthToken: conf["slack_auth_token"],
			Channel:   strings.Split(conf["slack_channel"], ","),
		}
		err := slack.New(opts).WithBody(msg).Send()
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugf("Notification successfully sent via %v. Notification: `%v`", slack.UID, msg)
		}
	}
}
