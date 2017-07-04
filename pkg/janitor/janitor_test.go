package janitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClusterSettings(t *testing.T) {
	expected := ClusterSettings{
		MonitoringStorageLifetime: 2222,
		LogStorageLifetime:        3333,
		LogIndexPrefix:            "test-",
	}

	path := os.Getenv("HOME") + "/temp"
	os.MkdirAll(path, 0777)
	defer os.RemoveAll(path)

	ioutil.WriteFile(path+"/"+MonitoringStorageLifetime, []byte(fmt.Sprintf("%v", expected.MonitoringStorageLifetime)), 0777)
	ioutil.WriteFile(path+"/"+LogStorageLifetime, []byte(fmt.Sprintf("%v", expected.LogStorageLifetime)), 0777)
	ioutil.WriteFile(path+"/"+LogIndexPrefix, []byte(expected.LogIndexPrefix), 0777)
	m, err := getClusterSettings(path)
	assert.Nil(t, err)
	assert.Equal(t, expected, m)
}
