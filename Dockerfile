FROM golang:1.20 AS builder
ENV GOPROXY https://goproxy.io
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -mod vendor -o /basicauth-proxy

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /basicauth-proxy /basicauth-proxy
CMD ["/basicauth-proxy"]
