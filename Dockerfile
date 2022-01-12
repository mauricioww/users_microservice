FROM golang:alpine AS base_app

RUN apk update && apk add --no-cache git

WORKDIR /src

COPY . .

RUN go mod download

RUN GOOS=linux go build -ldflags="-w -s" -o http ./http_srv/.
RUN GOOS=linux go build -ldflags="-w -s" -o user ./user_srv/.
RUN GOOS=linux go build -ldflags="-w -s" -o details ./user_details_srv/.

#############################################################################

FROM golang:alpine AS http_server

WORKDIR /

COPY --from=base_app ./src/http ./

ENTRYPOINT ["/http"]

#############################################################################

FROM golang:alpine AS grpc_user_server

WORKDIR /

COPY --from=base_app ./src/user ./

ENTRYPOINT ["/user"]


#############################################################################

FROM golang:alpine AS grpc_details_server

WORKDIR /

COPY --from=base_app ./src/details ./

ENTRYPOINT ["/details"]