# GO Images https://hub.docker.com/_/golang/
FROM golang:1.14.4-alpine3.12 AS build
# Install GIT
RUN apk add --no-cache git
# Work directory
WORKDIR /go/src/gost/
# Copy project
COPY ./ /go/src/gost/
# Build app
RUN go build -o gost main.go

