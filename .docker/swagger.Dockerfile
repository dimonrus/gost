# Get image from argument
ARG image=image
# Set allias for image
FROM ${image:-image} AS build
# Install curl
RUN apk add curl git
# Download swagger
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest
# Set up workdir
WORKDIR /go/src/gost/
# Generate spec
RUN make swagger-spec