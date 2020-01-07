package models

//Config 定义了从app.yaml 中读取的启动配置文件
type Config struct {
	ClusterName string `yaml:"clusterName"`
	RPCServer   string `yaml:"rpcServer"`
	KeeperAddr  string `yaml:"keeperAddr"`
	MetricsAddr string `yaml:"metricsAddr"`
	ProxyServer string `yaml:"proxyServer"`
	Charset     string `yaml:"charset"`
	DebugServer string `yaml:"debugServer"`
	DebugMode   bool   `yaml:"debugMode"`
	LogSQL      bool   `yaml:"logSQL"`
	LogConf     struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`
}
