package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
)

const (
	DefaultCheckInterval    time.Duration = 24 * time.Hour
	DefaultMinRemainingDays time.Duration = 24 * time.Hour * 7
	CaCertPath              string        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type CertWatcher struct {
	CheckInterval                       time.Duration
	MinRemainingDuration                time.Duration
	CertPath                            string
	ClusterKubedConfigSecretMountedPath string
}

func DefaultCertWatcher(clusterKubedSecretMountedPath string) *CertWatcher {
	return &CertWatcher{
		CheckInterval:        DefaultCheckInterval,
		MinRemainingDuration: DefaultMinRemainingDays,
		CertPath:             CaCertPath,
		ClusterKubedConfigSecretMountedPath: clusterKubedSecretMountedPath,
	}
}

func (c CertWatcher) Run() {
	for range time.NewTicker(c.CheckInterval).C {
		f, err := os.Open(c.CertPath)
		if err != nil {
			c.notify(fmt.Sprintf("Certificate  not found in path %v", c.CertPath))
			f.Close()
			continue
		}
		soonExp, days, err := c.isSoonExpired(f, c.MinRemainingDuration)
		f.Close()
		if err != nil {
			c.notify(fmt.Sprintf("Error while parsing certificate: %v", err))
			continue
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
	if int(cert.NotAfter.Sub(time.Now())) < 0 {
		return false, 0, errors.New("certificate already expired")
	}
	if cert.NotAfter.Sub(time.Now()) < minDuration {
		return true, int(cert.NotAfter.Sub(time.Now()) / (24 * time.Hour)), nil
	}
	return false, int(cert.NotAfter.Sub(time.Now()) / (24 * time.Hour)), nil
}

func (c CertWatcher) configuration() (map[string]string, error) {
	return util.MountedSecretToMap(c.ClusterKubedConfigSecretMountedPath)
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
