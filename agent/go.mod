module git.2dfire.net/zerodb/agent

go 1.12

require (
	git.2dfire.net/zerodb/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.1
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20190709130402-674ba3eaed22
)

replace git.2dfire.net/zerodb/common => ../common
