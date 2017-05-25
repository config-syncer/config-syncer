package dtypes

import (
	"strings"

	"github.com/appscode/errors"
	_env "github.com/appscode/go/env"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/xeipuuv/gojsonschema"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func statusErr(c codes.Code, err error) error {
	// if already a statusError, just return it (ignore c)
	if gs, ok := status.FromError(err); ok {
		return gs.Err()
	}
	e := errors.FromErr(err)
	// if the cause of traceable error is a statusError, just return it (ignore c)
	if e2 := e.Cause(); e2 != nil {
		if gs, ok := status.FromError(e2); ok {
			return gs.Err()
		}
	}
	// Well, we got to create a new statusError
	s := &spb.Status{
		Code:    int32(c),
		Message: e.Message(),
	}
	var details ErrorDetails
	if e.Cause() != nil {
		details.Cause = e.Cause().Error()
	}
	if !_env.FromHost().IsPublic() && e.Trace() != nil {
		details.StackTrace = &ErrorDetails_StackTrace{
			Frames: strings.Split(e.Trace().String(), "\n"),
		}
	}
	data, err := proto.Marshal(&details)
	if err == nil {
		s.Details = []*any.Any{{
			TypeUrl: proto.MessageName(&details),
			Value:   data,
		}}
	}
	return status.FromProto(s).Err()
}

// The operation was cancelled, typically by the caller.
//
// HTTP Mapping: 499 Client Closed Request
func Cancelled(err error) error {
	return statusErr(codes.Canceled, err)
}

// Unknown error.  For example, this error may be returned when
// a `Status` value received from another address space belongs to
// an error space that is not known in this address space.  Also
// errors raised by APIs that do not return enough error information
// may be converted to this error.
//
// HTTP Mapping: 500 Internal Server Error
func Unknown(err error) error {
	return statusErr(codes.Unknown, err)
}

// The client specified an invalid argument.  Note that this differs
// from `FAILED_PRECONDITION`.  `INVALID_ARGUMENT` indicates arguments
// that are problematic regardless of the state of the system
// (e.g., a malformed file name).
//
// HTTP Mapping: 400 Bad Request
func InvalidArgument(err error) error {
	return statusErr(codes.InvalidArgument, err)
}

// The deadline expired before the operation could complete. For operations
// that change the state of the system, this error may be returned
// even if the operation has completed successfully.  For example, a
// successful response from a server could have been delayed long
// enough for the deadline to expire.
//
// HTTP Mapping: 504 Gateway Timeout
func DeadlineExceeded(err error) error {
	return statusErr(codes.DeadlineExceeded, err)
}

// Some requested entity (e.g., file or directory) was not found.
// For privacy reasons, this code *may* be returned when the client
// does not have the access rights to the entity, though such usage is
// discouraged.
//
// HTTP Mapping: 404 Not Found
func NotFound(err error) error {
	return statusErr(codes.NotFound, err)
}

// The entity that a client attempted to create (e.g., file or directory)
// already exists.
//
// HTTP Mapping: 409 Conflict
func AlreadyExists(err error) error {
	return statusErr(codes.AlreadyExists, err)
}

// The caller does not have permission to execute the specified
// operation. `PERMISSION_DENIED` must not be used for rejections
// caused by exhausting some resource (use `RESOURCE_EXHAUSTED`
// instead for those errors). `PERMISSION_DENIED` must not be
// used if the caller can not be identified (use `UNAUTHENTICATED`
// instead for those errors).
//
// HTTP Mapping: 403 Forbidden
func PermissionDenied(err error) error {
	return statusErr(codes.PermissionDenied, err)
}

// The request does not have valid authentication credentials for the
// operation.
//
// HTTP Mapping: 401 Unauthorized
func Unauthenticated(err error) error {
	return statusErr(codes.Unauthenticated, err)
}

// Some resource has been exhausted, perhaps a per-user quota, or
// perhaps the entire file system is out of space.
//
// HTTP Mapping: 429 Too Many Requests
func ResourceExhausted(err error) error {
	return statusErr(codes.ResourceExhausted, err)
}

// The operation was rejected because the system is not in a state
// required for the operation's execution.  For example, the directory
// to be deleted is non-empty, an rmdir operation is applied to
// a non-directory, etc.
//
// Service implementors can use the following guidelines to decide
// between `FAILED_PRECONDITION`, `ABORTED`, and `UNAVAILABLE`:
//  (a) Use `UNAVAILABLE` if the client can retry just the failing call.
//  (b) Use `ABORTED` if the client should retry at a higher level
//      (e.g., restarting a read-modify-write sequence).
//  (c) Use `FAILED_PRECONDITION` if the client should not retry until
//      the system state has been explicitly fixed.  E.g., if an "rmdir"
//      fails because the directory is non-empty, `FAILED_PRECONDITION`
//      should be returned since the client should not retry unless
//      the files are deleted from the directory.
//
// HTTP Mapping: 400 Bad Request
func FailedPrecondition(err error) error {
	return statusErr(codes.FailedPrecondition, err)
}

// The operation was aborted, typically due to a concurrency issue such as
// a sequencer check failure or transaction abort.
//
// See the guidelines above for deciding between `FAILED_PRECONDITION`,
// `ABORTED`, and `UNAVAILABLE`.
//
// HTTP Mapping: 409 Conflict
func Aborted(err error) error {
	return statusErr(codes.Aborted, err)
}

// The operation was attempted past the valid range.  E.g., seeking or
// reading past end-of-file.
//
// Unlike `INVALID_ARGUMENT`, this error indicates a problem that may
// be fixed if the system state changes. For example, a 32-bit file
// system will generate `INVALID_ARGUMENT` if asked to read at an
// offset that is not in the range [0,2^32-1], but it will generate
// `OUT_OF_RANGE` if asked to read from an offset past the current
// file size.
//
// There is a fair bit of overlap between `FAILED_PRECONDITION` and
// `OUT_OF_RANGE`.  We recommend using `OUT_OF_RANGE` (the more specific
// error) when it applies so that callers who are iterating through
// a space can easily look for an `OUT_OF_RANGE` error to detect when
// they are done.
//
// HTTP Mapping: 400 Bad Request
func OutOfRange(err error) error {
	return statusErr(codes.OutOfRange, err)
}

// The operation is not implemented or is not supported/enabled in this
// service.
//
// HTTP Mapping: 501 Not Implemented
func Unimplemented(err error) error {
	return statusErr(codes.Unimplemented, err)
}

// Internal errors.  This means that some invariants expected by the
// underlying system have been broken.  This error code is reserved
// for serious errors.
//
// HTTP Mapping: 500 Internal Server Error
func Internal(err error) error {
	return statusErr(codes.Internal, err)
}

// The service is currently unavailable.  This is most likely a
// transient condition, which can be corrected by retrying with
// a backoff.
//
// See the guidelines above for deciding between `FAILED_PRECONDITION`,
// `ABORTED`, and `UNAVAILABLE`.
//
// HTTP Mapping: 503 Service Unavailable
func Unavailable(err error) error {
	return statusErr(codes.Unavailable, err)
}

// Unrecoverable data loss or corruption.
//
// HTTP Mapping: 500 Internal Server Error
func DataLoss(err error) error {
	return statusErr(codes.DataLoss, err)
}

// Ideally schema.py should generate the functions below, but it can't do it today.
// So, this is a manually written for now.
func (m *VoidRequest) IsValid() (*gojsonschema.Result, error) {
	return &gojsonschema.Result{}, nil
}
func (m *VoidRequest) IsRequest() {}
