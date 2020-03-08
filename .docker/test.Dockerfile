ARG image=image
FROM ${image:-image} AS build
RUN apk add --no-cache gcc musl-dev
CMD go test ./... -v

