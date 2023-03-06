package config

// Config 分库分表配置
type Config struct {
	Version string
	Basic   Basic `yaml:"basic" json:"basic"`

	StopService StopService `yaml:"stop_service" json:"stop_service"`

	Switch SwitchDB `yaml:"switch" json:"switch"`

	HostGroups []HostGroup `yaml:"host_groups" json:"host_groups"`

	HostGroupClusters []HostGroupCluster `yaml:"host_group_clusters" json:"host_group_clusters"`

	SchemaConfigs []SchemaConfig `yaml:"schema_configs" json:"schema_configs"`
}

type Basic struct {
	ConfigName   string `yaml:"config_name" json:"config_name"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	SlowLogTime  int    `yaml:"slow_log_time" json:"slow_log_time"`
	PushWhenFail bool   `yaml:"push_when_fail" json:"push_when_fail"`
}

type SwitchDB struct {
	NeedVote            bool    `yaml:"need_vote" json:"need_vote"`
	VoteApproveRatio    float64 `yaml:"vote_approve_ratio" json:"vote_approve_ratio"`
	NeedLoadCheck       bool    `yaml:"need_load_check" json:"need_load_check"`
	SafeLoad            float64 `yaml:"safe_load" json:"safe_load"`
	NeedBinlogCheck     bool    `yaml:"need_binlog_check" json:"need_binlog_check"`
	SafeBinlogDelay     float64 `yaml:"safe_binlog_delay" json:"safe_binlog_delay"`
	Binlogwaittime      int     `yaml:"binlog_wait_time" json:"binlog_wait_time"`
	Frequency           int     `yaml:"frequency" json:"frequency"`
	BackendPingInterval int     `yaml:"backend_ping_interval" json:"backend_ping_interval"`
}

type StopService struct {
	OfflineOnLostKeeper   bool `yaml:"offline_on_lost_keeper" json:"offline_on_lost_keeper"` //跟多台mysql host断开后是否下线自己默认必须为true
	OfflineSwhRejectedNum int  `yaml:"offline_swh_rejected_num" json:"offline_swh_rejected_num"`
	OfflineDownHostNum    int  `yaml:"offline_down_host_num" json:"offline_down_host_num"`
	OfflineRecover        bool `yaml:"offline_recover" json:"offline_recover"`
}

type HostGroup struct {
	Name         string `yaml:"name" json:"name"`
	MaxConn      int    `yaml:"max_conn" json:"max_conn"`
	InitConn     int    `yaml:"init_conn" json:"init_conn"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	Write        string `yaml:"write" json:"write"`
	ActiveWrite  int    `yaml:"active_write" json:"active_write"`
	Read         string `yaml:"read,omitempty" json:"read,omitempty"`
	EnableSwitch bool   `yaml:"enable_switch" json:"enable_switch"`
}

type HostGroupCluster struct {
	Name                 string   `yaml:"name" json:"name"`
	NonshardingHostGroup string   `yaml:"nonsharding_host_group" json:"nonsharding_host_group"`
	ShardingHostGroups   []string `yaml:"sharding_host_groups" json:"sharding_host_groups"`
}

type TableConfig struct {
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	ShardingKey string `yaml:"sharding_key,omitempty" json:"sharding_key,omitempty"`
	Rule        string `yaml:"rule,omitempty" json:"rule,omitempty"`
}

type SchemaConfig struct {
	Name               string `yaml:"name" josn:"name"`
	Custody            bool   `yaml:"custody,omitempty" json:"custody,omitempty"`
	RwSplit            bool   `yaml:"rw_split" json:"rw_split"`
	HostGroupCluster   string `yaml:"host_group_cluster" json:"hostgroupcluster"`
	SchemaSharding     int    `yaml:"schema_sharding,omitempty" json:"schema_sharding,omitempty"`
	TableSharding      int    `yaml:"table_sharding,omitempty" json:"table_sharding,omitempty"`
	MultiRoute         bool   `yaml:"multi_route" json:"multi_route"`
	InitConnMultiRoute bool   `yaml:"init_conn_multi_route" json:"init_conn_multi_route"`

	TableConfigs []TableConfig `yaml:"table_configs,omitempty" json:"table_configs,omitempty"`
}
