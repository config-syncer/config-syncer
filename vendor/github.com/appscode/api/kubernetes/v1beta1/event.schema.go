package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var eventRequestSchema *gojsonschema.Schema

func init() {
	var err error
	eventRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "EventRequestAncestors": {
      "properties": {
        "names": {
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "type": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "EventRequestObjectMeta": {
      "properties": {
        "ancestors": {
          "items": {
            "$ref": "#/definitions/EventRequestAncestors"
          },
          "type": "array"
        },
        "instance_external_ip": {
          "type": "string"
        },
        "instance_id": {
          "type": "string"
        },
        "instance_internal_ip": {
          "type": "string"
        },
        "kind": {
          "type": "string"
        },
        "labels": {
          "additionalProperties": {
            "type": "string"
          },
          "type": "object"
        },
        "namespace": {
          "type": "string"
        },
        "pod_ip": {
          "type": "string"
        }
      },
      "type": "object"
    }
  },
  "properties": {
    "cluster_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "event_type": {
      "type": "string"
    },
    "kube_namespace": {
      "type": "string"
    },
    "kube_object_name": {
      "type": "string"
    },
    "kube_object_type": {
      "type": "string"
    },
    "metadata": {
      "$ref": "#/definitions/EventRequestObjectMeta"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *EventRequest) IsValid() (*gojsonschema.Result, error) {
	return eventRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *EventRequest) IsRequest() {}

func (m *EventResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
