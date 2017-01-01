// Package cloudflare implements a DNS provider for solving the DNS-01
// challenge using cloudflare DNS.
package cloudflare

import (
	"errors"
	"net/http"
	"time"

	dp "github.com/appscode/go-dns/provider"
	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/xenolf/lego/acme"
)

// DNSProvider is an implementation of the acme.ChallengeProvider interface
type DNSProvider struct {
	api *cf.API
}

type Options struct {
	Email  string `json:"email" envconfig:"CLOUDFLARE_EMAIL" form:"cloudflare_email"`
	APIKey string `json:"api_key" envconfig:"CLOUDFLARE_API_KEY" form:"cloudflare_api_key"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance configured for cloudflare.
// Credentials must be passed in the environment variables: CLOUDFLARE_EMAIL
// and CLOUDFLARE_API_KEY.
func NewDNSProvider() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return NewDNSProviderCredentials(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for cloudflare.
func NewDNSProviderCredentials(opt Options) (*DNSProvider, error) {
	if opt.Email == "" || opt.APIKey == "" {
		return nil, errors.New("CloudFlare credentials missing")
	}

	api, err := cf.New(opt.APIKey, opt.Email, cf.HTTPClient(&http.Client{Timeout: time.Second * 10}))
	if err != nil {
		return nil, err
	}
	return &DNSProvider{api: api}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS
// propagation. Adjusting here to cope with spikes in propagation times.
func (c *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 120 * time.Second, 2 * time.Second
}

func (c *DNSProvider) getHostedZoneID(fqdn string) (string, error) {
	authZone, err := acme.FindZoneByFqdn(fqdn, acme.RecursiveNameservers)
	if err != nil {
		return "", err
	}

	return c.api.ZoneIDByName(acme.UnFqdn(authZone))
}

func (c *DNSProvider) EnsureARecord(domain string, ip string) error {
	zoneID, err := c.getHostedZoneID(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	records, err := c.api.DNSRecords(zoneID, cf.DNSRecord{
		Type:    "A",
		Name:    domain,
		Content: ip,
	})
	if err != nil {
		return err
	}
	if len(records) == 0 {
		_, err = c.api.CreateDNSRecord(zoneID, cf.DNSRecord{
			Type:    "A",
			Name:    domain,
			Content: ip,
			TTL:     300,
		})
		return err
	}
	return nil
}

func (c *DNSProvider) DeleteARecords(domain string) error {
	zoneID, err := c.getHostedZoneID(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	records, err := c.api.DNSRecords(zoneID, cf.DNSRecord{
		Type: "A",
		Name: domain,
	})
	if err != nil {
		return err
	}
	for _, record := range records {
		err = c.api.DeleteDNSRecord(zoneID, record.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
