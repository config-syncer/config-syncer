package notifier

import (
	"testing"

	_ "github.com/appscode/kubed/pkg/notifier/plivo"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func TestNotifierDriver(t *testing.T) {
	s := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysecret",
			Namespace: "kube-system",
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"username":         []byte("username"),
			"password":         []byte("password"),
			"notify_via":       []byte("plivo"),
			"plivo_auth_id":    []byte("auth_id"),
			"plivo_auth_token": []byte("auth_token"),
			"plivo_to":         []byte("admin,0111"),
			"plivo_from":       []byte("server"),
		},
	}
	driver, err := notificationDriver(fake.NewSimpleClientset(s), s.ObjectMeta.Name, s.ObjectMeta.Namespace)
	assert.Nil(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, "plivo", driver.Uid())
}

func TestNotifierDriverFromConfiguration(t *testing.T) {
	opts1 := map[string][]byte{
		"notify_via":       []byte("plivo"),
		"plivo_auth_id":    []byte("auth_id"),
		"plivo_auth_token": []byte("auth_token"),
		"plivo_to":         []byte("admin,0111"),
		"plivo_from":       []byte("server"),
	}
	driver, err := notifierDriverFromConfiguration(opts1)
	assert.Nil(t, err)
	assert.NotNil(t, driver)

	opts2 := map[string][]byte{
		"notify_via":       []byte("unknown"),
		"plivo_auth_id":    []byte("auth_id"),
		"plivo_auth_token": []byte("auth_token"),
		"plivo_to":         []byte("admin,0111"),
		"plivo_from":       []byte("server"),
	}

	_, err = notifierDriverFromConfiguration(opts2)
	assert.NotNil(t, err)

	opts3 := map[string][]byte{
		"notify_via":       []byte("plivo"),
		"plivo_auth_token": []byte("auth_token"),
		"plivo_to":         []byte("admin,0111"),
		"plivo_from":       []byte("server"),
	}

	_, err = notifierDriverFromConfiguration(opts3)
	assert.NotNil(t, err)
}
