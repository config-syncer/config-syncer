package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var diskDescribeRequestSchema *gojsonschema.Schema
var diskDeleteRequestSchema *gojsonschema.Schema
var diskListRequestSchema *gojsonschema.Schema
var diskCreateRequestSchema *gojsonschema.Schema

func init() {
	var err error
	diskDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
    "plugin": {
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
}

func (m *DiskDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return diskDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskDescribeRequest) IsRequest() {}

func (m *DiskDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return diskDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskDeleteRequest) IsRequest() {}

func (m *DiskListRequest) IsValid() (*gojsonschema.Result, error) {
	return diskListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskListRequest) IsRequest() {}

func (m *DiskCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return diskCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DiskCreateRequest) IsRequest() {}

func (m *DiskDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DiskListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
