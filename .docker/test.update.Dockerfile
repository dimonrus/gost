ARG image=image
FROM ${image:-image} AS build
RUN rm -rf /go/src/gost/
COPY ./ /go/src/gost/

