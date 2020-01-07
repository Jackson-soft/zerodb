package models

type Config struct {
	Persistence []string `yaml:"persistence"`
	RPCServer   string   `yaml:"rpcServer"`
	WebServer   string   `yaml:"webServer"`
	AgentPort   string   `yaml:"agentPort"`
	ProxyPort   string   `yaml:"proxyPort"`
	MetricsAddr string   `yaml:"metricsAddr"`
	OssURL      string   `yaml:"ossUrl"`
	LogConf     struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"log"`
	DebugServer string `yaml:"debugServer"`
	DebugMode   bool   `yaml:"debugMode"`
}

// the status of proxy
const (
	StatusReady = "ready"
	StatusUp    = "up"
	StatusLost  = "lost"
	StatusUnreg = "unreg"
)

// 状态信息
type StatusInfo struct {
	Status string `json:"status"`
	Time   int64  `json:"time"`
}
