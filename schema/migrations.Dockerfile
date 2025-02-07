# Migration stage
FROM golang:1.23.4-alpine AS migration_builder

ENV GO111MODULE=on

RUN adduser -D -u 1000 dev

ARG importPath

WORKDIR /go/src/${importPath}

COPY . .

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

# Runner stage
FROM alpine:latest

ARG importPath

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1000 dev

COPY --from=migration_builder /go/bin/goose /usr/local/bin/goose

USER dev

WORKDIR /go/src/${importPath}

ENTRYPOINT ["goose"]
CMD ["-dir", "/migrations", "postgres"]
