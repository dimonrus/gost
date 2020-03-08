# Get image from argument
ARG image=image
# Set allias for image
FROM ${image:-image} AS build
# Getting aline image
FROM alpine:latest
# Copy certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy bin file
COPY --from=build /go/src/gost/gost /gost/
# Copy etc files
COPY --from=build /go/src/gost/etc /
# Copy configs
COPY --from=build /go/src/gost/app/config/yaml /gost/app/config/yaml
# Copy resource
COPY --from=build /go/src/gost/resource/ /gost/resource/
# Entry
ENTRYPOINT ["gost/gost", "-app", "consumer"]

