package es

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/appscode/go/log"
	apis "github.com/appscode/kubed/pkg/apis/v1alpha1"
	"gopkg.in/olivere/elastic.v5"
)

type Janitor struct {
	Spec     apis.ElasticsearchSpec
	AuthInfo *apis.JanitorAuthInfo
	TTL      time.Duration
}

func (j *Janitor) Cleanup() error {
	var httpClient *http.Client

	if j.AuthInfo != nil {
		mTLSConfig := &tls.Config{}
		if j.AuthInfo.CACertData != nil {
			certs := x509.NewCertPool()
			certs.AppendCertsFromPEM(j.AuthInfo.CACertData)
			mTLSConfig.RootCAs = certs
			if j.AuthInfo.ClientCertData != nil && j.AuthInfo.ClientKeyData != nil {
				cert, err := tls.X509KeyPair(j.AuthInfo.ClientCertData, j.AuthInfo.ClientKeyData)
				if err == nil {
					mTLSConfig.Certificates = []tls.Certificate{cert}
				}
			}
		}

		if j.AuthInfo.InsecureSkipVerify {
			mTLSConfig.InsecureSkipVerify = true
		}

		// https://github.com/golang/go/blob/eca45997dfd6cd14a59fbdea2385f6648a0dc786/src/net/http/transport.go#L40
		tr := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       mTLSConfig,
		}
		httpClient = &http.Client{Transport: tr}
	}

	client, err := elastic.NewClient(
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(true),
		elastic.SetURL(j.Spec.Endpoint),
	)
	if err != nil {
		return err
	}

	indices, err := client.IndexNames()
	if err != nil {
		return err
	}

	indicesToDelete := make([]string, 0)

	now := time.Now().UTC()
	oldDate := now.Add(-j.TTL)

	for _, index := range indices {
		if strings.HasPrefix(index, j.Spec.LogIndexPrefix) && len(index) >= 10 { // len("2006.01.02") == 10
			t, err := time.Parse("2006.01.02", index[len(index)-10:])
			if err != nil {
				log.Warningf("Invalid format for Index [%s]", index)
				continue
			}
			if oldDate.After(t) {
				indicesToDelete = append(indicesToDelete, index)
			}
		}
	}

	if len(indicesToDelete) > 0 {
		if _, err := client.DeleteIndex(indicesToDelete...).Do(context.Background()); err != nil {
			log.Errorln(err)
			return err
		}
		log.Debugf("Old Indices [%s] deleted", strings.Join(indicesToDelete, ","))
	}

	log.Debugf("ElasticSearch cleanup process complete")
	return nil
}
