package plivo

import (
	"strings"
	"testing"

	"github.com/appscode/go-notify/plivo"
	"github.com/stretchr/testify/assert"
)

func TestGetOptions(t *testing.T) {
	b := biblio{}
	expOpt := plivo.Options{
		AuthID:    "auth_id",
		AuthToken: "auth-token",
		To:        []string{"admin", "maintainer", "1122"},
		From:      "server",
	}
	opts := map[string]string{
		"plivo_auth_id":       expOpt.AuthID,
		"plivo_auth_token":    expOpt.AuthToken,
		"cluster_admin_phone": strings.Join(expOpt.To, ","),
		"plivo_from":          expOpt.From,
	}
	err := b.SetOptions(opts)
	assert.Nil(t, err)
	assert.Equal(t, expOpt, b.opts)

	err = b.SetOptions(map[string]string{})
	assert.NotNil(t, err)
}
