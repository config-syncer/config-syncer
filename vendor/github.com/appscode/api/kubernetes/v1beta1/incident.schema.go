package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var incidentNotifyRequestSchema *gojsonschema.Schema
var incidentEventCreateRequestSchema *gojsonschema.Schema
var incidentListRequestSchema *gojsonschema.Schema
var incidentDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	incidentNotifyRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "alert_phid": {
      "type": "string"
    },
    "author": {
      "type": "string"
    },
    "comment": {
      "type": "string"
    },
    "host_name": {
      "type": "string"
    },
    "kubernetes_alert_name": {
      "type": "string"
    },
    "kubernetes_cluster": {
      "type": "string"
    },
    "output": {
      "type": "string"
    },
    "state": {
      "type": "string"
    },
    "time": {
      "description": "The time object is used in icinga to send request. This\nindicates detection time from icinga.",
      "type": "integer"
    },
    "type": {
      "type": "string"
    }
  },
  "title": "Next Id: 12",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	incidentEventCreateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "acknowledge": {
      "type": "boolean"
    },
    "comment": {
      "type": "string"
    },
    "phid": {
      "title": "Incident PHID",
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	incidentListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "kubernetes_cluster": {
      "type": "string"
    },
    "kubernetes_namespace": {
      "type": "string"
    },
    "kubernetes_object_name": {
      "type": "string"
    },
    "kubernetes_object_type": {
      "type": "string"
    },
    "states": {
      "items": {
        "type": "string"
      },
      "type": "array"
    }
  },
  "title": "Next Id: 6",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	incidentDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "phid": {
      "type": "string"
    }
  },
  "title": "Next Id: 2",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *IncidentNotifyRequest) IsValid() (*gojsonschema.Result, error) {
	return incidentNotifyRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IncidentNotifyRequest) IsRequest() {}

func (m *IncidentEventCreateRequest) IsValid() (*gojsonschema.Result, error) {
	return incidentEventCreateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IncidentEventCreateRequest) IsRequest() {}

func (m *IncidentListRequest) IsValid() (*gojsonschema.Result, error) {
	return incidentListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IncidentListRequest) IsRequest() {}

func (m *IncidentDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return incidentDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IncidentDescribeRequest) IsRequest() {}

