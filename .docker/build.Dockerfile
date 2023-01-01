# Get image from argument
ARG image=image
# Set allias for image
FROM ${image:-image} AS build
# Install Make and git
RUN apk add --update make git
# Remove codebase
RUN rm -rf /go/src/gost/
# Copy current codebase
COPY ./ /go/src/gost/
# Set work directory
WORKDIR /go/src/gost/
# Build app
RUN make project-build