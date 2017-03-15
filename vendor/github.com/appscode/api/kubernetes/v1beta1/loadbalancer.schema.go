package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var listRequestSchema *gojsonschema.Schema
var deleteRequestSchema *gojsonschema.Schema
var createRequestSchema *gojsonschema.Schema
var describeRequestSchema *gojsonschema.Schema
var updateRequestSchema *gojsonschema.Schema

func init() {
	var err error
	listRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
	deleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "kind": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "namespace": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	createRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "kubernetesv1beta1Status": {
      "properties": {
        "status": {
          "items": {
            "$ref": "#/definitions/v1beta1LoadBalancerStatus"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1HTTPLoadBalancerRule": {
      "properties": {
        "backend": {
          "$ref": "#/definitions/v1beta1LoadBalancerBackend"
        },
        "header_rules": {
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "path": {
          "type": "string"
        },
        "rewrite_rules": {
          "items": {
            "type": "string"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancer": {
      "properties": {
        "creation_timestamp": {
          "type": "integer"
        },
        "json": {
          "type": "string"
        },
        "kind": {
          "description": "'kind' defines is it the regular kubernetes instance or the\nappscode superset called Extended Ingress. This field will\nstrictly contains only those two values\n'ingress' - default kubernetes ingress object.\n'extendedIngress' - appscode superset of ingress.\nwhen creating a Loadbalancer from UI this field will always\nbe only 'extendedIngress.' List, Describe, Update and Delete\nwill support both two modes.\nCreate will support only extendedIngress.\nFor Creating or Updating an regular ingress one must use the\nkubectl or direct API calls directly to kubernetes.",
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "namespace": {
          "type": "string"
        },
        "options": {
          "additionalProperties": {
            "type": "string"
          },
          "type": "object"
        },
        "spec": {
          "$ref": "#/definitions/v1beta1Spec"
        },
        "status": {
          "$ref": "#/definitions/kubernetesv1beta1Status"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerBackend": {
      "properties": {
        "service_name": {
          "type": "string"
        },
        "service_port": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerRule": {
      "properties": {
        "SSL_secret_name": {
          "description": "ssl secret name to enable https on the host.\nssl secret must contain data with the certs pem file.",
          "type": "string"
        },
        "host": {
          "type": "string"
        },
        "http": {
          "items": {
            "$ref": "#/definitions/v1beta1HTTPLoadBalancerRule"
          },
          "type": "array"
        },
        "tcp": {
          "items": {
            "$ref": "#/definitions/v1beta1TCPLoadBalancerRule"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerStatus": {
      "properties": {
        "IP": {
          "type": "string"
        },
        "host": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "v1beta1Spec": {
      "properties": {
        "backend": {
          "$ref": "#/definitions/v1beta1HTTPLoadBalancerRule"
        },
        "rules": {
          "items": {
            "$ref": "#/definitions/v1beta1LoadBalancerRule"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1TCPLoadBalancerRule": {
      "properties": {
        "SSL_secret_name": {
          "type": "string"
        },
        "backend": {
          "$ref": "#/definitions/v1beta1LoadBalancerBackend"
        },
        "port": {
          "type": "string"
        },
        "secret_pem_name": {
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
    "load_balancer": {
      "$ref": "#/definitions/v1beta1LoadBalancer"
    },
    "name": {
      "type": "string"
    },
    "namespace": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	describeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "kind": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "namespace": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	updateRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "kubernetesv1beta1Status": {
      "properties": {
        "status": {
          "items": {
            "$ref": "#/definitions/v1beta1LoadBalancerStatus"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1HTTPLoadBalancerRule": {
      "properties": {
        "backend": {
          "$ref": "#/definitions/v1beta1LoadBalancerBackend"
        },
        "header_rules": {
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "path": {
          "type": "string"
        },
        "rewrite_rules": {
          "items": {
            "type": "string"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancer": {
      "properties": {
        "creation_timestamp": {
          "type": "integer"
        },
        "json": {
          "type": "string"
        },
        "kind": {
          "description": "'kind' defines is it the regular kubernetes instance or the\nappscode superset called Extended Ingress. This field will\nstrictly contains only those two values\n'ingress' - default kubernetes ingress object.\n'extendedIngress' - appscode superset of ingress.\nwhen creating a Loadbalancer from UI this field will always\nbe only 'extendedIngress.' List, Describe, Update and Delete\nwill support both two modes.\nCreate will support only extendedIngress.\nFor Creating or Updating an regular ingress one must use the\nkubectl or direct API calls directly to kubernetes.",
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "namespace": {
          "type": "string"
        },
        "options": {
          "additionalProperties": {
            "type": "string"
          },
          "type": "object"
        },
        "spec": {
          "$ref": "#/definitions/v1beta1Spec"
        },
        "status": {
          "$ref": "#/definitions/kubernetesv1beta1Status"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerBackend": {
      "properties": {
        "service_name": {
          "type": "string"
        },
        "service_port": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerRule": {
      "properties": {
        "SSL_secret_name": {
          "description": "ssl secret name to enable https on the host.\nssl secret must contain data with the certs pem file.",
          "type": "string"
        },
        "host": {
          "type": "string"
        },
        "http": {
          "items": {
            "$ref": "#/definitions/v1beta1HTTPLoadBalancerRule"
          },
          "type": "array"
        },
        "tcp": {
          "items": {
            "$ref": "#/definitions/v1beta1TCPLoadBalancerRule"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1LoadBalancerStatus": {
      "properties": {
        "IP": {
          "type": "string"
        },
        "host": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "v1beta1Spec": {
      "properties": {
        "backend": {
          "$ref": "#/definitions/v1beta1HTTPLoadBalancerRule"
        },
        "rules": {
          "items": {
            "$ref": "#/definitions/v1beta1LoadBalancerRule"
          },
          "type": "array"
        }
      },
      "type": "object"
    },
    "v1beta1TCPLoadBalancerRule": {
      "properties": {
        "SSL_secret_name": {
          "type": "string"
        },
        "backend": {
          "$ref": "#/definitions/v1beta1LoadBalancerBackend"
        },
        "port": {
          "type": "string"
        },
        "secret_pem_name": {
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
    "load_balancer": {
      "$ref": "#/definitions/v1beta1LoadBalancer"
    },
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

func (m *ListRequest) IsValid() (*gojsonschema.Result, error) {
	return listRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ListRequest) IsRequest() {}

func (m *DeleteRequest) IsValid() (*gojsonschema.Result, error) {
	return deleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DeleteRequest) IsRequest() {}

func (m *CreateRequest) IsValid() (*gojsonschema.Result, error) {
	return createRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CreateRequest) IsRequest() {}

func (m *DescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return describeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *DescribeRequest) IsRequest() {}

func (m *UpdateRequest) IsValid() (*gojsonschema.Result, error) {
	return updateRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *UpdateRequest) IsRequest() {}

func (m *ListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *DescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
