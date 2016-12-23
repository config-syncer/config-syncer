package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var sSHGetRequestSchema *gojsonschema.Schema

func init() {
	var err error
	sSHGetRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "instance_name": {
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    }
  },
  "title": "Use specific requests for protos",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *SSHGetRequest) IsValid() (*gojsonschema.Result, error) {
	return sSHGetRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SSHGetRequest) IsRequest() {}

func (m *SSHGetResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
