package etcdtool

//this is etcd key
const (
	ProxyPrefix = "zerodb/proxy/status/"
	AgentPrefix = "zerodb/agent/status/"
	//clusterName+ip:port 集群信息
	ProxyStatus = "zerodb/proxy/status/%s/%s"
	//clusterName
	ProxyList = "zerodb/proxy/list/%s"
	//clusterName
	StatusPrefix = "zerodb/proxy/status/%s/"
	//ip:port
	AgentStatus = "zerodb/agent/status/%s"
	//ip:port + name
	AgentInfor = "zerodb/agent/infor/%s/%s"
	//clusterName+snapshot+key
	ShardConfig = "zerodb/shardconf/%s/%s/%s"
	//clusterName+snapshot
	ShardPrefix = "zerodb/shardconf/%s/%s/"
	//clusterName
	ShardList = "zerodb/shardconf/%s/list"
	//clusterName + ip:port
	ProxyInfor = "zerodb/proxy/infor/%s/%s"
	//clusterName
	ProxyInforPrefix = "zerodb/proxy/infor/%s/"
	//cluster名称列表
	ClusterList = "zerodb/proxy/cluster/list"
	//心跳锁 clusterName+address(ip:port)
	HeartbeatLock = "zerodb/proxy/heartbeat/%s/%s"
	// 配置当前版本号1.9.2 clusterName
	ConfigVersion = "zerodb/shardconf/version/%s"
	// 集群的配置版本号 clusterName+address(ip:port)
	ProxyVersion = "zerodb/proxyconf/version/%s/%s"
	// 版本号前缀
	ProxyVersionPrefix = "zerodb/proxyconf/version/%s/"

	//分布式锁前缀
	HostLockPrefix = "zerodb/keeper/dislock/"
	//切换频率控制
	SwitchDBIntervalPrefix = "zerodb/keeper/interval/"
	//所有proxyIP的k前缀
	ProxyIPPrefix = "zerodb/proxy/ip/"
	//备份回滚数据前缀
	LastModifyPrefix = "zerodb/rollback/%s/"
)

// 配置子项key
const (
	DefaultConfig     = "default"
	Basic             = "basic"
	StopService       = "stopservice"
	Switchdb          = "switch"
	HostGroups        = "hostgroups"
	HostGroupClusters = "hostgroupclusters"
	SchemaConfigs     = "schemaconfigs"
)
