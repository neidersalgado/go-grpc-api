############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
# Fetch dependencies.
# Using go get.
RUN go mod download
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o grpcservice ./cmd/grpc-server/.

############################
# STEP 2 grpc service
############################
FROM golang:alpine as grpcserver

WORKDIR /

COPY --from=builder /src/grpcservice ./grpcservice

ENTRYPOINT ["/grpcservice"]
