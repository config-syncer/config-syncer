package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var agentCreateRequestSchema *gojsonschema.Schema
var agentDeleteRequestSchema *gojsonschema.Schema
var agentDescribeRequestSchema *gojsonschema.Schema
var agentRestartRequestSchema *gojsonschema.Schema
var agentListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	agentCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1PortInfo": {
      "properties": {
        "port_range": {
          "type": "string"
        },
        "protocol": {
          "type": "string"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "ci_starter_version": {
      "type": "string"
    },
    "executors": {
      "type": "integer"
    },
    "labels": {
      "type": "string"
    },
    "ports": {
      "items": {
        "$ref": "#/definitions/v1beta1PortInfo"
      },
      "type": "array"
    },
    "role": {
      "type": "string"
    },
    "saltbase_version": {
      "type": "string"
    },
    "sku": {
      "type": "string"
    },
    "user_startup_script": {
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
	agentDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	agentRestartRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	agentListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "status": {
      "items": {
        "type": "string"
      },
      "title": "List of status to get the agent filterd on the status\nvalues in\n  PENDING\n  FAILED\n  ONLINE\n  OFFLINE\n  DELETED",
      "type": "array"
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

func (m *AgentDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return agentDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AgentDescribeRequest) IsRequest() {}

func (m *AgentRestartRequest) IsValid() (*gojsonschema.Result, error) {
	return agentRestartRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AgentRestartRequest) IsRequest() {}

func (m *AgentListRequest) IsValid() (*gojsonschema.Result, error) {
	return agentListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AgentListRequest) IsRequest() {}

func (m *AgentRestartResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *AgentListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *AgentDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
