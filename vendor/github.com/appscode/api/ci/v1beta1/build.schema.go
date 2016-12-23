package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var buildListRequestSchema *gojsonschema.Schema
var buildDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	buildListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "job_name": {
      "type": "string"
    },
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	buildDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "job_name": {
      "type": "string"
    },
    "number": {
      "type": "integer"
    },
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *BuildListRequest) IsValid() (*gojsonschema.Result, error) {
	return buildListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *BuildListRequest) IsRequest() {}

func (m *BuildDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return buildDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *BuildDescribeRequest) IsRequest() {}

func (m *BuildDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *BuildListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
