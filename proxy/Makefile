### Makefile for ZeroDB

work_dir=$(shell pwd)
ZEROPROXY=${work_dir}/cmd/zeroproxy
ZERO_GOYACC=${work_dir}/cmd/goyacc

all: goyacc build

clean:
	rm -f $(ZERO_GOYACC)/*.out
	rm -f $(ZEROPROXY)/*.out

goyacc:
	go build -o ${ZERO_GOYACC}/goyacc.out ${ZERO_GOYACC}/
	${ZERO_GOYACC}/goyacc.out -o ${work_dir}/proxy/sqlparser/sql.go ${work_dir}/proxy/sqlparser/sql.y

build:
	go clean
	GOOS=linux GOARCH=amd64 GOCACHE=/tmp/gocache go build -mod=vendor -o $(ZEROPROXY)/zeroProxy.out  $(ZEROPROXY)/
