package dtypes

import (
	"errors"
	"strconv"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"github.com/xeipuuv/gojsonschema"
)

// Helper methods for status object.
const (
	StatusCodeOK string = "0"
)

func (s *Status) StatusCode() StatusCode {
	code, err := strconv.Atoi(s.Code)
	if err != nil {
		code = int(StatusCode_BADREQUEST)
	}
	return StatusCode(code)
}

// returns the status code string of the response status.
func (s *Status) StatusCodeString() string {
	code, err := strconv.Atoi(s.Code)
	if err != nil {
		code = int(StatusCode_BADREQUEST)
	}
	return proto.EnumName(StatusCode_name, int32(code))
}

func (s *Status) IsOK() bool {
	if s.Code == StatusCodeOK {
		return true
	}
	return false
}

func (s *Status) IsError() bool {
	if s.Code != StatusCodeOK {
		return true
	}
	return false
}

func (s *Status) Error() error {
	return errors.New(s.Status)
}

func (s *Status) StatusString() string {
	return s.Status
}

// Adds any proto message in the details field of the Status message.
// This uses google.protobuf.any to to hold and retried data.
func (a *Status) AddDetails(v ...proto.Message) {
	if len(a.Details) == 0 {
		a.Details = make([]*google_protobuf.Any, 0)
	}
	for _, val := range v {
		value, err := proto.Marshal(val)
		if err != nil {
			glog.V(1).Infoln("Marshaling any failed.")
			continue
		}
		anyValue := &google_protobuf.Any{
			TypeUrl: proto.MessageName(val),
			Value:   value,
		}
		a.Details = append(a.Details, anyValue)
	}
}

// Ideally schema.py should generate the functions below, but it can't do it today.
// So, this is a manually written for now.
func (m *VoidRequest) IsValid() (*gojsonschema.Result, error) {
	return &gojsonschema.Result{}, nil
}
func (m *VoidRequest) IsRequest() {}

func (m *LongRunningResponse) SetStatus(s *Status) {
	m.Status = s
}

func (m *VoidResponse) SetStatus(s *Status) {
	m.Status = s
}
