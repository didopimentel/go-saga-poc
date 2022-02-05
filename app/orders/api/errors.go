package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorType int

// Error is a common wrapper for any failure situation dealt with at the API level
// In the case we have no specific error to return, the NewInternalError should be
// used as fallback, corresponding to the HTTP 500 status code.
type Error struct {
	grpcCode   codes.Code
	httpStatus int
	message    string
}

func (e Error) Error() string {
	if e.message != "" {
		return e.message
	}

	return fmt.Sprintf("error. grpcCode: %d. httpStatus: %d", e.grpcCode, e.httpStatus)
}

func (e Error) HTTPStatus() int {
	return e.httpStatus
}

func (e Error) HTTPError() string {
	if len(e.message) == 0 {
		return ""
	}

	if isJSON(e.message) {
		return fmt.Sprintf(`{"error": %s}`, e.message)
	}

	message := removeEnclosingDoubleQuotes(e.message)
	message = strings.ReplaceAll(message, `"`, `\"`)

	return fmt.Sprintf(`{"error": "%s"}`, message)
}

func (e Error) GRPCError() error {
	return status.Error(e.grpcCode, e.message)
}

func isJSON(s string) bool {
	var js json.RawMessage

	return json.Unmarshal([]byte(s), &js) == nil

}

func removeEnclosingDoubleQuotes(s string) string {
	result := s
	re := regexp.MustCompile(`^".*"$`)

	if re.MatchString(result) {
		result = result[1 : len(result)-1]
	}

	return result
}

func NewPermissionDenied(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "permission denied"
	}
	grpcCode := codes.PermissionDenied

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewUnauthenticatedError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "no authentication provided"
	}
	grpcCode := codes.Unauthenticated

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewConflictError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "already exists"
	}
	grpcCode := codes.AlreadyExists

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewBadRequestError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "invalid request"
	}
	grpcCode := codes.InvalidArgument

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewFailedPreconditionError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "failed precondition"
	}

	// FailedPrecondition is translated to BadRequest on HTTP
	grpcCode := codes.FailedPrecondition

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewResourceExhaustedError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "resource exhausted"
	}
	grpcCode := codes.ResourceExhausted

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewNotFoundError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "not found"
	}
	grpcCode := codes.NotFound

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewDataLossError(messageFormat string, messageArgs ...interface{}) *Error {
	if messageFormat == "" {
		messageFormat = "failed to process headers"
	}
	grpcCode := codes.DataLoss

	return &Error{
		message:    fmt.Sprintf(messageFormat, messageArgs...),
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

func NewInternalError() *Error {
	grpcCode := codes.Internal

	return &Error{
		message:    "something wrong happened",
		grpcCode:   grpcCode,
		httpStatus: runtime.HTTPStatusFromCode(grpcCode),
	}
}

// FromGRPCError converts any golang error to an *Error with a corresponding grpcCode
func FromGRPCError(err error) *Error {
	grpcError := status.Convert(err)
	code := status.Code(err)

	return &Error{
		grpcCode:   code,
		httpStatus: runtime.HTTPStatusFromCode(code),
		message:    grpcError.Message(),
	}
}

// FromError returns an *Error from any given error
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		return e
	}

	return NewInternalError()
}

func NewFieldValidationError(value, field string) *Error {
	return NewBadRequestError("value %s is invalid for field %s", value, field)
}

type FieldValidationError struct {
	Field string
	Value string
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("value %s is invalid for field %s", e.Value, e.Field)
}
