// Package linode implements a DNS provider for solving the DNS-01 challenge
// using Linode DNS.
package linode

import (
	"errors"
	"log"
	"strings"
	"time"

	dp "github.com/appscode/go-dns/provider"
	"github.com/kelseyhightower/envconfig"
	"github.com/timewasted/linode/dns"
	"github.com/xenolf/lego/acme"
)

const (
	dnsMinTTLSecs      = 300
	dnsUpdateFreqMins  = 15
	dnsUpdateFudgeSecs = 120
)

type hostedZoneInfo struct {
	domainId     int
	resourceName string
}

// DNSProvider implements the acme.ChallengeProvider interface.
type DNSProvider struct {
	linode *dns.DNS
}

type Options struct {
	ApiKey string `json:"api_key" envconfig:"LINODE_API_KEY" form:"linode_api_key"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance configured for Linode.
// Credentials must be passed in the environment variable: LINODE_API_KEY.
func NewDNSProvider() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return NewDNSProviderCredentials(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for Linode.
func NewDNSProviderCredentials(opt Options) (*DNSProvider, error) {
	if len(opt.ApiKey) == 0 {
		return nil, errors.New("Linode credentials missing")
	}

	return &DNSProvider{
		linode: dns.New(opt.ApiKey),
	}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS
// propagation.  Adjusting here to cope with spikes in propagation times.
func (p *DNSProvider) Timeout() (timeout, interval time.Duration) {
	// Since Linode only updates their zone files every X minutes, we need
	// to figure out how many minutes we have to wait until we hit the next
	// interval of X.  We then wait another couple of minutes, just to be
	// safe.  Hopefully at some point during all of this, the record will
	// have propagated throughout Linode's network.
	minsRemaining := dnsUpdateFreqMins - (time.Now().Minute() % dnsUpdateFreqMins)

	timeout = (time.Duration(minsRemaining) * time.Minute) +
		(dnsMinTTLSecs * time.Second) +
		(dnsUpdateFudgeSecs * time.Second)
	interval = 15 * time.Second
	return
}

func (p *DNSProvider) EnsureARecord(domain string, ip string) error {
	zone, err := p.getHostedZoneInfo(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	records, err := p.linode.GetResourcesByType(zone.domainId, "A")
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == zone.resourceName && record.Target == ip {
			log.Println("DNS is already configured. No DNS related change is necessary.")
			return nil
		}
	}
	_, err = p.linode.CreateDomainResourceA(zone.domainId, zone.resourceName, ip, 300)
	return err
}

func (p *DNSProvider) DeleteARecords(domain string) error {
	zone, err := p.getHostedZoneInfo(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	records, err := p.linode.GetResourcesByType(zone.domainId, "A")
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == zone.resourceName {
			_, err = p.linode.DeleteDomainResource(record.DomainID, record.ResourceID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *DNSProvider) getHostedZoneInfo(fqdn string) (*hostedZoneInfo, error) {
	// Lookup the zone that handles the specified FQDN.
	authZone, err := acme.FindZoneByFqdn(fqdn, acme.RecursiveNameservers)
	if err != nil {
		return nil, err
	}
	resourceName := strings.TrimSuffix(fqdn, "."+authZone)

	// Query the authority zone.
	domain, err := p.linode.GetDomain(acme.UnFqdn(authZone))
	if err != nil {
		return nil, err
	}

	return &hostedZoneInfo{
		domainId:     domain.DomainID,
		resourceName: resourceName,
	}, nil
}
