package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var volumeListRequestSchema *gojsonschema.Schema

func init() {
	var err error
	volumeListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "glusterfs_cluster": {
      "type": "string"
    },
    "kube_cluster": {
      "type": "string"
    },
    "kube_namespace": {
      "type": "string"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *VolumeListRequest) IsValid() (*gojsonschema.Result, error) {
	return volumeListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *VolumeListRequest) IsRequest() {}

func (m *VolumeListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
