package server

import (
	"time"

	"git.2dfire.net/zerodb/common/zeroproto/pkg/proxy"
)

// proxy的客户端map
type ProxyCli map[string]proxy.ProxyClient

const (
	TimeOut                       = 3 * time.Second
	WebTimeOut                    = 20 * time.Second
	PushConfigTimeOut             = 100 * time.Second
	HeartbeatCheckInterval        = 2 * time.Second
	DefaultSwitchDBInterval int64 = 30
	DialTimeOut                   = 5 * time.Second
	MaxRetryNum             uint  = 3
	// StatusLostTime 判定丢失的心跳间隔
	HBLost = int64(3 * 5)
	// Lost后10秒内没有启动，判定为Down
	HBDown    = int64(HBLost + 10)
	Heartbeat = 2 * time.Second
	UnregNum  = 1
	MysqlKey  = "mysql"
	SystemKey = "system"

	UploadAPI = "/uploadfile"
)
