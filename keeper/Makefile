### Makefile for ZeroDB

ZEROKEEPER_CMD=./cmd
ZEROPROTO_GO=../common/zeroproto/pkg
ZEROKEEPER_AGENTCLINT_MOCK=./server/mock/agentmock
ZEROKEEPER_PROXYCLINT_MOCK=./server/mock/proxymock
ZEROKEEPER_KEEPERCLINT_MOCK=./server/mock/keepermock

all: clean build

clean:
	rm -f $(ZEROKEEPER_CMD)/*.out
	rm -f $(ZEROKEEPER_AGENTCLINT_MOCK)/*.go
	rm -f $(ZEROKEEPER_PROXYCLINT_MOCK)/*.go
	rm -f $(ZEROKEEPER_KEEPERCLINT_MOCK)/*.go

build:
	go clean
	GOOS=linux GOARCH=amd64 GOCACHE=/tmp/gocache go build -mod=vendor -o $(ZEROKEEPER_CMD)/zeroKeeper.out $(ZEROKEEPER_CMD)/main.go

mock:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	mkdir -p $(ZEROKEEPER_AGENTCLINT_MOCK) $(ZEROKEEPER_PROXYCLINT_MOCK) $(ZEROKEEPER_KEEPERCLINT_MOCK)
	mockgen git.2dfire.net/zerodb/common/zeroproto/pkg/agent AgentClient  > $(ZEROKEEPER_AGENTCLINT_MOCK)/mock_AgentClient.go
	mockgen git.2dfire.net/zerodb/common/zeroproto/pkg/keeper KeeperClient  > $(ZEROKEEPER_KEEPERCLINT_MOCK)/mock_KeeperClient.go
	mockgen git.2dfire.net/zerodb/common/zeroproto/pkg/proxy ProxyClient  >  $(ZEROKEEPER_PROXYCLINT_MOCK)/mock_ProxyClient.go
