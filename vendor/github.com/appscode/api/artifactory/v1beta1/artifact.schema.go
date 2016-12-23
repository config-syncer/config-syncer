package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var artifactSearchRequestSchema *gojsonschema.Schema
var artifactListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	artifactSearchRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "query": {
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
	artifactListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
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

func (m *ArtifactSearchRequest) IsValid() (*gojsonschema.Result, error) {
	return artifactSearchRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ArtifactSearchRequest) IsRequest() {}

func (m *ArtifactListRequest) IsValid() (*gojsonschema.Result, error) {
	return artifactListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ArtifactListRequest) IsRequest() {}

func (m *ArtifactSearchResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ArtifactListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
