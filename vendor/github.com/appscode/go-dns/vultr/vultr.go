// Package vultr implements a DNS provider for solving the DNS-01 challenge using
// the vultr DNS.
// See https://www.vultr.com/api/#dns
package vultr

import (
	"errors"
	"fmt"
	"log"
	"strings"

	vultr "github.com/JamesClonk/vultr/lib"
	dp "github.com/appscode/go-dns/provider"
	"github.com/kelseyhightower/envconfig"
	"github.com/xenolf/lego/acme"
)

// DNSProvider is an implementation of the acme.ChallengeProvider interface.
type DNSProvider struct {
	client *vultr.Client
}

type Options struct {
	ApiKey string `json:"api_key" envconfig:"VULTR_API_KEY" form:"vultr_api_key"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance with a configured Vultr client.
// Authentication uses the VULTR_API_KEY environment variable.
func NewDNSProvider() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return NewDNSProviderCredentials(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a DNSProvider
// instance configured for Vultr.
func NewDNSProviderCredentials(opt Options) (*DNSProvider, error) {
	if opt.ApiKey == "" {
		return nil, errors.New("Vultr credentials missing")
	}

	c := &DNSProvider{
		client: vultr.NewClient(opt.ApiKey, nil),
	}

	return c, nil
}

func (c *DNSProvider) EnsureARecord(domain string, ip string) error {
	zoneDomain, err := c.getHostedZone(domain)
	if err != nil {
		return err
	}
	relative := toRelativeRecord(domain, zoneDomain)

	records, err := c.client.GetDNSRecords(zoneDomain)
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == relative && record.Data == ip {
			log.Println("DNS is already configured. No DNS related change is necessary.")
			return nil
		}
	}
	return c.client.CreateDNSRecord(zoneDomain, relative, "A", ip, 0, 300)
}

func (c *DNSProvider) DeleteARecords(domain string) error {
	zoneDomain, err := c.getHostedZone(domain)
	if err != nil {
		return err
	}
	relative := toRelativeRecord(domain, zoneDomain)

	records, err := c.client.GetDNSRecords(zoneDomain)
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == relative {
			err = c.client.DeleteDNSRecord(zoneDomain, record.RecordID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *DNSProvider) getHostedZone(domain string) (string, error) {
	domains, err := c.client.GetDNSDomains()
	if err != nil {
		return "", fmt.Errorf("Vultr API call failed: %v", err)
	}

	var hostedDomain vultr.DNSDomain
	for _, d := range domains {
		if strings.HasSuffix(domain, d.Domain) {
			if len(d.Domain) > len(hostedDomain.Domain) {
				hostedDomain = d
			}
		}
	}
	if hostedDomain.Domain == "" {
		return "", fmt.Errorf("No matching Vultr domain found for domain %s", domain)
	}

	return hostedDomain.Domain, nil
}

func (c *DNSProvider) extractRecordName(fqdn, domain string) string {
	name := acme.UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}
	return name
}

// Returns the relative record to the domain
func toRelativeRecord(domain, zone string) string {
	return acme.UnFqdn(strings.TrimSuffix(domain, zone))
}
