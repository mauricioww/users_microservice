package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	gokitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/mauricioww/user_microsrv/errors"
	"google.golang.org/grpc/status"
)

// NewHTTPServer returns the server with the endpoints and the specifications for each one
func NewHTTPServer(ctx context.Context, endpoints HTTPEndpoints) http.Handler {
	root := mux.NewRouter()
	root.Use(middleware)

	userRouter := root.PathPrefix("/users").Subrouter()

	opt := gokitHttp.ServerOption(gokitHttp.ServerErrorEncoder(encodeError))

	userRouter.Methods("GET").Path("/{id}").Handler(gokitHttp.NewServer(
		endpoints.GetUser,
		decodeGetUserRequest,
		encodeResponse,
		opt,
	))

	userRouter.Methods("POST").Handler(gokitHttp.NewServer(
		endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeResponse,
		opt,
	))

	userRouter.Methods("PUT").Path("/{id}").Handler(gokitHttp.NewServer(
		endpoints.UpdateUser,
		decodeUpdateUserRequest,
		encodeResponse,
		opt,
	))

	userRouter.Methods("DELETE").Path("/{id}").Handler(gokitHttp.NewServer(
		endpoints.DeleteUser,
		decodeDeleteUserRequest,
		encodeResponse,
		opt,
	))

	root.Methods("POST").Path("/auth").Handler(gokitHttp.NewServer(
		endpoints.Authenticate,
		decodeAuthenticateRequest,
		encodeResponse,
		opt,
	))

	return root
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAuthenticateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request AuthenticateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request UpdateUserRequest
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.UserID = id
	return request, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return nil, err
	}

	request := GetUserRequest{UserID: id}
	return request, nil
}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return nil, err
	}

	request := DeleteUserRequest{UserID: id}
	return request, nil
}

func encodeResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	e, ok := status.FromError(err)
	if !ok {
		panic("Unsupported error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errors.ResolveHTTP(e.Code()))
	json.NewEncoder(w).Encode(map[string]string{"error": e.Message()})
}
