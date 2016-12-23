package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var jobListRequestSchema *gojsonschema.Schema
var jobCopyRequestSchema *gojsonschema.Schema
var jobDescribeRequestSchema *gojsonschema.Schema
var jobCreateRequestSchema *gojsonschema.Schema
var jobBuildRequestSchema *gojsonschema.Schema
var jobDeleteRequestSchema *gojsonschema.Schema

func init() {
	var err error
	jobListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	jobCopyRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "destination": {
      "type": "string"
    },
    "source": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	jobDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	jobCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "parents": {
      "type": "string"
    },
    "sh_file": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	jobBuildRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "param": {
      "type": "string"
    },
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	jobDeleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "parents": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *JobListRequest) IsValid() (*gojsonschema.Result, error) {
	return jobListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobListRequest) IsRequest() {}

func (m *JobCopyRequest) IsValid() (*gojsonschema.Result, error) {
	return jobCopyRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobCopyRequest) IsRequest() {}

func (m *JobDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return jobDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobDescribeRequest) IsRequest() {}

func (m *JobCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return jobCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobCreateRequest) IsRequest() {}

func (m *JobBuildRequest) IsValid() (*gojsonschema.Result, error) {
	return jobBuildRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobBuildRequest) IsRequest() {}

func (m *JobDeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return jobDeleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *JobDeleteRequest) IsRequest() {}

func (m *JobDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *JobListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
