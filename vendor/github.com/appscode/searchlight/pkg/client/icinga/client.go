package icinga

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/appscode/k8s-addons/api/install"
)

type IcingaConfig struct {
	Endpoint  string
	BasicAuth struct {
		Username string
		Password string
	}
	CaCert []byte
}

type IcingaClient struct {
	config     *IcingaConfig
	pathPrefix string
}

type IcingaApiRequest struct {
	client *http.Client

	uri      string
	suffix   string
	params   map[string]string
	userName string
	password string
	verb     string

	Err  error
	req  *http.Request
	resp *http.Response

	Status       int
	ResponseBody []byte
}

type IcingaApiResponse struct {
	Err          error
	Status       int
	ResponseBody []byte
}

func NewClient(icingaConfig *IcingaConfig) *IcingaClient {
	c := &IcingaClient{
		config: icingaConfig,
	}
	return c
}

func (c *IcingaClient) SetEndpoint(endpoint string) *IcingaClient {
	c.config.Endpoint = endpoint
	return c
}

func (c *IcingaClient) Objects() *IcingaClient {
	c.pathPrefix = "/objects"
	return c
}

func (c *IcingaClient) Hosts(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/hosts/" + hostName)
}

func (c *IcingaClient) HostGroups(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/hostgroups/" + hostName)
}

func (c *IcingaClient) Service(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/services/" + hostName)
}

func (c *IcingaClient) Actions(action string) *IcingaApiRequest {
	c.pathPrefix = ""
	return c.newIcingaRequest("/actions/" + action)
}

func (c *IcingaClient) Notifications(hostName string) *IcingaApiRequest {
	return c.newIcingaRequest("/notifications/" + hostName)
}

func (c *IcingaClient) Check() *IcingaApiRequest {
	c.pathPrefix = ""
	return c.newIcingaRequest("")
}

func addUri(uri string, name []string) string {
	for _, v := range name {
		uri = uri + "!" + v
	}
	return uri
}

func (r *IcingaApiRequest) Get(name []string, jsonBody ...string) *IcingaApiRequest {
	if len(jsonBody) == 0 {
		r.req, r.Err = r.newRequest("GET", addUri(r.uri, name), nil)
	} else if len(jsonBody) == 1 {
		r.req, r.Err = r.newRequest("GET", addUri(r.uri, name), bytes.NewBuffer([]byte(jsonBody[0])))
	} else {
		r.Err = errors.New("Invalid request")
	}
	return r
}

func (r *IcingaApiRequest) Create(name []string, jsonBody string) *IcingaApiRequest {
	r.req, r.Err = r.newRequest("PUT", addUri(r.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return r
}

func (r *IcingaApiRequest) Update(name []string, jsonBody string) *IcingaApiRequest {
	r.req, r.Err = r.newRequest("POST", addUri(r.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return r
}

func (r *IcingaApiRequest) Delete(name []string, jsonBody string) *IcingaApiRequest {
	r.req, r.Err = r.newRequest("DELETE", addUri(r.uri, name), bytes.NewBuffer([]byte(jsonBody)))
	return r
}

func (r *IcingaApiRequest) Params(param map[string]string) *IcingaApiRequest {
	p := r.req.URL.Query()
	for k, v := range param {
		p.Add(k, v)
	}
	r.req.URL.RawQuery = p.Encode()
	fmt.Println(r.req.URL)
	return r
}
