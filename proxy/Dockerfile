FROM golang:1.13-alpine as builder

RUN apk --no-cache add git

RUN git clone http://git.2dfire.net/zerodb/common.git /root/common

RUN mkdir -p /root/proxy
WORKDIR /root/proxy
COPY . .
RUN GOCACHE=/tmp/gocache go build -mod=vendor -o ./cmd/zeroproxy/zeroProxy.out ./cmd/zeroproxy/main.go

FROM alpine:latest

RUN mkdir -p /opt/proxy
WORKDIR /opt/proxy

COPY --from=builder /root/proxy/cmd/zeroproxy/zeroProxy.out ./zeroProxy
COPY --from=builder /root/proxy/etc/*.yaml ./etc/
COPY entrypoint.sh ./entrypoint.sh

RUN chmod -R 777 /opt/proxy

ENTRYPOINT ["/opt/proxy/entrypoint.sh"]
