module git.2dfire.net/zerodb/keeper

go 1.12

require (
	git.2dfire.net/zerodb/common v0.0.0
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.7.7
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/mock v1.1.1
	github.com/google/uuid v1.1.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.3.0
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v3.3.18+incompatible
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20190709130402-674ba3eaed22
)

replace git.2dfire.net/zerodb/common => ../common
