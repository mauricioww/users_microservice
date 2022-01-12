package transport_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	srv_mock := new(transport.ServiceMock)
	endpoints := transport.MakeHttpEndpoints(srv_mock)
	s := transport.NewHTTPServer(context.Background(), endpoints)
	server := httptest.NewServer(s)

	defer server.Close()

	test_cases := []struct {
		test_name  string
		body       string
		data       transport.CreateUserRequest
		res        int
		err        error
		httpStatus int
	}{
		{
			test_name: "user created succes",
			body: `
				{
					"email": "example@email.com",
					"password": "querty", 
					"age": 25
				}`,
			data: transport.CreateUserRequest{
				Email:    "example@email.com",
				Password: "querty",
				Age:      25,
			},
			res:        1,
			err:        nil,
			httpStatus: 200,
		},
		{
			test_name: "no password error",
			body: `
				{
					"email": "example@email.com",
					"age": 25
				}
			`,
			data: transport.CreateUserRequest{
				Email: "example@email.com",
				Age:   25,
			},
			res:        -1,
			err:        status.Error(codes.FailedPrecondition, "Missing field 'password'"),
			httpStatus: 400,
		},
		{
			test_name: "no email error",
			body: `
				{
					"password": "qwerty",
					"age": 25
				}
			`,
			data: transport.CreateUserRequest{
				Password: "qwerty",
				Age:      25,
			},
			res:        -1,
			err:        status.Error(codes.FailedPrecondition, "Missing field 'email'"),
			httpStatus: 400,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			srv_mock.On("CreateUser", mock.Anything, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details).
				Return(tc.res, tc.err)
			res, _ := http.Post(server.URL+"/users", "application/json", strings.NewReader(tc.body))

			// assert
			assert.Equal(tc.httpStatus, res.StatusCode)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	srv_mock := new(transport.ServiceMock)
	endpoints := transport.MakeHttpEndpoints(srv_mock)
	s := transport.NewHTTPServer(context.Background(), endpoints)
	server := httptest.NewServer(s)

	defer server.Close()

	test_cases := []struct {
		test_name  string
		body       string
		data       transport.AuthenticateRequest
		res        string
		err        error
		httpStatus int
	}{
		{
			test_name: "user auth succes",
			body: `
				{
					"email": "example@email.com",
					"password": "querty"
				}`,
			data: transport.AuthenticateRequest{
				Email:    "example@email.com",
				Password: "querty",
			},
			res:        "jwt_token",
			err:        nil,
			httpStatus: 200,
		},
		{
			test_name: "no password error",
			body: `
				{
					"email": "example@email.com"
				}
			`,
			data: transport.AuthenticateRequest{
				Email: "example@email.com",
			},
			err:        status.Error(codes.FailedPrecondition, "Missing field 'password'"),
			httpStatus: 400,
		},
		{
			test_name: "no email error",
			body: `
				{
					"password": "qwerty"
				}
			`,
			data: transport.AuthenticateRequest{
				Password: "qwerty",
			},
			err:        status.Error(codes.FailedPrecondition, "Missing field 'email'"),
			httpStatus: 400,
		},
		{
			test_name: "auth error",
			body: `
				{
					"email": "example@email.com",
					"password": "qwerty"
				}
			`,
			data: transport.AuthenticateRequest{
				Email:    "example@email.com",
				Password: "qwerty",
			},
			err:        status.Error(codes.Unauthenticated, "Password or email error"),
			httpStatus: 401,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			srv_mock.On("Authenticate", mock.Anything, tc.data.Email, tc.data.Password).Return(tc.res, tc.err)

			req, _ := http.NewRequest("GET", server.URL+"/auth", strings.NewReader(tc.body))
			res, _ := http.DefaultClient.Do(req)

			// assert
			assert.Equal(tc.httpStatus, res.StatusCode)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	srv_mock := new(transport.ServiceMock)
	endpoints := transport.MakeHttpEndpoints(srv_mock)
	s := transport.NewHTTPServer(context.Background(), endpoints)
	server := httptest.NewServer(s)

	defer server.Close()

	test_cases := []struct {
		test_name  string
		user_id    int
		body       string
		data       transport.UpdateUserRequest
		res        bool
		err        error
		httpStatus int
	}{
		{
			test_name: "user update succes",
			user_id:   0,
			body: `
				{
					"email": "example@email.com",
					"password": "querty", 
					"age": 25
				}`,
			data: transport.UpdateUserRequest{
				Email:    "example@email.com",
				Password: "querty",
				Age:      25,
			},
			res:        true,
			err:        nil,
			httpStatus: 200,
		},
		{
			test_name: "no password error",
			user_id:   1,
			body: `
				{
					"email": "example@email.com",
					"age": 25
				}
			`,
			data: transport.UpdateUserRequest{
				Email: "example@email.com",
				Age:   25,
			},
			err:        status.Error(codes.FailedPrecondition, "Missing field 'password'"),
			httpStatus: 400,
		},
		{
			test_name: "no email error",
			user_id:   2,
			body: `
				{
					"password": "qwerty",
					"age": 25
				}
			`,
			data: transport.UpdateUserRequest{
				Password: "qwerty",
				Age:      25,
			},
			err:        status.Error(codes.FailedPrecondition, "Missing field 'email'"),
			httpStatus: 400,
		},
		{
			test_name: "user not found",
			user_id:   100,
			body: `
				{
					"email": "example@email.com",
					"password": "qwerty",
					"age": 25
				}
			`,
			data: transport.UpdateUserRequest{
				Email:    "example@email.com",
				Password: "qwerty",
				Age:      25,
			},
			err:        status.Error(codes.NotFound, "User not found"),
			httpStatus: 404,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			srv_mock.On("UpdateUser", mock.Anything, tc.user_id, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details).Return(tc.res, tc.err)

			uri := fmt.Sprintf("%v/users/%v", server.URL, tc.user_id)
			req, _ := http.NewRequest("PUT", uri, strings.NewReader(tc.body))
			res, _ := http.DefaultClient.Do(req)

			// assert
			assert.Equal(tc.httpStatus, res.StatusCode)
		})
	}
}

func TestGetUser(t *testing.T) {
	srv_mock := new(transport.ServiceMock)
	endpoints := transport.MakeHttpEndpoints(srv_mock)
	s := transport.NewHTTPServer(context.Background(), endpoints)
	server := httptest.NewServer(s)

	defer server.Close()

	test_cases := []struct {
		test_name  string
		user_id    int
		res        entities.User
		err        error
		httpStatus int
	}{
		{
			test_name: "user found success",
			user_id:   0,
			res: entities.User{
				Email:    "email@domain.com",
				Password: "passsword",
			},
			err:        nil,
			httpStatus: 200,
		},
		{
			test_name:  "user not found error",
			user_id:    1,
			err:        status.Error(codes.NotFound, "User not found"),
			httpStatus: 404,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			srv_mock.On("GetUser", mock.Anything, tc.user_id).Return(tc.res, tc.err)

			uri := fmt.Sprintf("%v/users/%v", server.URL, tc.user_id)
			res, _ := http.Get(uri)

			// assert
			assert.Equal(res.StatusCode, tc.httpStatus)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	srv_mock := new(transport.ServiceMock)
	endpoints := transport.MakeHttpEndpoints(srv_mock)
	s := transport.NewHTTPServer(context.Background(), endpoints)
	server := httptest.NewServer(s)

	defer server.Close()

	test_cases := []struct {
		test_name  string
		user_id    int
		res        bool
		err        error
		httpStatus int
	}{
		{
			test_name:  "user deleted success",
			user_id:    0,
			res:        true,
			err:        nil,
			httpStatus: 200,
		},
		{
			test_name:  "user not found error",
			user_id:    1,
			err:        status.Error(codes.NotFound, "User not found"),
			httpStatus: 404,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			srv_mock.On("DeleteUser", mock.Anything, tc.user_id).Return(tc.res, tc.err)

			uri := fmt.Sprintf("%v/users/%v", server.URL, tc.user_id)
			req, _ := http.NewRequest("DELETE", uri, http.NoBody)
			res, _ := http.DefaultClient.Do(req)

			// assert
			assert.Equal(res.StatusCode, tc.httpStatus)
		})
	}
}
