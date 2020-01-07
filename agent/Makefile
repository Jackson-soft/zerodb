#Makefile for ZeroDB

work_dir=$(shell pwd)
ZEROAGENT_CMD=${work_dir}/cmd

all: build

clean:
	rm -f $(ZEROAGENT_CMD)/*.out

build:
	go clean
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o ${ZEROAGENT_CMD}/zeroAgent.out ${ZEROAGENT_CMD}/

