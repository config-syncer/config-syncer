package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var clientReconfigureRequestSchema *gojsonschema.Schema

func init() {
	var err error
	clientReconfigureRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
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
}

func (m *ClientReconfigureRequest) IsValid() (*gojsonschema.Result, error) {
	return clientReconfigureRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *ClientReconfigureRequest) IsRequest() {}

