package smtp

import (
	"testing"

	"github.com/appscode/go-notify/smtp"
	"github.com/stretchr/testify/assert"
)

func TestGetOptions(t *testing.T) {
	b := biblio{}
	expOpt := smtp.Options{
		Host:     "www.example.com",
		Port:     3344,
		Username: "admin",
		Password: "demo_pass",
		From:     "server@example.com",
		To:       []string{"admin@example.com"},
	}
	ts := map[string]string{
		"smtp_host":           "www.example.com",
		"smtp_port":           "3344",
		"smtp_username":       "admin",
		"smtp_password":       "demo_pass",
		"smtp_from":           "server@example.com",
		"cluster_admin_email": "admin@example.com",
	}
	err := b.SetOptions(ts)
	assert.Nil(t, err)
	assert.Equal(t, expOpt, b.opts)

	err = b.SetOptions(map[string]string{})
	assert.NotNil(t, err)
}
