package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var purchaseBeginRequestSchema *gojsonschema.Schema
var purchaseCompleteRequestSchema *gojsonschema.Schema
var purchaseCloseRequestSchema *gojsonschema.Schema

func init() {
	var err error
	purchaseBeginRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "count": {
      "type": "integer"
    },
    "product_sku": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	purchaseCompleteRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "failed": {
      "type": "boolean"
    },
    "metadata": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "object_phid": {
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
	purchaseCloseRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "object_phid": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *PurchaseBeginRequest) IsValid() (*gojsonschema.Result, error) {
	return purchaseBeginRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PurchaseBeginRequest) IsRequest() {}

func (m *PurchaseCompleteRequest) IsValid() (*gojsonschema.Result, error) {
	return purchaseCompleteRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PurchaseCompleteRequest) IsRequest() {}

func (m *PurchaseCloseRequest) IsValid() (*gojsonschema.Result, error) {
	return purchaseCloseRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *PurchaseCloseRequest) IsRequest() {}

func (m *PurchaseBeginResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
