# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.22.2 as builder

ARG GOARCH=''

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

RUN go install github.com/rakyll/statik@latest

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY statik/ statik/
COPY web/ web/

ARG TARGETOS
ARG TARGETARCH

# Build
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o garden-universe main.go

FROM alpine:3.19.1
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /workspace/garden-universe .
USER 65532:65532
ENTRYPOINT [ "/garden-universe"]