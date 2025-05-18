# Get image from argument
ARG image=image
# Set allias for image
FROM ${image:-image} AS build
# Getting aline image
FROM alpine:latest
# Create gost user
RUN addgroup -S gost && adduser -u 1000 -S -D -H gost -G gost
# Copy certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy bin file
COPY --chown=gost --from=build /go/src/gost/gost /gost/
# Copy etc files
COPY --chown=gost --from=build /go/src/gost/etc /
# Copy configs
COPY --chown=gost --from=build /go/src/gost/app/config/yaml /gost/app/config/yaml
# Copy resource
COPY --chown=gost --from=build /go/src/gost/resource/ /gost/resource/
# Run as user
USER 1000
# Entry
ENTRYPOINT ["gost/gost", "-app", "consumer"]

