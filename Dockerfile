FROM golang:1.17.6 AS builder
WORKDIR /go/src/github.com/afritzler/garden-universe
RUN go get github.com/rakyll/statik
COPY . .
RUN make

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates=20191127-r7
WORKDIR /
COPY --from=builder /go/src/github.com/afritzler/garden-universe/garden-universe .
ENTRYPOINT [ "/garden-universe"]