# Latest version
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION} AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /main ./cmd/proto/main.go

# the running container.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /srv
COPY --from=builder /main /srv/main

EXPOSE 50051

ENTRYPOINT ["/srv/main"]

