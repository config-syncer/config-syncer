package v1beta2

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var secretEditRequestSchema *gojsonschema.Schema
var persistentVolumeClaimRegisterRequestSchema *gojsonschema.Schema
var diskListRequestSchema *gojsonschema.Schema
var createResourceRequestSchema *gojsonschema.Schema
var updateResourceRequestSchema *gojsonschema.Schema
var diskDescribeRequestSchema *gojsonschema.Schema
var persistentVolumeUnRegisterRequestSchema *gojsonschema.Schema
var copyResourceRequestSchema *gojsonschema.Schema
var describeResourceRequestSchema *gojsonschema.Schema
var configMapEditRequestSchema *gojsonschema.Schema
var listResourceRequestSchema *gojsonschema.Schema
var persistentVolumeClaimUnRegisterRequestSchema *gojsonschema.Schema
var diskCreateRequestSchema *gojsonschema.Schema
var diskDeleteRequestSchema *gojsonschema.Schema
var persistentVolumeRegisterRequestSchema *gojsonschema.Schema
var deleteResourceRequestSchema *gojsonschema.Schema

func init() {
	var err error
	secretEditRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "add": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "cluster": {
      "type": "string"
    },
    "deleted": {
      "items": {
        "type": "string"
      },
      "type": "array"
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
    "update": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	persistentVolumeClaimRegisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
    "size_gb": {
      "type": "integer"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	diskListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	createResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta2Raw": {
      "properties": {
        "data": {
          "type": "string"
        },
        "format": {
          "type": "string"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "raw": {
      "$ref": "#/definitions/v1beta2Raw"
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
	updateResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta2Raw": {
      "properties": {
        "data": {
          "type": "string"
        },
        "format": {
          "type": "string"
        }
      },
      "type": "object"
    }
  },
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
      "$ref": "#/definitions/v1beta2Raw"
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
	diskDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "type": "string"
    },
    "provider": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	persistentVolumeUnRegisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	copyResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta2KubeObject": {
      "properties": {
        "cluster": {
          "type": "string"
        },
        "name": {
          "maxLength": 63,
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
      "$ref": "#/definitions/v1beta2KubeObject"
    },
    "source": {
      "$ref": "#/definitions/v1beta2KubeObject"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	describeResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "include_metrics": {
      "type": "boolean"
    },
    "name": {
      "maxLength": 63,
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "raw": {
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
	configMapEditRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "add": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "cluster": {
      "type": "string"
    },
    "deleted": {
      "items": {
        "type": "string"
      },
      "type": "array"
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
    "update": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	listResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "ListResourceRequestAncestor": {
      "properties": {
        "name": {
          "maxLength": 63,
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
    "ancestor": {
      "$ref": "#/definitions/ListResourceRequestAncestor"
    },
    "cluster": {
      "type": "string"
    },
    "include_metrics": {
      "type": "boolean"
    },
    "label_selector": {
      "additionalProperties": {
        "type": "string"
      }
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
}`))
	if err != nil {
		glog.Fatal(err)
	}
	persistentVolumeClaimUnRegisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	diskCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "disk_type": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "size_gb": {
      "type": "integer"
    },
    "zone": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	diskDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	persistentVolumeRegisterRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	deleteResourceRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "type": "string"
    },
    "namespace": {
      "maxLength": 63,
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

func (m *SecretEditRequest) IsValid() (*gojsonschema.Result, error) {
	return secretEditRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SecretEditRequest) IsRequest() {}

func (m *PersistentVolumeClaimRegisterRequest) IsValid() (*gojsonschema.Result, error) {
	return persistentVolumeClaimRegisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PersistentVolumeClaimRegisterRequest) IsRequest() {}

func (m *DiskListRequest) IsValid() (*gojsonschema.Result, error) {
	return diskListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskListRequest) IsRequest() {}

func (m *CreateResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return createResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CreateResourceRequest) IsRequest() {}

func (m *UpdateResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return updateResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *UpdateResourceRequest) IsRequest() {}

func (m *DiskDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return diskDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskDescribeRequest) IsRequest() {}

func (m *PersistentVolumeUnRegisterRequest) IsValid() (*gojsonschema.Result, error) {
	return persistentVolumeUnRegisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PersistentVolumeUnRegisterRequest) IsRequest() {}

func (m *CopyResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return copyResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CopyResourceRequest) IsRequest() {}

func (m *DescribeResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return describeResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DescribeResourceRequest) IsRequest() {}

func (m *ConfigMapEditRequest) IsValid() (*gojsonschema.Result, error) {
	return configMapEditRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ConfigMapEditRequest) IsRequest() {}

func (m *ListResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return listResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ListResourceRequest) IsRequest() {}

func (m *PersistentVolumeClaimUnRegisterRequest) IsValid() (*gojsonschema.Result, error) {
	return persistentVolumeClaimUnRegisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PersistentVolumeClaimUnRegisterRequest) IsRequest() {}

func (m *DiskCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return diskCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskCreateRequest) IsRequest() {}

func (m *DiskDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return diskDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskDeleteRequest) IsRequest() {}

func (m *PersistentVolumeRegisterRequest) IsValid() (*gojsonschema.Result, error) {
	return persistentVolumeRegisterRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PersistentVolumeRegisterRequest) IsRequest() {}

func (m *DeleteResourceRequest) IsValid() (*gojsonschema.Result, error) {
	return deleteResourceRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DeleteResourceRequest) IsRequest() {}

func (m *DiskDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DescribeResourceResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DiskListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ListResourceResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
