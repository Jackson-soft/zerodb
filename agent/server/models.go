package server

//Config 配置项
type Config struct {
	KeeperAddr    string  `yaml:"keeperAddr"`
	RPCServer     string  `yaml:"rpcServer"`
	MetricsServer string  `yaml:"metricsServer"`
	MysqlConn     string  `yaml:"mysql"`
	LogConf       LogConf `yaml:"log"`
}

//LogConf 日志相关配置项
type LogConf struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}
