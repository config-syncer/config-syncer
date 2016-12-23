package util

import (
	"encoding/base64"
)

const (
	NAMESPACE_SEPARATOR = "."
	TOKEN_SEPARATOR     = ":"
)

// Basic auth encoder for the appscode schema
func BasicEncode(namespace, username, password string) string {
	data := namespace + NAMESPACE_SEPARATOR + username + TOKEN_SEPARATOR + password
	return base64.StdEncoding.EncodeToString([]byte(data))
}
