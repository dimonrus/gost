# Get image from argument
ARG image=image
# Set allias for image
FROM ${image:-image} AS build
# Install curl
RUN apk add --no-cache curl
# Download swagger
RUN curl -o /swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.22.0/swagger_linux_amd64
# Permission for swagger
RUN chmod +x /swagger
# Set up workdir
WORKDIR /go/src/gost/
# Generate spec
RUN /swagger generate spec -m -o resource/swagger.json