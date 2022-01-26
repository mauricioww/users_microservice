package errors

import (
	"google.golang.org/grpc/codes"
)

const (
	unknownError       = -1
	badRequestEmail    = 0
	badRequestPassword = 1
	unauthenticated    = 2
	unauthorized       = 3
	userNotFound       = 4
	serverFail         = 5
)

var message = map[int]string{
	unknownError:       "Unsupported error",
	badRequestEmail:    "Missing field 'email'",
	badRequestPassword: "Missing field 'password'",
	unauthenticated:    "Password or email error",
	unauthorized:       "Unauthorized user",
	userNotFound:       "User not found",
	serverFail:         "Internal server error",
}

func messageError(code int) string {
	return message[code]
}

func resolveGrpc(err int) codes.Code {
	switch err {
	case unknownError:
		return codes.Unknown
	case badRequestEmail, badRequestPassword:
		return codes.FailedPrecondition
	case unauthenticated, unauthorized:
		return codes.Unauthenticated
	case userNotFound:
		return codes.NotFound
	default:
		return codes.Internal
	}
}

// ResolveHTTP is a function to transform gRPC codes to HTTP codes
func ResolveHTTP(c codes.Code) int {
	switch c {
	case codes.FailedPrecondition:
		return 400
	case codes.Unauthenticated:
		return 401
	case codes.NotFound:
		return 404
	default:
		return 500
	}
}

// UnknownError used when the current error is not within these custom errors
type UnknownError int

// BadRequestEmailError used when field email is missing
type BadRequestEmailError int

// BadRequestPasswordError used when field password is missing
type BadRequestPasswordError int

// InternalError used when server error happens
type InternalError int

// UserNotFoundError used when an user does not exits
type UserNotFoundError int

// UnauthorizedError used when an user does not have the grants to do operations
type UnauthorizedError int

// UnauthenticatedError used when a login error happens
type UnauthenticatedError int

// ErrorResolver is an interfaz shared between my custom errors to handle errors in the services
type ErrorResolver interface {
	GrpcCode() codes.Code
}

// NewUnknownError returns a badRequestEmail error type
func NewUnknownError() UnknownError {
	return unknownError
}

// NewBadRequestEmailError returns a badRequestEmail error type
func NewBadRequestEmailError() BadRequestEmailError {
	return badRequestEmail
}

// NewBadRequestPasswordError returns a badRequestPassword error type
func NewBadRequestPasswordError() BadRequestPasswordError {
	return badRequestPassword
}

// NewInternalError returns a internal error type
func NewInternalError() InternalError {
	return serverFail
}

// NewUserNotFoundError returns a userNotFound error type
func NewUserNotFoundError() UserNotFoundError {
	return userNotFound
}

// NewUnauthorizedError returns a unauthorized error type
func NewUnauthorizedError() UnauthorizedError {
	return unauthorized
}

// NewUnauthenticatedError returns a unauthenticated error type
func NewUnauthenticatedError() UnauthenticatedError {
	return unauthenticated
}

func (e UnknownError) Error() string {
	return messageError(int(e))
}

func (e BadRequestEmailError) Error() string {
	return messageError(int(e))
}

func (e BadRequestPasswordError) Error() string {
	return messageError(int(e))
}

func (e InternalError) Error() string {
	return messageError(int(e))
}

func (e UserNotFoundError) Error() string {
	return messageError(int(e))
}

func (e UnauthorizedError) Error() string {
	return messageError(int(e))
}

func (e UnauthenticatedError) Error() string {
	return messageError(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e UnknownError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e BadRequestEmailError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e BadRequestPasswordError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e InternalError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e UserNotFoundError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e UnauthorizedError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

// GrpcCode translate from HTTP code to gRPC code
func (e UnauthenticatedError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}
