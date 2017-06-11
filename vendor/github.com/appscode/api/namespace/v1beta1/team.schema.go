package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var createRequestSchema *gojsonschema.Schema
var getRequestSchema *gojsonschema.Schema
var isAvailableRequestSchema *gojsonschema.Schema

func init() {
	var err error
	createRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "v1beta1Address": {
      "properties": {
        "company": {
          "type": "string"
        },
        "country_code_numeric": {
          "title": "Ref https://developers.braintreepayments.com/reference/general/countries/ruby",
          "type": "string"
        },
        "extended_address": {
          "type": "string"
        },
        "first_name": {
          "type": "string"
        },
        "last_name": {
          "type": "string"
        },
        "locality": {
          "type": "string"
        },
        "postal_code": {
          "type": "string"
        },
        "region": {
          "type": "string"
        },
        "street_address": {
          "type": "string"
        }
      },
      "title": "Example\nresult = Braintree::Address.create(\n  :first_name          => 'Jenna',\n  :last_name           => 'Smith',\n  :company             => 'Braintree',\n  :street_address      => '1 E Main St',\n  :extended_address    => 'Suite 403',\n  :locality            => 'Chicago',\n  :region              => 'Illinois',\n  :postal_code         => '60622',\n  :country_code_numeric => '840'\n)",
      "type": "object"
    }
  },
  "properties": {
    "billing_address": {
      "$ref": "#/definitions/v1beta1Address"
    },
    "display_name": {
      "type": "string"
    },
    "email": {
      "type": "string"
    },
    "initial_units": {
      "type": "integer"
    },
    "invite_emails": {
      "items": {
        "type": "string"
      },
      "type": "array"
    },
    "name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "password": {
      "type": "string"
    },
    "payment_method_nonce": {
      "type": "string"
    },
    "subscription": {
      "type": "string"
    },
    "user_name": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	getRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
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
	isAvailableRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "name": {
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

func (m *CreateRequest) IsValid() (*gojsonschema.Result, error) {
	return createRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *CreateRequest) IsRequest() {}

func (m *GetRequest) IsValid() (*gojsonschema.Result, error) {
	return getRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *GetRequest) IsRequest() {}

func (m *IsAvailableRequest) IsValid() (*gojsonschema.Result, error) {
	return isAvailableRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *IsAvailableRequest) IsRequest() {}

