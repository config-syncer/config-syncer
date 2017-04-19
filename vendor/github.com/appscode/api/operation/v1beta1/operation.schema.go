package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var describeRequestSchema *gojsonschema.Schema
var logDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	describeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "phid": {
      "type": "string"
    },
    "timestamp": {
      "type": "integer"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	logDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "log_id": {
      "type": "string"
    },
    "phid": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *DescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return describeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DescribeRequest) IsRequest() {}

func (m *LogDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return logDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *LogDescribeRequest) IsRequest() {}

func (m *LogDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
