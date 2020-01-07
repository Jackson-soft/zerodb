package frontend

import (
	"io/ioutil"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"gopkg.in/yaml.v3"
)

var testProxyEngine *ProxyEngine
var testRouter *Router

// init test proxy server
func init() {
	err := glog.CreateLogs("../../zlog_test", "info")
	if err != nil {
		panic(err)
	}
	rt, err := newTestRouter("../test-conf/router_conf_test.yaml")
	if err != nil {
		panic(err)
	}
	testRouter = rt

	ps, err := newTestProxyService("../test-conf/proxy_conf_test.yaml")
	if err != nil {
		panic(err)
	}
	testProxyEngine = ps
	started := make(chan int)
	go testProxyEngine.ServeInBlocking(func() {
		started <- 1
	})
	<-started
}

func newTestProxyService(filename string) (*ProxyEngine, error) {
	cfg, err := getCfgFromFile(filename)
	if err != nil {
		return nil, err
	}
	dummyService, err := NewProxyEngine("127.0.0.1:9696", "utf8mb4", true)
	if err != nil {
		return nil, err
	}
	err = dummyService.ParserConfig(cfg)
	if err != nil {
		return nil, err
	}
	return dummyService, nil
}

func newTestRouter(filename string) (*Router, error) {
	cfg, err := getCfgFromFile(filename)

	dummyService, err := NewProxyEngine("127.0.0.1:9696", "utf8mb4", true)
	if err := dummyService.ParseHostGroupNodes(cfg, false); err != nil {
		return nil, err
	}

	if err := dummyService.ParseHostGroupNodeClusters(cfg); err != nil {
		return nil, err
	}

	if err := dummyService.ParseSchemas(cfg); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return dummyService.GetRouter(), nil
}

func getCfgFromFile(filename string) (*config.Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := new(config.Config)

	if err = yaml.Unmarshal(file, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
