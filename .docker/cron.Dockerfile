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
# Create folder for symlink
RUN mkdir -p /root/gost/app/config
# Create symlink
RUN ln -s /gost/app/config/yaml /root/gost/app/config/yaml
# Cronbab permission
RUN chmod 0644 cron/crontab
# Cron bash permission
RUN chmod 0777 cron/cronap.sh
# Update crontab
RUN crontab cron/crontab
# Entry
ENTRYPOINT ["sh", "cron/cronap.sh"]

