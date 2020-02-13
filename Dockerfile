# Buildtime
FROM golang:1.13 AS builder

WORKDIR /workspace
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

# Runtime
# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/bin/kubevent .
USER nonroot:nonroot

ENTRYPOINT ["/kubevent"]
