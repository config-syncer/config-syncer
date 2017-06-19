package janitor

import (
	"testing"

	"github.com/appscode/vultr/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func TestConfigMapToClusterSettings(t *testing.T) {
	cnf1 := apiv1.Secret{
		Data: map[string][]byte{
			"log-index-prefix":            []byte("test-prefix"),
			"log-storage-lifetime":        []byte("3333"),
			"monitoring-storage-lifetime": []byte("2222"),
		},
	}
	expected := ClusterSettings{
		LogIndexPrefix:            "test-prefix",
		LogStorageLifetime:        3333,
		MonitoringStorageLifetime: 2222,
	}
	c, err := SecretToClusterSettings(cnf1)
	assert.Nil(t, err)
	assert.Equal(t, expected, c)

	cnf2 := apiv1.Secret{
		Data: map[string][]byte{
			"log-index-prefix":            []byte("test-prefix"),
			"log-storage-lifetime":        []byte("err-data"),
			"monitoring-storage-lifetime": []byte("2222"),
		},
	}
	c, err = SecretToClusterSettings(cnf2)
	assert.NotNil(t, err)
}
