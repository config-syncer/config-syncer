package provider

type Provider interface {
	EnsureARecord(domain string, ip string) error
	DeleteARecords(domain string) error
}
