package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var backupUnscheduleRequestSchema *gojsonschema.Schema
var backupScheduleRequestSchema *gojsonschema.Schema
var snapshotRestoreRequestSchema *gojsonschema.Schema
var snapshotListRequestSchema *gojsonschema.Schema
var snapshotDescribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	backupUnscheduleRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 3",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	backupScheduleRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "auth_secret_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "bucket_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "cluster": {
      "type": "string"
    },
    "credential_uid": {
      "type": "string"
    },
    "schedule_cron_expr": {
      "type": "string"
    },
    "snapshot_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 13",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	snapshotRestoreRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "auth_secret_name": {
      "maxLength": 63,
      "pattern": "^[a-z0-9](?:[a-z0-9\\-]{0,61}[a-z0-9])?$",
      "type": "string"
    },
    "cluster": {
      "type": "string"
    },
    "snapshot_phid": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 19",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	snapshotListRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 3",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
	snapshotDescribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "cluster": {
      "type": "string"
    },
    "phid": {
      "type": "string"
    },
    "uid": {
      "type": "string"
    }
  },
  "title": "Next Id: 4",
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *BackupUnscheduleRequest) IsValid() (*gojsonschema.Result, error) {
	return backupUnscheduleRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *BackupUnscheduleRequest) IsRequest() {}

func (m *BackupScheduleRequest) IsValid() (*gojsonschema.Result, error) {
	return backupScheduleRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *BackupScheduleRequest) IsRequest() {}

func (m *SnapshotRestoreRequest) IsValid() (*gojsonschema.Result, error) {
	return snapshotRestoreRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SnapshotRestoreRequest) IsRequest() {}

func (m *SnapshotListRequest) IsValid() (*gojsonschema.Result, error) {
	return snapshotListRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SnapshotListRequest) IsRequest() {}

func (m *SnapshotDescribeRequest) IsValid() (*gojsonschema.Result, error) {
	return snapshotDescribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SnapshotDescribeRequest) IsRequest() {}

func (m *SnapshotDescribeResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *SnapshotListResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
