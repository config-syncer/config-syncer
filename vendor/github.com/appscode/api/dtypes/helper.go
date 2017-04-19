package dtypes

import (
	"strconv"
	"strings"

	"github.com/appscode/errors"
	_env "github.com/appscode/go/env"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
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

func NewStatusOK(message ...string) *Status {
	glog.V(4).Infoln("Sending OK response with message ", message)
	return &Status{
		Code:    statusCodeString(int32(StatusCode_OK)),
		Status:  StatusCode_OK.String(),
		Message: strings.Join(message, ";"),
	}
}

func NewStatusFromError(err error) *Status {
	e := errors.Parse(err)
	if e == nil {
		// not an error, so sending ok response
		return NewStatusOK()
	}
	glog.V(4).Infoln("Sending response ", e.Code(), " with message ", e.Message())
	s := &Status{
		Code:    statusCodeString(errorToAPICode(e.Code())),
		Status:  e.Code(),
		Message: statusMessage(e.Code(), e.Messages()),
	}

	if errorHelp := e.Help(); errorHelp != nil {
		s.Help = &Help{
			Url:         errorHelp.Url,
			Description: errorHelp.Description,
		}
	}
	if e.Trace() != nil {
		if !_env.FromHost().IsPublic() {
			s.AddDetails(&ErrorDetails{
				RequestedResource: e.Error(),
				Stacktrace:        e.TraceString(),
			})
		}
	}
	glog.V(4).Infoln("Sending EROR response with message", e.Messages())
	return s
}

func NewStatusUnauthorized() *Status {
	err := errors.New().Unauthorized()
	return NewStatusFromError(err)
}

func NewStatusBadRequest(message ...string) *Status {
	err := errors.New().WithMessage(message...).BadRequest()
	return NewStatusFromError(err)
}

func statusCodeString(code int32) string {
	return strconv.FormatInt(int64(code), 10)
}

func errorToAPICode(code string) int32 {
	c, ok := StatusCode_value[code]
	if ok {
		return c
	}
	return -1
}

func statusMessage(code string, msg []string) string {
	switch code {
	case errors.Internal:
		return "Internal Error. Contact support, support@appscode.com"
	case errors.Unauthorized:
		return "Unauthorized. Please Login."
	}
	return strings.Join(msg, ";")
}

// Adds any proto message in the details field of the Status message.
// This uses google.protobuf.any to to hold and retried data.
func (a *Status) AddDetails(v ...proto.Message) {
	if len(a.Details) == 0 {
		a.Details = make([]*any.Any, 0)
	}
	for _, val := range v {
		value, err := proto.Marshal(val)
		if err != nil {
			glog.V(1).Infoln("Marshaling any failed.")
			continue
		}
		anyValue := &any.Any{
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
