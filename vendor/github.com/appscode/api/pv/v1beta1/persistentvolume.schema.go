package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var pVUnregisterRequestSchema *gojsonschema.Schema
var pVRegisterRequestSchema *gojsonschema.Schema
var pVDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	pVUnregisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
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
	pVRegisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "endpoint": {
      "type": "string"
    },
    "identifier": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "plugin": {
      "type": "string"
    },
    "size_gb": {
      "type": "integer"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	pVDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
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

func (m *PVUnregisterRequest) IsValid() (*gojsonschema.Result, error) {
	return pVUnregisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PVUnregisterRequest) IsRequest() {}

func (m *PVRegisterRequest) IsValid() (*gojsonschema.Result, error) {
	return pVRegisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PVRegisterRequest) IsRequest() {}

func (m *PVDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return pVDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PVDescribeRequest) IsRequest() {}

func (m *PVDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
