package server

import "time"

const (
	//默认rpc服务地址
	DefaultrpcServer = "0.0.0.0:5004"
	//默认Metrics服务地址
	DefaultMetricsServer = "0.0.0.0:5005"
	//默认的代理服务地址
	DefaultProxyServer = "0.0.0.0:5003"
	//默认配置文件地址
	DefaultConfigPath = "app.yaml"
	//默认日志级别
	DefaultLogLevel = "error"
	//默认日志路径
	DefaultLogPath = "/opt/logs/golang/proxy.log"
	//默认日志保存时间
	DefaultMaxLogDay = 7
	//TimeOut rpc请求超时时间
	TimeOut     = 3 * time.Second
	DialTimeOut = 5 * time.Second
	//HearbeatTime 心跳频率
	HearbeatTime      = 5 * time.Second
	MaxRetryNum  uint = 3
)

// the status of proxy
const (
	StatusReady   = "ready"
	StatusUp      = "up"
	StatusLost    = "lost"
	StatusUnreg   = "unreg"
	StatusOffline = "offline"
	StatusOnline  = "online"
)
