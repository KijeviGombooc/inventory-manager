# Stage 1: Development
FROM alpine:3.20 AS dev
# Install dependencies
RUN apk add --no-cache \
    git \
    openssh \
    go \
    docker-cli
# Install delve debugger for Go
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest