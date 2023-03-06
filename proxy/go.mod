module git.2dfire.net/zerodb/proxy

go 1.12

require (
	git.2dfire.net/zerodb/common v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20190709130402-674ba3eaed22
)

replace git.2dfire.net/zerodb/common => ../common
