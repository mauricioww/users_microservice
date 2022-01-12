package errors

import (
	"google.golang.org/grpc/codes"
)

const (
	badRequestEmail    = 0
	badRequestPassword = 1
	unauthenticated    = 2
	unauthorized       = 3
	userNotFound       = 4
	serverFail         = 5
)

var message = map[int]string{
	badRequestEmail:    "Missing field 'email'",
	badRequestPassword: "Missing field 'password'",
	unauthenticated:    "Password or email error",
	unauthorized:       "Unauthorized user",
	userNotFound:       "User not found",
	serverFail:         "Internal server error",
}

func MessageError(code int) string {
	return message[code]
}

func resolveGrpc(err int) codes.Code {
	switch err {
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

func ResolveHttp(c codes.Code) int {
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

type (
	BadRequestEmailError    int
	BadRequestPasswordError int
	InternalError           int
	UserNotFoundError       int
	UnauthorizedError       int
	UnauthenticatedError    int

	ErrorResolver interface {
		GrpcCode() codes.Code
	}
)

func NewBadRequestEmailError() BadRequestEmailError {
	return badRequestEmail
}

func NewBadRequestPasswordError() BadRequestPasswordError {
	return badRequestPassword
}

func NewInternalError() InternalError {
	return serverFail
}

func NewUserNotFoundError() UserNotFoundError {
	return userNotFound
}

func NewUnauthorizedError() UnauthorizedError {
	return unauthorized
}

func NewUnauthenticatedError() UnauthenticatedError {
	return unauthenticated
}

func (e BadRequestEmailError) Error() string {
	return MessageError(int(e))
}

func (e BadRequestPasswordError) Error() string {
	return MessageError(int(e))
}

func (e InternalError) Error() string {
	return MessageError(int(e))
}

func (e UserNotFoundError) Error() string {
	return MessageError(int(e))
}

func (e UnauthorizedError) Error() string {
	return MessageError(int(e))
}

func (e UnauthenticatedError) Error() string {
	return MessageError(int(e))
}

func (e BadRequestEmailError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

func (e BadRequestPasswordError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

func (e InternalError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

func (e UserNotFoundError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

func (e UnauthorizedError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}

func (e UnauthenticatedError) GrpcCode() codes.Code {
	return resolveGrpc(int(e))
}
