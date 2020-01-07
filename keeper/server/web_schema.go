package server

import (
	"context"
	"net/http"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"git.2dfire.net/zerodb/keeper/pkg/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) addSchema(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var parameter addSchema
	var err error
	if err = c.BindJSON(&parameter); err != nil {
		glog.GLog.Errorf("add schemas error: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	// 过滤界面传来的无用的table信息
	for i := 0; i < len(parameter.Schemas); i++ {
		if parameter.Schemas[i].Custody {
			if len(parameter.Schemas[i].TableConfigs) > 0 {
				msg.Code = codes.Failed
				msg.ErrorCode = codes.CommonError
				msg.ErrMsg = "parameter error."
				return
			}
		} else {
			for j := 0; j < len(parameter.Schemas[i].TableConfigs); j++ {
				if parameter.Schemas[i].TableConfigs[j].Name == "" || parameter.Schemas[i].TableConfigs[j].Rule == "" || parameter.Schemas[i].TableConfigs[j].ShardingKey == "" {
					msg.Code = codes.Failed
					msg.ErrorCode = codes.CommonError
					msg.ErrMsg = parameterIsNil
					return
				}
			}
		}
	}

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + parameter.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(parameter.Cluster)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithField("cluster", parameter).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster).Errorln(err)
		return
	}

	var data []byte
	sMsgs := make([]string, 0)
	bSucceed := true

	for _, cfg := range parameter.Schemas {
		in := new(proxy.AddSchemaRequest)
		in.Version = nextVersion
		data, err = utils.UnLoadYaml(cfg)
		if err != nil {
			glog.GLog.Errorf("add schema unloadyaml error: %v", err)
			continue
		}
		in.SchemaCfgData = data

		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := cli.AddSchema(ctx, in)
			cancel()
			if err != nil {
				bSucceed = false
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithField("cluster", parameter.Cluster).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.AddSchema(parameter.Cluster, parameter.Schemas); err != nil {
			glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			return
		}
	}
	msg.Data = sMsgs

	glog.GLog.Infof("add schema successfully")
}

func (s *Server) deleteSchema(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var err error
	var parameter deleteByNames
	if err = c.BindJSON(&parameter); err != nil {
		glog.GLog.Errorf("delete schema failed: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + parameter.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(parameter.Cluster)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithField("cluster", parameter.Cluster).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster).Errorln(err)
		return
	}

	sMsgs := make([]string, 0)
	bSucceed := true

	in := new(proxy.DeleteSchemaRequest)
	in.Version = nextVersion

	for _, name := range parameter.Names {
		in.SchemaName = name
		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			defer cancel()
			resp, err := cli.DeleteSchema(ctx, in)
			if err != nil {
				glog.GLog.WithFields("ip", ip, "cluster", parameter.Cluster, "name", name).Errorln(err)
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("ip", ip, "cluster", parameter.Cluster, "name", name).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.DeleteSchema(parameter.Cluster, parameter.Names); err != nil {
			glog.GLog.WithFields("cluster", parameter.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = err.Error()
			return
		}
	}

	msg.Data = sMsgs

	glog.GLog.Infof("delete schema successfully")
}

func (s *Server) addTable(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var err error
	var parameter addTable

	if err = c.BindJSON(&parameter); err != nil {
		glog.GLog.Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + parameter.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(parameter.Cluster)
	if err != nil {
		msg.Code = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithField("cluster", parameter.Cluster).Errorln(err)
		return
	}
	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithField("cluster", parameter.Cluster).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster).Errorln(err)
		return
	}

	var data []byte
	in := new(proxy.AddTableRequest)
	in.SchemaName = parameter.SchemaName
	in.Version = nextVersion

	sMsgs := make([]string, 0)
	bSucceed := true

	for _, value := range parameter.Tables {
		data, err = utils.UnLoadYaml(value)
		if err != nil {
			glog.GLog.Errorf("add table unloadyaml error: %v", err)
			continue
		}
		in.TableCfgData = data

		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := cli.AddTable(ctx, in)
			cancel()
			if err != nil {
				glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.SchemaName).Errorln(err)
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.SchemaName).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.AddTable(parameter.Cluster, parameter.SchemaName, parameter.Tables); err != nil {
			glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.SchemaName).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = err.Error()
			return
		}
	}

	msg.Data = sMsgs

	glog.GLog.Infof("add table successfully")
}

func (s *Server) deleteTable(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var err error
	var parameter deleteTable
	if err = c.BindJSON(&parameter); err != nil {
		glog.GLog.Errorf("delete table failed: %v", err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	clis, ok := s.proxyClients.Load(parameter.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + parameter.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(parameter.Cluster)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.Schema, "table", parameter.TableName).Errorln(err)
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.Schema, "table", parameter.TableName).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(parameter.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.Schema, "table", parameter.TableName).Errorln(err)
		return
	}

	sMsgs := make([]string, 0)
	bSucceed := true
	in := new(proxy.DeleteTableRequest)
	in.SchemaName = parameter.Schema
	in.Version = nextVersion
	in.TableName = parameter.TableName
	for ip, cli := range clis.(ProxyCli) {
		ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
		resp, err := cli.DeleteTable(ctx, in)
		cancel()
		if err != nil {
			glog.GLog.WithFields("cluster", parameter.Cluster, "ip", ip, "schema", parameter.Schema, "table", parameter.TableName).Errorln(err)
			bSucceed = false
			sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
		} else {
			sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
			if resp.Code != codes.OK {
				bSucceed = false
				glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.Schema, "table", parameter.TableName, "ip", ip).Errorln(resp.ErrMsg)
			}
		}
	}

	if bSucceed {
		if err = s.store.DeleteTable(parameter.Cluster, parameter.Schema, parameter.TableName); err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = err.Error()
			glog.GLog.WithFields("cluster", parameter.Cluster, "schema", parameter.Schema, "table", parameter.TableName).Errorln(err)
			return
		}
	}
	msg.Data = sMsgs

	glog.GLog.Infof("delete table successfully")
}

