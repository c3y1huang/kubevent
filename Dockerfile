# Buildtime
FROM golang:alpine AS builder

ADD . "$GOPATH/src/github.com/innobead/kubevent"
WORKDIR "$GOPATH/src/github.com/innobead/kubevent"

RUN apk update && \
    apk add git build-base && \
    cd "$GOPATH/src/github.com/innobead/kubevent" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/kubevent -o /kubevent

# Runtime
FROM alpine:3.10

RUN apk add --update ca-certificates

COPY --from=builder kubevent /bin/kubevent
COPY --from=builder configs /$HOME/.kubevent

ENTRYPOINT ["/bin/kubevent"]