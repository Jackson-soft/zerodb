package server

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/agent"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"

	"github.com/golang/mock/gomock"
)

var (
	agentMockclients = []*mock_agent.MockAgentClient{}
)

func TestServer_isBinlogFullySynced(t *testing.T) {
	conf := &Config{
		AgentPort: "5002",
	}
	creatlog(conf, t)
	createAgentmock(t)
	store, err := createTestStore()
	if err != nil {
		t.Errorf("new store failed %v", err)
	}
	type fields struct {
		conf         *Config
		store        *etcdtool.Store
		proxyClients sync.Map
		agentClients sync.Map
		pushFailed   bool
	}
	type args struct {
		ctx             context.Context
		clusterName     string
		hostgroup       string
		from            int32
		to              int32
		SafeBinlogDalay float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				ctx:             context.TODO(),
				clusterName:     "cluster_test",
				hostgroup:       "hostGroup1",
				from:            0,
				to:              1,
				SafeBinlogDalay: 1000,
			},
			fields: fields{
				store: store,
				conf:  conf,
			},
			wantErr: false,
		},
		{
			name: "case2",
			args: args{
				ctx:             context.TODO(),
				clusterName:     "cluster_test",
				hostgroup:       "hostGroup1",
				from:            0,
				to:              1,
				SafeBinlogDalay: 1000,
			},
			fields: fields{
				store: store,
				conf:  conf,
			},
			wantErr: true,
		},
		{
			name: "case3",
			args: args{
				ctx:             context.TODO(),
				clusterName:     "cluster_test",
				hostgroup:       "hostGroup1",
				from:            0,
				to:              1,
				SafeBinlogDalay: 1000,
			},
			fields: fields{
				store: store,
				conf:  conf,
			},
			wantErr: true,
		},
	}
	tests[0].fields.agentClients.Store("10.1.22.1:5002", agentMockclients[0])
	tests[0].fields.agentClients.Store("10.1.21.79:5002", agentMockclients[1])
	tests[1].fields.agentClients.Store("10.1.22.1:5002", agentMockclients[2])
	tests[1].fields.agentClients.Store("10.1.21.79:5002", agentMockclients[3])
	tests[2].fields.agentClients.Store("10.1.22.1:5002", agentMockclients[4])
	tests[2].fields.agentClients.Store("10.1.21.79:5002", agentMockclients[5])
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				conf:         tt.fields.conf,
				store:        tt.fields.store,
				proxyClients: tt.fields.proxyClients,
				agentClients: tt.fields.agentClients,
				pushFailed:   tt.fields.pushFailed,
			}
			if err := s.isBinlogFullySynced(tt.args.ctx, tt.args.clusterName, tt.args.hostgroup, tt.args.from, tt.args.to, tt.args.SafeBinlogDalay); (err != nil) != tt.wantErr {
				t.Errorf("Server.isBinlogFullySynced() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_getSwitchConfig(t *testing.T) {
	store, err := createTestStore()
	if err != nil {
		t.Errorf("new store failed %v", err)
	}
	type fields struct {
		conf         *Config
		store        *etcdtool.Store
		proxyClients sync.Map
		agentClients sync.Map
		pushFailed   bool
	}
	type args struct {
		clusterName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    config.SwitchDB
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				store: store,
			},
			args: args{
				clusterName: "cluster_test",
			},
			want: config.SwitchDB{
				NeedLoadCheck:    true,
				NeedVote:         false,
				NeedBinlogCheck:  true,
				SafeBinlogDelay:  1000,
				VoteApproveRatio: 0,
				SafeLoad:         8,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				conf:         tt.fields.conf,
				store:        tt.fields.store,
				proxyClients: tt.fields.proxyClients,
				agentClients: tt.fields.agentClients,
				pushFailed:   tt.fields.pushFailed,
			}
			got, err := s.getSwitchConfig(tt.args.clusterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.getSwitchConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.getSwitchConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createAgentmock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i := 0; i < 6; i++ {
		agent := mock_agent.NewMockAgentClient(ctrl)
		agentMockclients = append(agentMockclients, agent)
	}
	agentMockclients[0].EXPECT().GetLoad(gomock.Any(), gomock.Any()).Return(new(basic.BasicResponse), nil).AnyTimes()
	agentMockclients[0].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file1",
		Position: 10000,
	}, nil).AnyTimes()
	agentMockclients[1].EXPECT().GetLoad(gomock.Any(), gomock.Any()).Return(new(basic.BasicResponse), nil).AnyTimes()
	agentMockclients[1].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file1",
		Position: 10000,
	}, nil).AnyTimes()
	agentMockclients[2].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file1",
		Position: 10000,
	}, nil).AnyTimes()
	agentMockclients[3].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file2",
		Position: 10000,
	}, nil).AnyTimes()
	agentMockclients[4].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file1",
		Position: 10000,
	}, nil).AnyTimes()
	agentMockclients[5].EXPECT().GetBinLog(gomock.Any(), gomock.Any()).Return(&agent.BinLogResponse{
		File:     "file1",
		Position: 9500,
	}, nil).AnyTimes()
}

func createTestStore() (*etcdtool.Store, error) {
	endpoints := []string{"10.1.22.2:2379", "10.1.22.3:2379", "10.1.22.4:2379"}
	store, err := etcdtool.NewStore(endpoints)
	if err != nil {
		return nil, err
	}
	return store, nil
}
