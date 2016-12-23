// Package digitalocean implements a DNS provider for solving the DNS-01
// challenge using digitalocean DNS.
package digitalocean

import (
	"errors"
	"fmt"
	"log"
	"strings"

	dp "github.com/appscode/go-dns/provider"
	"github.com/digitalocean/godo"
	"github.com/kelseyhightower/envconfig"
	"github.com/xenolf/lego/acme"
	"golang.org/x/oauth2"
)

// DNSProvider is an implementation of the acme.ChallengeProvider interface
// that uses DigitalOcean's REST API to manage TXT records for a domain.
type DNSProvider struct {
	client *godo.Client
}

type Options struct {
	AuthToken string `json:"auth_token" envconfig:"DO_AUTH_TOKEN"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance configured for Digital
// Ocean. Credentials must be passed in the environment variable:
// DO_AUTH_TOKEN.
func NewDNSProvider() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return NewDNSProviderCredentials(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for Digital Ocean.
func NewDNSProviderCredentials(opt Options) (*DNSProvider, error) {
	if opt.AuthToken == "" {
		return nil, errors.New("DigitalOcean credentials missing")
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: opt.AuthToken,
	}))
	return &DNSProvider{
		client: godo.NewClient(oauthClient),
	}, nil
}

func (c *DNSProvider) EnsureARecord(domain string, ip string) error {
	authZone, err := acme.FindZoneByFqdn(acme.ToFqdn(domain), acme.RecursiveNameservers)
	if err != nil {
		return fmt.Errorf("Could not determine zone for domain: '%s'. %s", domain, err)
	}
	authZone = acme.UnFqdn(authZone)
	relative := toRelativeRecord(domain, authZone)

	records, _, err := c.client.Domains.Records(authZone, &godo.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == relative && record.Data == ip {
			log.Println("DNS is already configured. No DNS related change is necessary.")
			return nil
		}
	}
	_, _, err = c.client.Domains.CreateRecord(authZone, &godo.DomainRecordEditRequest{
		Type: "A",
		Name: relative,
		Data: ip,
	})
	return err
}

func (c *DNSProvider) DeleteARecords(domain string) error {
	authZone, err := acme.FindZoneByFqdn(acme.ToFqdn(domain), acme.RecursiveNameservers)
	if err != nil {
		return fmt.Errorf("Could not determine zone for domain: '%s'. %s", domain, err)
	}
	authZone = acme.UnFqdn(authZone)
	relative := toRelativeRecord(domain, authZone)

	records, _, err := c.client.Domains.Records(authZone, &godo.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == relative {
			_, err = c.client.Domains.DeleteRecord(authZone, record.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Returns the relative record to the domain
func toRelativeRecord(domain, zone string) string {
	return acme.UnFqdn(strings.TrimSuffix(domain, zone))
}
