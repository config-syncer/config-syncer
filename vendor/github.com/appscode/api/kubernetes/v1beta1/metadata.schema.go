package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var zoneListRequestSchema *gojsonschema.Schema
var bucketListRequestSchema *gojsonschema.Schema
var regionListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	zoneListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cloud_credential": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "region": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	bucketListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cloud_credential": {
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
	regionListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cloud_credential": {
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
}

func (m *ZoneListRequest) IsValid() (*gojsonschema.Result, error) {
	return zoneListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ZoneListRequest) IsRequest() {}

func (m *BucketListRequest) IsValid() (*gojsonschema.Result, error) {
	return bucketListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *BucketListRequest) IsRequest() {}

func (m *RegionListRequest) IsValid() (*gojsonschema.Result, error) {
	return regionListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *RegionListRequest) IsRequest() {}

func (m *RegionListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *BucketListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *ZoneListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
