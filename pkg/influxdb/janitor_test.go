package influx

import (
	"fmt"
	"testing"
	"time"

	apis "github.com/appscode/kubed/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestInfluxJanitor(t *testing.T) {
	host := ""
	user := ""
	pass := ""

	j := Janitor{
		Spec: apis.InfluxDBSpec{
			Endpoint: fmt.Sprintf("http://%s:8086", host),
			Username: user,
			Password: pass,
		},
	}
	j.TTL, _ = time.ParseDuration("24h")

	err := j.Cleanup()
	assert.Nil(t, err)
}
