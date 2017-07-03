package notifier

import (
	"testing"

	_ "github.com/appscode/kubed/pkg/notifier/plivo"
	_ "github.com/appscode/kubed/pkg/notifier/smtp"
	"github.com/stretchr/testify/assert"
)

func TestNotificationDriver(t *testing.T) {
	opts1 := map[string]string{
		"notify_via":          "plivo",
		"plivo_auth_id":       "auth_id",
		"plivo_auth_token":    "auth_token",
		"cluster_admin_phone": "admin,0111",
		"plivo_from":          "server",
	}
	driver, err := New(opts1).notificationDriver()
	assert.Nil(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, "plivo", driver.Uid())

	opts2 := map[string]string{
		"notify_via":       "unknown",
		"plivo_auth_id":    "auth_id",
		"plivo_auth_token": "auth_token",
		"plivo_to":         "admin,0111",
		"plivo_from":       "server",
	}

	_, err = New(opts2).notificationDriver()
	assert.NotNil(t, err)

	opts3 := map[string]string{
		"notify_via":       "plivo",
		"plivo_auth_token": "auth_token",
		"plivo_to":         "admin,0111",
		"plivo_from":       "server",
	}
	driver, err = New(opts3).notificationDriver()
	assert.NotNil(t, err)
}
