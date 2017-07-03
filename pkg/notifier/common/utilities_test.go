package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureRequiredKeys(t *testing.T) {
	mp := map[string]string{
		"host":        "example.com",
		"admin_email": "admin@example.com",
		"admin_phone": "********",
	}
	assert.Nil(t, EnsureRequiredKeys(mp, []string{"host", "admin_email", "admin_phone"}))
	assert.Nil(t, EnsureRequiredKeys(mp, []string{"admin_email", "admin_phone"}))
	assert.NotNil(t, EnsureRequiredKeys(mp, []string{"host", "admin_email", "admin_name"}))

	assert.Equal(t, "admin_name not found", EnsureRequiredKeys(mp, []string{"admin_name"}).Error())
}
