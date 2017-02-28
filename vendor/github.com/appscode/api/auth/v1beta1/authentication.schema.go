package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var tokenRequestSchema *gojsonschema.Schema
var logoutRequestSchema *gojsonschema.Schema
var loginRequestSchema *gojsonschema.Schema

func init() {
	var err error
	tokenRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "token": {
      "type": "string"
    }
  },
  "title": "Next Id 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	logoutRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "token": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	loginRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "issue_token": {
      "type": "boolean"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "secret": {
      "type": "string"
    },
    "username": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *TokenRequest) IsValid() (*gojsonschema.Result, error) {
	return tokenRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *TokenRequest) IsRequest() {}

func (m *LogoutRequest) IsValid() (*gojsonschema.Result, error) {
	return logoutRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *LogoutRequest) IsRequest() {}

func (m *LoginRequest) IsValid() (*gojsonschema.Result, error) {
	return loginRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *LoginRequest) IsRequest() {}

func (m *CSRFTokenResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *TokenResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *LogoutResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *LoginResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
