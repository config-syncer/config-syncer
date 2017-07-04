package util

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMountedSecretToMap(t *testing.T) {
	path := os.Getenv("HOME") + "/temp"
	os.MkdirAll(path, 0777)
	defer os.RemoveAll(path)
	expected := map[string]string{
		"username": "ausername",
		"password": "astrongpassword",
	}
	ioutil.WriteFile(path+"/username", []byte(expected["username"]), 0777)
	ioutil.WriteFile(path+"/password", []byte(expected["password"]), 0777)
	m, err := MountedSecretToMap(path)
	assert.Nil(t, err)
	assert.Equal(t, expected, m)
}
