#!/usr/bin/env bash

PROGRAM=$(basename "$0")

if [ -z "$GOPATH" ]; then
    printf "Error: the environment variable GOPATH is not set, please set it before running %s\n" $PROGRAM >/dev/stderr
    exit 1
fi

if hash protoc-gen-gofast 2>/dev/null; then
    echo "exists protoc-gen-gofast"
else
    echo "install protoc-gen-gofast ..."
    go get -u github.com/golang/protobuf/{proto,protoc-gen-gofast}
fi

if hash goimports 2>/dev/null; then
    echo 'exists goimports'
else
    echo "install goimports ..."
    go get -u golang.org/x/tools/cmd/goimports
fi

GO_PREFIX_PATH=./zeroproto/pkg
GO_OUT_M=

cd ./zeroproto/proto
for file in $(ls *.proto); do
    base_name=$(basename "$file" ".proto")
    mkdir -p ../pkg/"$base_name"
    if [ -z $GO_OUT_M ]; then
        GO_OUT_M="M$file=$GO_PREFIX_PATH/$base_name"
    else
        GO_OUT_M="$GO_OUT_M,M$file=$GO_PREFIX_PATH/$base_name"
    fi
done

echo "generate go code..."
ret=0
for file in $(ls *.proto); do
    base_name=$(basename "$file" ".proto")
    protoc -I. --gofast_out=plugins=grpc,"$GO_OUT_M":../pkg/"$base_name" "$file" || ret=$?
    cd ../pkg/"$base_name" || exit
    sed -i.bak -E 's/import _ \"gogoproto\"//g' *.pb.go
    sed -i.bak -E 's/import fmt \"fmt\"//g' *.pb.go
    sed -i.bak -E 's/import io \"io\"//g' *.pb.go
    sed -i.bak -E 's/import math \"math\"//g' *.pb.go
    sed -i -E 's#./zeroproto/pkg/basic#git.2dfire.net/zerodb/common/zeroproto/pkg/basic#g' *.pb.go
    rm -f *.bak
    rm -f *.go-E
    goimports -w *.pb.go
    cd ../../proto
done

exit $ret
