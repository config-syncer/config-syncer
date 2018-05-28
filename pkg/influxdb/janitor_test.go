package influx

import (
	"fmt"
	"testing"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func XTestInfluxJanitor(t *testing.T) {
	host := ""
	user := ""
	pass := ""

	j := Janitor{
		Spec: api.InfluxDBSpec{
			Endpoint: fmt.Sprintf("http://%s:8086", host),
			Username: user,
			Password: pass,
		},
	}
	j.TTL, _ = time.ParseDuration("24h")

	err := j.Cleanup()
	assert.Nil(t, err)
}
