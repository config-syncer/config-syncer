package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var agentCreateRequestSchema *gojsonschema.Schema
var agentDeleteRequestSchema *gojsonschema.Schema

func init() {
	var err error
	agentCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "external_ip": {
      "type": "string"
    },
    "git_ssh_public_key": {
      "type": "string"
    },
    "internal_ip": {
      "type": "string"
    },
    "jenkins_jnlp_port": {
      "type": "integer"
    },
    "name": {
      "type": "string"
    },
    "role": {
      "type": "string"
    },
    "ssh_port": {
      "type": "integer"
    },
    "ssh_user": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	agentDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *AgentCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return agentCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AgentCreateRequest) IsRequest() {}

func (m *AgentDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return agentDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AgentDeleteRequest) IsRequest() {}

func (m *AgentListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *AgentCreateResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
