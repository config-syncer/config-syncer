package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var databaseDeleteRequestSchema *gojsonschema.Schema
var databaseUpdateRequestSchema *gojsonschema.Schema
var databaseCreateRequestSchema *gojsonschema.Schema
var databaseListRequestSchema *gojsonschema.Schema
var databaseScaleRequestSchema *gojsonschema.Schema
var databaseRecoverRequestSchema *gojsonschema.Schema
var databaseDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	databaseDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "destroy": {
      "type": "boolean"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseUpdateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "do_not_delete": {
      "type": "boolean"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "DatabaseCreateRequestInitialScript": {
      "properties": {
        "config_map": {
          "$ref": "#/definitions/InitialScriptConfigMap"
        },
        "git_repo": {
          "$ref": "#/definitions/InitialScriptGitRepo"
        },
        "script_path": {
          "type": "string"
        },
        "secret": {
          "$ref": "#/definitions/InitialScriptSecret"
        }
      },
      "type": "object"
    },
    "InitialScriptConfigMap": {
      "properties": {
        "name": {
          "maxLength": 63,
          "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
          "type": "string"
        }
      },
      "type": "object"
    },
    "InitialScriptGitRepo": {
      "properties": {
        "directory": {
          "type": "string"
        },
        "repository": {
          "type": "string"
        },
        "revision": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "InitialScriptSecret": {
      "properties": {
        "name": {
          "maxLength": 63,
          "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
          "type": "string"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "auth_secret_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "cluster": {
      "type": "string"
    },
    "hostname": {
      "type": "string"
    },
    "ignore_validation": {
      "type": "boolean"
    },
    "initial_script": {
      "$ref": "#/definitions/DatabaseCreateRequestInitialScript"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "node_count": {
      "type": "integer"
    },
    "pv_size_gb": {
      "type": "integer"
    },
    "sku": {
      "type": "string"
    },
    "snapshot_phid": {
      "type": "string"
    },
    "storage_class": {
      "type": "string"
    },
    "type": {
      "type": "string"
    },
    "version": {
      "type": "string"
    }
  },
  "title": "Next Id: 19",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "status": {
      "items": {
        "type": "string"
      },
      "title": "List of status to get the agent filterd on the status\nvalues in\n  PENDING\n  FAILED\n  READY\n  DELETING\n  DELETED\n  DESTROYED",
      "type": "array"
    },
    "type": {
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseScaleRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "node_count": {
      "type": "integer"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseRecoverRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 3",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	databaseDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 3",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *DatabaseDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseDeleteRequest) IsRequest() {}

func (m *DatabaseUpdateRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseUpdateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseUpdateRequest) IsRequest() {}

func (m *DatabaseCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseCreateRequest) IsRequest() {}

func (m *DatabaseListRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseListRequest) IsRequest() {}

func (m *DatabaseScaleRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseScaleRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseScaleRequest) IsRequest() {}

func (m *DatabaseRecoverRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseRecoverRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseRecoverRequest) IsRequest() {}

func (m *DatabaseDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return databaseDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DatabaseDescribeRequest) IsRequest() {}

func (m *DatabaseListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DatabaseDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