func (s *Server) updateShardRW(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	var err error
	var shard shardRW

	if err = c.BindJSON(&shard); err != nil {
		glog.GLog.Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		return
	}

	glog.GLog.WithFields("clusterName", shard.Cluster).Infoln("update shard read write.")

	clis, ok := s.proxyClients.Load(shard.Cluster)
	if !ok || len(clis.(ProxyCli)) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = proxyNotExist + " [" + shard.Cluster + "]"
		return
	}

	status, err := s.store.GetClusterStatus(shard.Cluster)
	if err != nil {
		glog.GLog.WithFields("cluster", shard.Cluster).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	if status != models.StatusUp {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = statusError + "expected: [" + models.StatusUp + "] actual: [" + status + "]"
		glog.GLog.WithField("cluster", shard.Cluster).Warnln(statusError)
		return
	}

	nextVersion, err := s.store.GetNextVersion(shard.Cluster, false)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", shard.Cluster).Errorln(err)
		return
	}

	in := new(proxy.UpdateSchemaRWSplitRequest)
	in.Version = nextVersion
	sMsgs := make([]string, 0)
	bSucceed := true

	for key, value := range shard.ShardRW {
		in.SchemaName = key
		in.RwSplit = value
		for ip, cli := range clis.(ProxyCli) {
			ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
			resp, err := cli.UpdateSchemaRWSplit(ctx, in)
			cancel()
			if err != nil {
				sMsgs = append(sMsgs, rpcMessage(ip, err.Error()))
				bSucceed = false
				glog.GLog.WithFields("cluster", shard.Cluster, "ip", ip, "schema", key).Errorln(err)
			} else {
				sMsgs = append(sMsgs, rpcMessage(ip, resp.ErrMsg))
				if resp.Code != codes.OK {
					bSucceed = false
					glog.GLog.WithFields("cluster", shard.Cluster, "ip", ip, "schema", key).Errorln(resp.ErrMsg)
				}
			}
		}
	}

	if bSucceed {
		if err = s.store.UpdateShardRW(shard.Cluster, shard.ShardRW); err != nil {
			glog.GLog.WithFields("clusterName", shard.Cluster).Errorln(err)
			msg.Code = codes.Failed
			msg.ErrorCode = codes.EtcdPutFailed
			msg.ErrMsg = inconsistentError + err.Error()
			return
		}
	}
	msg.Data = sMsgs

	glog.GLog.Infof("update read/write split successfully")
}
