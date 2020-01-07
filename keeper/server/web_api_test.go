package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/keeper/pkg/glog"

	"github.com/golang/mock/gomock"

	"github.com/pkg/errors"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
)

var dataConf = `
persistence: [10.1.134.71:22379]
rpcAddr: 127.0.0.1:5006
keeperAddr: 127.0.0.1:5003
metricsAddr: 0.0.0.0:5008
proxyAddr: 0.0.0.0:9696
charset: utf8mb4
log:
  file: zlog/log
  level: info
`

func InitServer(req *http.Request) (interface{}, error) {
	svr, err := creatSvr()
	if err != nil {
		return nil, err
	}
	return checkData(svr, req)
}

func creatSvr() (*Server, error) {
	var err error
	conf := Config{}
	if err = utils.LoadYaml([]byte(dataConf), &conf); err != nil {
		return nil, err
	}

	if err = glog.CreateLogs("text", conf.LogConf.Path, conf.LogConf.Level); err != nil {
		return nil, err
	}

	svr, err := NewServer(&conf)
	if err != nil {
		return nil, err
	}

	return svr, nil
}

func creatMockSvr(t *testing.T) (*Server, []*mock_proxy.MockProxyClient, error) {
	svr, err := creatSvr()
	if err != nil {
		return nil, nil, err
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	var proxys []*mock_proxy.MockProxyClient

	keys := []string{"10.1.21.80:5006", "10.1.21.82:5006"}

	for _, key := range keys {
		mockProxy := mock_proxy.NewMockProxyClient(mockCtl)

		cli := make(ProxyCli)
		cli[key] = mockProxy

		svr.proxyClients.Store("cluster_test", cli)

		proxys = append(proxys, mockProxy)
	}
	return svr, proxys, nil
}

func checkData(svr *Server, req *http.Request) (interface{}, error) {
	w := httptest.NewRecorder()

	svr.creatHander().ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := APIResult{}
	if err = utils.LoadJSON(body, &ret); err != nil {
		return nil, err
	}

	if ret.Code != codes.OK {
		return ret.Data, errors.New(ret.ErrMsg)
	}
	return ret.Data, nil
}

func TestLogin(t *testing.T) {
	v := url.Values{}
	v.Set("userName", "admin")
	v.Set("password", "admin")
	req, err := http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	data, err := InitServer(req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/proxy_cluster/status", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("clusterName", "cluster_test")
	req.URL.RawQuery = q.Encode()

	data, err := InitServer(req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestInitConfig(t *testing.T) {
	bodyBuf := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuf)

	fileName := "../../zero-proxy/proxy/test-conf/proxy_conf_test.yaml"
	file, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}

	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		t.Error(err)
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		t.Error(err)
	}

	contentType := bodyWriter.FormDataContentType()

	if err = bodyWriter.WriteField("clusterName", "cluster_test"); err != nil {
		t.Error(err)
	}
	if err = bodyWriter.WriteField("force", "true"); err != nil {
		t.Error(err)
	}

	bodyWriter.Close()

	req, err := http.NewRequest("POST", "/api/shardconf_init", &bodyBuf)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", contentType)

	data, err := InitServer(req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)

}

func TestGetConfig(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/config", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("clusterName", "cluster_test")
	q.Add("snapshotName", "")
	req.URL.RawQuery = q.Encode()

	data, err := InitServer(req)
	if err != nil {
		t.Error(err)
	}
	cfg := data.(map[string]interface{})
	t.Log(cfg["user"])
}

func TestSwitch(t *testing.T) {
	req, err := http.NewRequest("PUT", "/api/proxy_cluster/switch", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("hostGroup", "hostGroup1")
	q.Add("from", "1")
	q.Add("to", "2")
	q.Add("clusterName", "cluster_test")
	q.Add("reason", "this is a reason")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	svr, err := creatSvr()
	if err != nil {
		t.Error(err)
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockProxy := mock_proxy.NewMockProxyClient(mockCtl)

	mockProxy.EXPECT().SwitchDataSource(gomock.Any(), gomock.Any()).Return(
		&proxy.SwitchDatasourceResponse{}, nil)

	cli := make(ProxyCli)
	cli["10.1.21.80:5006"] = mockProxy

	svr.proxyClients.Store("cluster_test", cli)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestPushConfig(t *testing.T) {
	req, err := http.NewRequest("PUT", "/api/push_config", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("snapshotName", "")
	q.Add("clusterName", "cluster_test")

	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	svr, err := creatSvr()
	if err != nil {
		t.Error(err)
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockProxy := mock_proxy.NewMockProxyClient(mockCtl)

	cli := make(ProxyCli)
	cli["10.1.21.80:5006"] = mockProxy

	svr.proxyClients.Store("cluster_test", cli)

	mockProxy.EXPECT().PushConfig(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code: codes.OK,
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestAddHosts(t *testing.T) {
	payload := strings.NewReader("{\n    \"clusterName\": \"cluster_test\",\n    \"groups\": [{\n        \"name\": \"hostGroup100\",\n        \"max_conn\": 100,\n        \"init_conn\": 100,\n        \"user\": \"zerodb\",\n        \"password\": \"zerodb@552208\",\n        \"write\": \"10.1.22.1:3306,10.1.21.79:3306\"\n    }]\n}")

	req, err := http.NewRequest("POST", "/api/add_hostgroups", payload)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().AddHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().AddHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.CommonError,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestDeleHosts(t *testing.T) {
	p := deleteByNames{
		Cluster: "cluster_test",
		Names:   []string{"hostGroup6", "hostGroup100"},
	}

	b, err := utils.UnLoadJSON(p)
	if err != nil {
		t.Error(err)
	}

	payload := strings.NewReader(string(b))

	req, err := http.NewRequest("POST", "/api/delete_hostgroups", payload)
	if err != nil {
		t.Error(err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().DeleteHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().DeleteHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.OK,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func createPostReq(data interface{}, url string) (*http.Request, error) {
	b, err := utils.UnLoadJSON(data)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	payload := strings.NewReader(string(b))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	return req, nil
}

func TestStopService(t *testing.T) {
	p := stopService{
		Cluster: "cluster_test",
		Service: config.StopService{
			UnregOnSwhRejected:    true,
			UnregRejectedNum:      4,
			UnregDownHostNum:      3,
			OfflineOnLostKeeper:   true,
			OfflineSwhRejectedNum: 4,
			OfflineDownHostNum:    3,
			OfflineRecover:        false,
		},
	}

	req, err := createPostReq(p, "/api/update_stopservice")
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().DeleteHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().DeleteHostGroup(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.OK,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestAddSchema(t *testing.T) {
	p := addSchema{
		Cluster: "cluster_test",
		Schemas: []config.SchemaConfig{
			{
				Name:                 "xxx_test",
				Custody:              true,
				NonshardingHostGroup: "hostGroup1",
				ShardingHostGroups:   []string{"hostGroup1", "hostGroup2", "hostGroup3", "hostGroup4"},
				SchemaSharding:       64,
				TableSharding:        2,
				TableConfigs: []config.TableConfig{
					{
						Name:        "ship",
						ShardingKey: "entity_id",
						Rule:        "string",
					},
				},
			},
		},
	}

	req, err := createPostReq(p, "/api/add_schema")
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().AddSchema(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().AddSchema(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.CommonError,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestDeleSchema(t *testing.T) {
	p := deleteByNames{
		Cluster: "cluster_test",
		Names:   []string{"xxx_test", "member"},
	}

	req, err := createPostReq(p, "/api/delete_schema")
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().DeleteSchema(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().DeleteSchema(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.CommonError,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestAddTable(t *testing.T) {
	p := addTable{
		Cluster:    "cluster_test",
		SchemaName: "empty",
		Tables: []config.TableConfig{
			{
				Name:        "xxx_text",
				ShardingKey: "shard_key",
				Rule:        "hash",
			},
		},
	}

	req, err := createPostReq(p, "/api/add_table")
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().AddTable(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().AddTable(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.OK,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestDeleTable(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/delete_table", nil)
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().DeleteTable(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().DeleteTable(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.CommonError,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}

func TestShardRW(t *testing.T) {
	p := shardRW{
		Cluster: "cluster_test",
		ShardRW: map[string]bool{
			"hostGroup1": true,
			"hostGroup2": false,
		},
	}

	req, err := createPostReq(p, "/api/update_shardrw")
	if err != nil {
		t.Error(err)
	}

	svr, proxys, err := creatMockSvr(t)
	if err != nil {
		t.Error(err)
	}

	proxys[0].EXPECT().UpdateSchemaRWSplit(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{}, nil)

	proxys[1].EXPECT().UpdateSchemaRWSplit(gomock.Any(), gomock.Any()).Return(
		&basic.BasicResponse{
			Code:   codes.CommonError,
			ErrMsg: "this is a error msg",
		}, nil)

	data, err := checkData(svr, req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v\n", data)
}
