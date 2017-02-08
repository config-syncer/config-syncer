package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var clusterScaleRequestSchema *gojsonschema.Schema
var clusterInstanceListRequestSchema *gojsonschema.Schema
var clusterInstanceByIPRequestSchema *gojsonschema.Schema
var clusterUpgradeRequestSchema *gojsonschema.Schema
var clusterDeleteRequestSchema *gojsonschema.Schema
var clusterCreateRequestSchema *gojsonschema.Schema
var clusterUpdateRequestSchema *gojsonschema.Schema
var clusterDescribeRequestSchema *gojsonschema.Schema
var clusterListRequestSchema *gojsonschema.Schema
var clusterStartupConfigRequestSchema *gojsonschema.Schema
var clusterClientConfigRequestSchema *gojsonschema.Schema

func init() {
	var err error
	clusterScaleRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1InstanceGroup": {
      "properties": {
        "count": {
          "type": "integer"
        },
        "sku": {
          "type": "string"
        },
        "use_spot_instances": {
          "type": "boolean"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "node_groups": {
      "items": {
        "$ref": "#/definitions/v1beta1InstanceGroup"
      },
      "type": "array"
    },
    "node_set": {
      "additionalProperties": {
        "type": "integer"
      },
      "title": "New node configuration for the cluster",
      "type": "object"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clusterInstanceListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster_name": {
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
	clusterInstanceByIPRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "external_ip": {
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
	clusterUpgradeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "hostfacts_version": {
      "type": "string"
    },
    "kube_saltbase_version": {
      "type": "string"
    },
    "kube_server_version": {
      "type": "string"
    },
    "kube_starter_version": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "version": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clusterDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "delete_dynamic_volumes": {
      "type": "boolean"
    },
    "force": {
      "type": "boolean"
    },
    "keep_lodabalancers": {
      "type": "boolean"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "release_reserved_ip": {
      "type": "boolean"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clusterCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1InstanceGroup": {
      "properties": {
        "count": {
          "type": "integer"
        },
        "sku": {
          "type": "string"
        },
        "use_spot_instances": {
          "type": "boolean"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "cloud_credential": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "cloud_credential_data": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "default_access_level": {
      "title": "Default access level is to allow permission to the cluster\nwhen no Role matched for that specif user or group. This can\nset as\n   - team-admins       // to allow ns admin access\n   - cluster-admins    // to allow admin access\n   - cluster-editors   // to allow editor access\n   - cluster-deployers // to allow deployer access\n   - cluster-viewers   // to allow viewer access\n   - no-access         // to allow no default access\nIf not set this will set \"\"",
      "type": "string"
    },
    "do_not_delete": {
      "type": "boolean"
    },
    "gce_project": {
      "type": "string"
    },
    "hostfacts_version": {
      "type": "string"
    },
    "kube_saltbase_version": {
      "type": "string"
    },
    "kube_server_version": {
      "type": "string"
    },
    "kube_starter_version": {
      "type": "string"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "node_groups": {
      "items": {
        "$ref": "#/definitions/v1beta1InstanceGroup"
      },
      "type": "array"
    },
    "node_set": {
      "additionalProperties": {
        "type": "integer"
      },
      "type": "object"
    },
    "provider": {
      "type": "string"
    },
    "version": {
      "type": "string"
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
	clusterUpdateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1ClusterSettings": {
      "properties": {
        "log_index_prefix": {
          "type": "string"
        },
        "log_storage_lifetime": {
          "title": "Number of secs logs will be stored in ElasticSearch",
          "type": "integer"
        },
        "monitoring_storage_lifetime": {
          "title": "Number of secs logs will be stored in InfluxDB",
          "type": "integer"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "default_access_level": {
      "title": "Default access level is to allow permission to the cluster\nwhen no Role matched for that specif user or group. This can\nset as\n   - v:cluster-admins    // to allow admin access\n   - v:cluster-deployer  // to allow deployer access\n   - v:cluster-viewer    // to allow viewer access\n   - \"\"                  // empty value stands for no access",
      "type": "string"
    },
    "do_not_delete": {
      "type": "boolean"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "settings": {
      "$ref": "#/definitions/v1beta1ClusterSettings"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clusterDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	clusterListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "status": {
      "items": {
        "type": "string"
      },
      "type": "array"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	clusterStartupConfigRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "context_version": {
      "type": "integer"
    },
    "role": {
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
	clusterClientConfigRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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

func (m *ClusterScaleRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterScaleRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterScaleRequest) IsRequest() {}

func (m *ClusterInstanceListRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterInstanceListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterInstanceListRequest) IsRequest() {}

func (m *ClusterInstanceByIPRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterInstanceByIPRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterInstanceByIPRequest) IsRequest() {}

func (m *ClusterUpgradeRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterUpgradeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterUpgradeRequest) IsRequest() {}

func (m *ClusterDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterDeleteRequest) IsRequest() {}

func (m *ClusterCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterCreateRequest) IsRequest() {}

func (m *ClusterUpdateRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterUpdateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterUpdateRequest) IsRequest() {}

func (m *ClusterDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterDescribeRequest) IsRequest() {}

func (m *ClusterListRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterListRequest) IsRequest() {}

func (m *ClusterStartupConfigRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterStartupConfigRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterStartupConfigRequest) IsRequest() {}

func (m *ClusterClientConfigRequest) IsValid() (*gojsonschema.Result, error) {
	return clusterClientConfigRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClusterClientConfigRequest) IsRequest() {}

func (m *ClusterStartupConfigResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ClusterDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ClusterListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ClusterClientConfigResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ClusterInstanceListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ClusterInstanceResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
