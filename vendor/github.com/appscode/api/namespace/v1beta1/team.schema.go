package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var createRequestSchema *gojsonschema.Schema
var getRequestSchema *gojsonschema.Schema
var isAvailableRequestSchema *gojsonschema.Schema

func init() {
	var err error
	createRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "display_name": {
      "type": "string"
    },
    "email": {
      "type": "string"
    },
    "initial_units": {
      "type": "integer"
    },
    "invite_emails": {
      "items": {
        "type": "string"
      },
      "type": "array"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "password": {
      "type": "string"
    },
    "payment_method_nonce": {
      "type": "string"
    },
    "subscription": {
      "type": "string"
    },
    "user_name": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	getRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	isAvailableRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *CreateRequest) IsValid() (*gojsonschema.Result, error) {
	return createRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CreateRequest) IsRequest() {}

func (m *GetRequest) IsValid() (*gojsonschema.Result, error) {
	return getRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *GetRequest) IsRequest() {}

func (m *IsAvailableRequest) IsValid() (*gojsonschema.Result, error) {
	return isAvailableRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IsAvailableRequest) IsRequest() {}

