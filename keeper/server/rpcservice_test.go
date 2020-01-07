package server

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"

	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"github.com/golang/mock/gomock"
)

var proxyMockClients []*mock_proxy.MockProxyClient

func TestServer_SwitchDB(t *testing.T) {
	conf := &Config{}
	creatlog(conf, t)
	store, err := createTestStore()
	if err != nil {
		t.Fatal(err)
	}
	createAgentmock(t)
	createProxyMock(t)

	type fields struct {
		conf         *Config
		store        *etcdtool.Store
		proxyClients sync.Map
		agentClients sync.Map
		pushFailed   bool
	}
	type args struct {
		ctx context.Context
		in  *keeper.SwitchDBRequest
	}
	a := args{
		in: &keeper.SwitchDBRequest{
			ClusterName: "cluster_test",
			Hostgroup:   "hostGroup1",
			From:        0,
			To:          1,
			ProxyIP:     "10.1.1.1",
		},
		ctx: context.TODO(),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.SwitchDBResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "case1",
			args:    a,
			want:    createSwitchDBResp(codes.OK),
			wantErr: false,
			fields: fields{
				store: store,
			},
		},
	}
	conf.AgentPort = "5002"
	tests[0].fields.conf = conf
	tests[0].fields.agentClients.Store("10.1.22.1:5002", agentMockclients[0])
	tests[0].fields.agentClients.Store("10.1.21.79:5002", agentMockclients[1])
	cli := make(ProxyCli)
	cli["10.1.1.1"] = proxyMockClients[0]
	tests[0].fields.proxyClients.Store("cluster_test", cli)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				conf:         tt.fields.conf,
				store:        tt.fields.store,
				proxyClients: tt.fields.proxyClients,
				agentClients: tt.fields.agentClients,
				pushFailed:   tt.fields.pushFailed,
			}
			got, err := s.SwitchDB(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.SwitchDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.SwitchDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

//create Mock  proxy object
func createProxyMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i := 0; i < 3; i++ {
		mockproxy := mock_proxy.NewMockProxyClient(ctrl)
		proxyMockClients = append(proxyMockClients, mockproxy)
	}
	proxyMockClients[0].EXPECT().GetVote(gomock.Any(), gomock.Any()).Return(&proxy.GetVoteResponse{
		BasicResp: new(basic.BasicResponse),
		From:      0,
		Vote:      true,
	}, nil).AnyTimes()
	proxyMockClients[0].EXPECT().RecoverWritingAbility(gomock.Any(), gomock.Any()).Return(new(basic.BasicResponse), nil).AnyTimes()
	proxyMockClients[0].EXPECT().StopWritingAbility(gomock.Any(), gomock.Any()).Return(new(basic.BasicResponse), nil).AnyTimes()
	proxyMockClients[0].EXPECT().SwitchDataSource(gomock.Any(), gomock.Any()).Return(&proxy.SwitchDatasourceResponse{new(basic.BasicResponse)}, nil).AnyTimes()
}

func creatlog(conf *Config, t *testing.T) {
	conf.LogConf.Level = "info"
	conf.LogConf.Path = "zlog"
	err := glog.CreateLogs("text", conf.LogConf.Path, conf.LogConf.Level)
	if err != nil {
		t.Errorf("create log %v", err)
		return
	}
}
