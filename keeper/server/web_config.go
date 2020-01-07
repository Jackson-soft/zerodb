package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	"github.com/gin-gonic/gin"
)

// 查看配置
func (s *Server) apiConfig(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	clusterName := c.Query("clusterName")
	if len(clusterName) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	name := c.Query("snapshotName")

	data, err := s.store.GetShardConf(clusterName, name)
	if err != nil {
		glog.GLog.WithFields("clusterName", clusterName, "snapshotName", name).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		return
	}

	cfg := config.Config{}
	if err = utils.LoadYaml(data, &cfg); err != nil {
		glog.GLog.WithFields("clusterName", clusterName, "snapshotName", name).Errorln(err)
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		return
	}

	msg.Data = cfg
}

// basic
func (s *Server) apiBasic(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.Basic)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := config.Basic{}
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

// apiStopService
func (s *Server) apiStopService(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.StopService)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := config.StopService{}
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

// apiSwitch
func (s *Server) apiSwitch(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.Switchdb)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := config.SwitchDB{}
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

//获取hostgroups
func (s *Server) apiHostGroups(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.HostGroups)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := make([]config.HostGroup, 0)
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

// apiHostGroupClusters
func (s *Server) apiHostGroupClusters(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.HostGroupClusters)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := make([]config.HostGroupCluster, 0)
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

// apiSchemaConfigs
func (s *Server) apiSchemaConfigs(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.SchemaConfigs)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	conf := make([]config.SchemaConfig, 0)
	if err = utils.LoadYaml(data, &conf); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.YamlParseFail
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	msg.Data = conf
}

// apiTableConfigs
func (s *Server) apiTableConfigs(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")
	schema := c.Query("schemaName")

	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}
	if len(schema) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [schemaName]"
		return
	}
	if len(snapshot) == 0 {
		snapshot = etcdtool.DefaultConfig
	}

	key := fmt.Sprintf(etcdtool.ShardConfig, cluster, snapshot, etcdtool.SchemaConfigs)
	data, err := s.store.GetData(key)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot, "schema", schema).Errorln(err)
		return
	}

	schemas := make([]config.SchemaConfig, 0)
	if err := utils.LoadYaml(data, &schemas); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot, "schema", schema).Errorln(err)
		return
	}

	var tablesResult []config.TableConfig
	for _, s := range schemas {
		if s.Name == schema {
			tablesResult = s.TableConfigs
			break
		}
	}

	msg.Data = tablesResult
}

// 导出配置
func (s *Server) exportConfig(c *gin.Context) {
	msg := NewResult()
	defer c.JSON(http.StatusOK, &msg)

	cluster := c.Query("clusterName")
	snapshot := c.Query("snapshotName")
	if len(cluster) == 0 {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.ParameterIsNil
		msg.ErrMsg = parameterIsNil + " [clusterName]"
		return
	}

	exportName := c.Query("exportName")

	data, err := s.store.GetShardConf(cluster, snapshot)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}

	if len(exportName) == 0 {
		exportName = etcdtool.DefaultConfig
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", exportName+".yaml")
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	if _, err = io.Copy(fileWriter, bytes.NewReader(data)); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.EtcdGetFailed
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}

	params := map[string]string{
		"path":        "zerodb",
		"rename":      exportName,
		"projectName": "OssDownload",
		"domain":      "1",
	}
	for key, val := range params {
		bodyWriter.WriteField(key, val)
	}
	bodyWriter.Close()

	resp, err := http.Post(s.conf.OssURL+UploadAPI, bodyWriter.FormDataContentType(), bodyBuf)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
	if err = utils.LoadJSON(body, &msg); err != nil {
		msg.Code = codes.Failed
		msg.ErrorCode = codes.CommonError
		msg.ErrMsg = err.Error()
		glog.GLog.WithFields("cluster", cluster, "snapshot", snapshot).Errorln(err)
		return
	}
}
