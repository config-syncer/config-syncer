package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var subscriptionSubscribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	subscriptionSubscribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "auto_extend": {
      "type": "boolean"
    },
    "product_id": {
      "type": "string"
    },
    "start_time": {
      "type": "integer"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *SubscriptionSubscribeRequest) IsValid() (*gojsonschema.Result, error) {
	return subscriptionSubscribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SubscriptionSubscribeRequest) IsRequest() {}

