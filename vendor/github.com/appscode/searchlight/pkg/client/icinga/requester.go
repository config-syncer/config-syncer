package icinga

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (r *IcingaApiRequest) Do() *IcingaApiResponse {
	if r.Err != nil {
		return &IcingaApiResponse{
			Err: r.Err,
		}
	}
	r.req.Header.Set("Accept", "application/json")

	if r.userName != "" && r.password != "" {
		r.req.SetBasicAuth(r.userName, r.password)
	}

	r.resp, r.Err = r.client.Do(r.req)
	if r.Err != nil {
		return &IcingaApiResponse{
			Err: r.Err,
		}
	}

	r.Status = r.resp.StatusCode
	r.ResponseBody, r.Err = ioutil.ReadAll(r.resp.Body)
	if r.Err != nil {
		return &IcingaApiResponse{
			Err: r.Err,
		}
	}
	return &IcingaApiResponse{
		Status:       r.Status,
		ResponseBody: r.ResponseBody,
	}
}

func (r *IcingaApiResponse) Into(to interface{}) (int, error) {
	if r.Err != nil {
		return r.Status, r.Err
	}
	err := json.Unmarshal(r.ResponseBody, to)
	if err != nil {
		return r.Status, err
	}
	return r.Status, nil
}

func (c *IcingaClient) newIcingaRequest(path string) *IcingaApiRequest {
	mTLSConfig := &tls.Config{}

	if c.config.CaCert != nil {
		certs := x509.NewCertPool()
		certs.AppendCertsFromPEM(c.config.CaCert)
		mTLSConfig.RootCAs = certs
	} else {
		mTLSConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		TLSClientConfig: mTLSConfig,
	}
	client := &http.Client{Transport: tr}

	c.pathPrefix = c.pathPrefix + path
	return &IcingaApiRequest{
		uri:      c.config.Endpoint + c.pathPrefix,
		client:   client,
		userName: c.config.BasicAuth.Username,
		password: c.config.BasicAuth.Password,
	}
}

func (c *IcingaApiRequest) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	if strings.HasSuffix(urlStr, "/") {
		urlStr = strings.TrimRight(urlStr, "/")
	}

	return http.NewRequest(method, urlStr, body)
}
