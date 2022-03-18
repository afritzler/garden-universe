FROM golang:1.18.0 AS builder
WORKDIR /go/src/github.com/afritzler/garden-universe
RUN go install github.com/rakyll/statik@latest
COPY . .
RUN make

FROM alpine:3.15.1
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/afritzler/garden-universe/garden-universe .
ENTRYPOINT [ "/garden-universe"]