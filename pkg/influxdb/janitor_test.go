package influx

import (
	"fmt"
	"testing"
	"time"

	"github.com/appscode/kubed/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInfluxJanitor(t *testing.T) {
	host := ""
	user := ""
	pass := ""

	j := Janitor{
		Spec: config.InfluxDBSpec{
			Endpoint: fmt.Sprintf("http://%s:8086", host),
			Username: user,
			Password: pass,
		},
	}
	j.TTL, _ = time.ParseDuration("24h")

	err := j.Cleanup()
	assert.Nil(t, err)
}
