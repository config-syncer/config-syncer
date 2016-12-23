package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var serverDeleteRequestSchema *gojsonschema.Schema
var serverCreateRequestSchema *gojsonschema.Schema

func init() {
	var err error
	serverDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "namespace": {
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
	serverCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "bucket_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "cluster": {
      "type": "string"
    },
    "credential": {
      "type": "string"
    },
    "disk": {
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "title": "TODO will be removed",
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *ServerDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return serverDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ServerDeleteRequest) IsRequest() {}

func (m *ServerCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return serverCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ServerCreateRequest) IsRequest() {}

