package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var projectListRequestSchema *gojsonschema.Schema
var projectMemberListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	projectListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "members": {
      "items": {
        "type": "string"
      },
      "type": "array"
    },
    "with_member": {
      "type": "boolean"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	projectMemberListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "uid": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *ProjectListRequest) IsValid() (*gojsonschema.Result, error) {
	return projectListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ProjectListRequest) IsRequest() {}

func (m *ProjectMemberListRequest) IsValid() (*gojsonschema.Result, error) {
	return projectMemberListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ProjectMemberListRequest) IsRequest() {}

