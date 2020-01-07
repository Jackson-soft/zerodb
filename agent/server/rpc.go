package server

import (
	"context"

	"github.com/pkg/errors"

	"git.2dfire.net/zerodb/agent/glog"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/system/load"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/agent"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/basic"
)

func (s *Server) GetBinLog(ctx context.Context, in *agent.BinLogRequest) (*agent.BinLogResponse, error) {
	sql := "show master status"
	if in.Role == "slave" {
		sql = "show slave status"
	}

	mData, err := s.sqlQuery(sql)
	if err != nil {
		glog.Glog.Errorf("agent get binlog failed: %v\n", err)
		return nil, err
	}

	if mData != nil {
		out := new(agent.BinLogResponse)
		if in.Role == "slave" {
			out.File = mData["Relay_Master_Log_File"].(string)
			out.Position = mData["Exec_Master_Log_Pos"].(int64)
		} else {
			out.File = mData["File"].(string)
			out.Position = mData["Position"].(int64)
		}
		return out, nil
	}

	glog.Glog.Infof("binglog sql: %s, data: %v\n", sql, mData)

	return nil, errors.New("binglog data is nil")
}

func (s *Server) GetLoad(ctx context.Context, in *agent.LoadRequest) (*basic.BasicResponse, error) {
	avg, err := load.LoadAvg()
	if err != nil {
		glog.Glog.Errorf("agent get system load failed: %+v", err)
		return nil, err
	}

	out := new(basic.BasicResponse)
	if in.LoadLimit >= avg.Avg1min {
		glog.Glog.Infoln("agent loadlimit: %d, avg1min: %d, avg5min: %d, avg15min: %d", in.LoadLimit, avg.Avg1min, avg.Avg5min, avg.Avg15min)
		out.Code = codes.OK
	} else {
		out.Code = codes.SystemLoadHigh
		out.ErrMsg = "agent system load high"
		glog.Glog.Warnf("agent loadlimit: %d, avg1min: %d, avg5min: %d, avg15min: %d", in.LoadLimit, avg.Avg1min, avg.Avg5min, avg.Avg15min)
	}

	return out, nil
}
