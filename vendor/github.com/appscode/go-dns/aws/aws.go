// Package route53 implements a DNS provider for solving the DNS-01 challenge
// using AWS Route 53 DNS.
package aws

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	dp "github.com/appscode/go-dns/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/xenolf/lego/acme"
)

const (
	maxRetries = 5
	ttl        = 300
)

// DNSProvider implements the acme.ChallengeProvider interface
type DNSProvider struct {
	client *route53.Route53
}

var _ dp.Provider = &DNSProvider{}

// customRetryer implements the client.Retryer interface by composing the
// DefaultRetryer. It controls the logic for retrying recoverable request
// errors (e.g. when rate limits are exceeded).
type customRetryer struct {
	client.DefaultRetryer
}

// RetryRules overwrites the DefaultRetryer's method.
// It uses a basic exponential backoff algorithm that returns an initial
// delay of ~400ms with an upper limit of ~30 seconds which should prevent
// causing a high number of consecutive throttling errors.
// For reference: Route 53 enforces an account-wide(!) 5req/s query limit.
func (d customRetryer) RetryRules(r *request.Request) time.Duration {
	retryCount := r.RetryCount
	if retryCount > 7 {
		retryCount = 7
	}

	delay := (1 << uint(retryCount)) * (rand.Intn(50) + 200)
	return time.Duration(delay) * time.Millisecond
}

type Options struct {
	AccessKeyId     string `json:"access_key_id" envconfig:"AWS_ACCESS_KEY_ID" form:"aws_access_key_id"`
	SecretAccessKey string `json:"secret_access_key" envconfig:"AWS_SECRET_ACCESS_KEY" form:"aws_secret_access_key"`
	Region          string `json:"region" envconfig:"AWS_REGION" form:"aws_region"`
}

// NewDNSProvider returns a DNSProvider instance configured for the AWS
// Route 53 service.
//
// AWS Credentials are automatically detected in the following locations
// and prioritized in the following order:
// 1. Environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY,
//    AWS_REGION, [AWS_SESSION_TOKEN]
// 2. Shared credentials file (defaults to ~/.aws/credentials)
// 3. Amazon EC2 IAM role
//
// See also: https://github.com/aws/aws-sdk-go/wiki/configuring-sdk
func NewDNSProvider() (*DNSProvider, error) {
	r := customRetryer{}
	r.NumMaxRetries = maxRetries
	config := request.WithRetryer(aws.NewConfig(), r)
	client := route53.New(session.New(config))

	return &DNSProvider{client: client}, nil
}

func NewDNSProviderCredentials(opt Options) (*DNSProvider, error) {
	r := customRetryer{}
	r.NumMaxRetries = maxRetries
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(opt.AccessKeyId, opt.SecretAccessKey, " "),
		Region:      aws.String(opt.Region),
		Retryer:     r,
	}
	client := route53.New(session.New(config))
	return &DNSProvider{client: client}, nil
}

func (r *DNSProvider) EnsureARecord(domain string, ip string) error {
	fqdn := acme.ToFqdn(domain)
	hostedZoneID, err := getHostedZoneID(fqdn, r.client)
	if err != nil {
		return fmt.Errorf("Failed to determine Route 53 hosted zone ID: %v", err)
	}

	resp, err := r.client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneID),
		StartRecordName: aws.String(fqdn),
		StartRecordType: aws.String(route53.RRTypeA),
	})
	if err != nil {
		return err
	}

	log.Println("Updating A record for cluster", domain)
	reqParams := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("Managed by AppsCode"),
			Changes: make([]*route53.Change, 0),
		},
	}
	if len(resp.ResourceRecordSets) == 0 || !contains(resp.ResourceRecordSets[0].ResourceRecords, ip) {
		rrecords := []*route53.ResourceRecord{
			{
				Value: aws.String(ip),
			},
		}
		if len(resp.ResourceRecordSets) > 0 {
			rrecords = append(rrecords, resp.ResourceRecordSets[0].ResourceRecords...)
		}

		log.Println("Adding A record ", []string{ip})
		reqParams.ChangeBatch.Changes = append(reqParams.ChangeBatch.Changes, &route53.Change{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(fqdn),
				Type:            aws.String(route53.RRTypeA),
				ResourceRecords: rrecords,
				TTL:             aws.Int64(ttl),
			},
		})
	}
	if len(reqParams.ChangeBatch.Changes) == 0 {
		log.Println("DNS is already configured. No DNS related change is necessary.")
		return nil
	}
	_, err = r.client.ChangeResourceRecordSets(reqParams)
	return err
}

func (r *DNSProvider) DeleteARecords(domain string) error {
	fqdn := acme.ToFqdn(domain)
	hostedZoneID, err := getHostedZoneID(fqdn, r.client)
	if err != nil {
		return fmt.Errorf("Failed to determine Route 53 hosted zone ID: %v", err)
	}

	resp, err := r.client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneID),
		StartRecordName: aws.String(fqdn),
		StartRecordType: aws.String(route53.RRTypeA),
	})
	if err != nil {
		return err
	}
	if len(resp.ResourceRecordSets) == 0 {
		log.Println("No A record found. No DNS related change is necessary.")
		return nil
	}
	recordSet := &route53.ResourceRecordSet{
		Name:            aws.String(fqdn),
		Type:            aws.String(route53.RRTypeA),
		TTL:             aws.Int64(ttl),
		ResourceRecords: resp.ResourceRecordSets[0].ResourceRecords,
	}
	reqParams := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("Managed by AppsCode"),
			Changes: []*route53.Change{
				{
					Action:            aws.String(route53.ChangeActionDelete),
					ResourceRecordSet: recordSet,
				},
			},
		},
	}
	_, err = r.client.ChangeResourceRecordSets(reqParams)
	return err
}

func getHostedZoneID(fqdn string, client *route53.Route53) (string, error) {
	authZone, err := acme.FindZoneByFqdn(fqdn, acme.RecursiveNameservers)
	if err != nil {
		return "", err
	}

	// .DNSName should not have a trailing dot
	reqParams := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(acme.UnFqdn(authZone)),
	}
	resp, err := client.ListHostedZonesByName(reqParams)
	if err != nil {
		return "", err
	}

	var hostedZoneID string
	for _, hostedZone := range resp.HostedZones {
		// .Name has a trailing dot
		if !*hostedZone.Config.PrivateZone && *hostedZone.Name == authZone {
			hostedZoneID = *hostedZone.Id
			break
		}
	}

	if len(hostedZoneID) == 0 {
		return "", fmt.Errorf("Zone %s not found in Route 53 for domain %s", authZone, fqdn)
	}

	if strings.HasPrefix(hostedZoneID, "/hostedzone/") {
		hostedZoneID = strings.TrimPrefix(hostedZoneID, "/hostedzone/")
	}

	return hostedZoneID, nil
}

func contains(records []*route53.ResourceRecord, s string) bool {
	for _, record := range records {
		if *record.Value == s {
			return true
		}
	}
	return false
}
