package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	gokit_http "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/mauricioww/user_microsrv/errors"
	"google.golang.org/grpc/status"
)

func NewHTTPServer(ctx context.Context, http_endpoints HttpEndpoints) http.Handler {
	root := mux.NewRouter()
	root.Use(middleware)

	user_router := root.PathPrefix("/users").Subrouter()

	opt := gokit_http.ServerOption(gokit_http.ServerErrorEncoder(encodeError))

	user_router.Methods("GET").Path("/{id}").Handler(gokit_http.NewServer(
		http_endpoints.GetUser,
		decodeGetUserRequest,
		encodeResponse,
		opt,
	))

	user_router.Methods("POST").Handler(gokit_http.NewServer(
		http_endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeResponse,
		opt,
	))

	user_router.Methods("PUT").Path("/{id}").Handler(gokit_http.NewServer(
		http_endpoints.UpdateUser,
		decodeUpdateUserRequest,
		encodeResponse,
		opt,
	))

	user_router.Methods("DELETE").Path("/{id}").Handler(gokit_http.NewServer(
		http_endpoints.DeleteUser,
		decodeDeleteUserRequest,
		encodeResponse,
		opt,
	))

	root.Methods("POST").Path("/auth").Handler(gokit_http.NewServer(
		http_endpoints.Authenticate,
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
	id_param := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_param)

	if err != nil {
		return nil, nil
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}

	request.UserId = id
	return request, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id_param := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_param)

	if err != nil {
		return nil, err
	}

	request := GetUserRequest{UserId: id}
	return request, nil
}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id_param := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_param)

	if err != nil {
		return nil, err
	}

	request := DeleteUserRequest{UserId: id}
	return request, nil
}

func encodeResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	e, _ := status.FromError(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errors.ResolveHttp(e.Code()))
	json.NewEncoder(w).Encode(map[string]string{"error": e.Message()})
}
