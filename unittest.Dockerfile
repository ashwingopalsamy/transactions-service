FROM golang:1.23.4-alpine

# Build args
ARG importPath

# Configure WORKDIR in Go convention
WORKDIR $GOPATH/src/${importPath}

# Copy contents
COPY . .
