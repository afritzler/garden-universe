FROM golang:1.13.0 AS builder
WORKDIR /go/src/github.com/afritzler/garden-universe
RUN go get github.com/rakyll/statik
COPY . .
RUN make

FROM alpine:3.13.4
RUN apk --no-cache add ca-certificates=20191127
WORKDIR /root/
COPY --from=builder /go/src/github.com/afritzler/garden-universe/garden-universe .
ENTRYPOINT [ "/garden-universe"]