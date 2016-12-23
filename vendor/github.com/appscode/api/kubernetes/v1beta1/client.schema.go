package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var appDescribeRequestSchema *gojsonschema.Schema
var configMapDescribeRequestSchema *gojsonschema.Schema
var secretDescribeRequestSchema *gojsonschema.Schema
var copyResourceRequestSchema *gojsonschema.Schema
var clientRequestSchema *gojsonschema.Schema

func init() {
	var err error
	appDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
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
	configMapDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "raw": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	secretDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "raw": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	copyResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1KubeResource": {
      "properties": {
        "cluster": {
          "type": "string"
        },
        "name": {
          "maxLength": 63,
          "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
          "type": "string"
        },
        "namespace": {
          "maxLength": 63,
          "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "destination": {
      "$ref": "#/definitions/v1beta1KubeResource"
    },
    "source": {
      "$ref": "#/definitions/v1beta1KubeResource"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clientRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *AppDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return appDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *AppDescribeRequest) IsRequest() {}

func (m *ConfigMapDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return configMapDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ConfigMapDescribeRequest) IsRequest() {}

func (m *SecretDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return secretDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SecretDescribeRequest) IsRequest() {}

func (m *CopyResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return copyResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CopyResourceRequest) IsRequest() {}

func (m *ClientRequest) IsValid() (*gojsonschema.Result, error) {
	return clientRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClientRequest) IsRequest() {}

func (m *ConfigMapListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *SecretListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ConfigMapDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *NamespaceListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ServiceListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *AppListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *AppDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *NodeListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *SecretDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *CopyResourceResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *JobListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ReplicationControllerListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *PodListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
