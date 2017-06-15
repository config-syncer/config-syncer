package janitor

import (
	"testing"

	"github.com/appscode/vultr/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func TestConfigMapToClusterSettings(t *testing.T) {
	cnf1 := apiv1.ConfigMap{
		Data: map[string]string{
			"log-index-prefix":            "test-prefix",
			"log-storage-lifetime":        "3333",
			"monitoring-storage-lifetime": "2222",
		},
	}
	expected := ClusterSettings{
		LogIndexPrefix:            "test-prefix",
		LogStorageLifetime:        3333,
		MonitoringStorageLifetime: 2222,
	}
	c, err := ConfigMapToClusterSettings(cnf1)
	assert.Nil(t, err)
	assert.Equal(t, expected, c)

	cnf2 := apiv1.ConfigMap{
		Data: map[string]string{
			"log-index-prefix":            "test-prefix",
			"log-storage-lifetime":        "err-data",
			"monitoring-storage-lifetime": "2222",
		},
	}
	c, err = ConfigMapToClusterSettings(cnf2)
	assert.NotNil(t, err)
}
