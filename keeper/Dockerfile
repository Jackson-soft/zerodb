FROM golang:1.13-alpine as builder

RUN apk --no-cache add git

RUN git clone http://git.2dfire.net/zerodb/common.git /root/common

RUN mkdir -p /root/keeper
WORKDIR /root/keeper

COPY . .
RUN GOCACHE=/tmp/gocache go build -mod=vendor -o ./cmd/zeroKeeper.out ./cmd/main.go

FROM alpine:latest

RUN mkdir -p /opt/keeper
WORKDIR /opt/keeper

COPY --from=builder /root/keeper/cmd/zeroKeeper.out ./zeroKeeper
COPY --from=builder /root/keeper/etc/*.yaml ./etc/
COPY entrypoint.sh ./entrypoint.sh

RUN chmod -R 777 /opt/keeper/

ENTRYPOINT ["/opt/keeper/entrypoint.sh"]
