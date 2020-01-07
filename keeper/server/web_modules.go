package server

import (
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/statics/codes"
)

type APIResult struct {
	Code      uint32      `json:"code"`
	ErrorCode uint32      `json:"errorCode"`
	ErrMsg    string      `json:"message"`
	Data      interface{} `json:"data"`
}

func NewResult() APIResult {
	return APIResult{
		Code:      codes.Succeed,
		ErrorCode: codes.OK,
		ErrMsg:    "",
		Data:      nil,
	}
}

type systemInfo struct {
	IP          string  `json:"ip"`
	Status      string  `json:"status"`
	ConfVersion string  `jsong:"confVersion"`
	CpuLoad     float64 `json:"cpuLoad,omitempty"`
	MemLoad     float64 `json:"memLoad,omitempty"`
	LoadAvg     struct {
		Avg1Min  float64 `json:"avg1min,omitempty"`
		Avg5Min  float64 `json:"avg5min,omitempty"`
		Avg15Min float64 `json:"avg15min,omitempty"`
	} `json:"loadAvg,omitempty"`
}

type shardRW struct {
	Cluster string          `json:"clusterName"`
	ShardRW map[string]bool `json:"shardrw"`
}

type addHost struct {
	Cluster string             `json:"clusterName"`
	Groups  []config.HostGroup `json:"groups"`
}

type deleteByNames struct {
	Cluster string   `json:"clusterName"`
	Names   []string `json:"names"`
}

type stopService struct {
	Cluster string             `json:"clusterName"`
	Service config.StopService `json:"service"`
}

type addSchema struct {
	Cluster string                `json:"clusterName"`
	Schemas []config.SchemaConfig `json:"schemas"`
}

type addTable struct {
	Cluster    string               `json:"clusterName"`
	SchemaName string               `json:"schemaName"`
	Tables     []config.TableConfig `json:"tables"`
}

type deleteTable struct {
	Cluster   string `json:"clusterName"`
	Schema    string `json:"schemaName"`
	TableName string `json:"tableName"`
}

type addhostGroupCluster struct {
	Cluster   string                  `json:"clusterName" binding:"required"`
	Hgcluster config.HostGroupCluster `json:"hostgroupcluster" binding:"required"`
}

type delHostGroupCluster struct {
	Cluster               string   `json:"clusterName" binding:"required"`
	HostGroupClusterNames []string `json:"hostgroupclusterName" binding:"required"`
	SnapsHot              string   `json:"snapshot"`
}

type updateBasic struct {
	Cluster string       `json:"clusterName"`
	Basic   config.Basic `json:"basic"`
}

type updateSwitch struct {
	Cluster string          `json:"clusterName"`
	Switch  config.SwitchDB `json:"switch"`
}

type unregProxy struct {
	Cluster string `json:"clusterName" binding:"required"`
	Address string `json:"address" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

type cluster struct {
	Cluster string `json:"clusterName" binding:"required"`
}

type SwitchDataSource struct {
	Cluster   string `json:"clusterName" binding:"required"`
	HostGroup string `json:"hostGroup" binding:"required"`
	From      *int   `json:"from"  binding:"exists"`
	To        *int   `json:"to" binding:"exists"`
	Reason    string `json:"reason" binding:"required"`
}
