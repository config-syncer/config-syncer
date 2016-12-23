package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var versionDescribeRequestSchema *gojsonschema.Schema
var versionListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	versionDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "type": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	versionListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "type": "string"
    },
    "type": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *VersionDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return versionDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *VersionDescribeRequest) IsRequest() {}

func (m *VersionListRequest) IsValid() (*gojsonschema.Result, error) {
	return versionListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *VersionListRequest) IsRequest() {}

func (m *VersionListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *VersionDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
