# Builder stage
FROM golang:1.23.4-alpine AS builder

ENV GO111MODULE=on

RUN adduser -D -u 1000 dev

ARG importPath
ARG pkg

WORKDIR /go/src/${importPath}

COPY . .

# Optimized Go build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app ${pkg}

# Runner stage
FROM golang:1.23.4-alpine

ARG importPath

RUN adduser -D -u 1000 dev

COPY --from=builder /go/src/${importPath}/app .

USER dev

EXPOSE 8080

ENTRYPOINT ["./app"]