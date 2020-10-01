FROM golang:1.15.2-alpine3.12 AS tester

WORKDIR /src
# Prevent updates to source from causing modules to download again
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test ./...
