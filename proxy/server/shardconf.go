package server

import (
	"encoding/json"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/glog"

	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var ErrBackupVersionIsNil = errors.New("backup config is nil")

//ShardConfig 为Proxy提供配置数据
type ShardConfig struct {
	mutex       sync.RWMutex
	currentConf *config.Config
	//oldConf     *config.Config
	oldConf *config.Config
}

//NewShardConf new
func NewShardConf() *ShardConfig {
	c := new(ShardConfig)
	c.currentConf = new(config.Config)
	//c.oldConf = new(config.Config)
	c.oldConf = new(config.Config)
	return c
}

//LoadConfig 加载配置文件，如已有则备份现有的配置
func (s *ShardConfig) LoadConfig(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//启动初始化配置

	if err := yaml.Unmarshal(data, s.currentConf); err != nil {
		return err
	}
	s.oldConf = nil
	return nil
}

//Rollback 回滚配置
func (s *ShardConfig) Rollback() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	/*	if s.oldConf.Version == s.currentConf.Version {
		//return ErrBackupVersionIsNil
		return nil
	}*/
	if s.oldConf == nil {
		return ErrBackupVersionIsNil
	}
	s.currentConf = s.oldConf
	s.oldConf = nil
	/*	if len(s.oldConf.Version) != 0 {
		s.currentConf = s.oldConf
		return nil
	}*/
	return nil
}

//GetConfig 获取当前的配置
func (s *ShardConfig) GetConfig() *config.Config {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.currentConf
}

//upddate shardconfig version

func (s *ShardConfig) SetVersion(version string) {
	s.mutex.Lock()
	s.currentConf.Version = version
	s.mutex.Unlock()
}

func (s *ShardConfig) GetVersion() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.currentConf.Version
}

type SyncData struct {
	Stype string
	Data  []byte
}

func (s *ShardConfig) DelHostGroup(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	hgs := s.currentConf.HostGroups
	for i, hg := range hgs {
		if hg.Name == name {
			s.currentConf.HostGroups = append(hgs[:i], hgs[i+1:]...)
			glog.Glog.Infof("del hostgroup %s in current shardconfig", name)
			return
		}
	}
	glog.Glog.Errorf("can't find  hostgroup %s in current shardconfig", name)
}

func (s *ShardConfig) AddHostGroup(hg config.HostGroup) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currentConf.HostGroups = append(s.currentConf.HostGroups, hg)
	glog.Glog.Infof("add hostgroup %s to current shardconfig", hg.Name)
}

func (s *ShardConfig) AddSchema(schemaConfig config.SchemaConfig) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.currentConf.SchemaConfigs = append(s.currentConf.SchemaConfigs, schemaConfig)
	glog.Glog.Infof("add schemaConfig %s to current shardconfig", schemaConfig.Name)
}

func (s *ShardConfig) DelSchema(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	schemaConfigs := s.currentConf.SchemaConfigs
	//这里对slice的append操作需要注意
	for i, sch := range schemaConfigs {
		if sch.Name == name {
			s.currentConf.SchemaConfigs = append(s.currentConf.SchemaConfigs[:i], s.currentConf.SchemaConfigs[i+1:]...)
			glog.Glog.Infof("del schema  %s from current shardconfig", name)
			return
		}
	}
	glog.Glog.Errorf("can't find  schema %s in current shardconfig", name)
}
func (s *ShardConfig) AddTable(name string, table config.TableConfig) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//check duplicate
	for _, sc := range s.currentConf.SchemaConfigs {
		if sc.Name == name {
			for _, t := range sc.TableConfigs {
				if t.Name == table.Name {
					return
				}
			}
		}
	}
	schemaConfigs := s.currentConf.SchemaConfigs
	for _, sch := range schemaConfigs {
		if sch.Name == name {
			sch.TableConfigs = append(sch.TableConfigs, table)
			glog.Glog.Infof("add %s tableConfig  to current shardconfig", name)
			return
		}
	}
	glog.Glog.Errorf("can't find  schema %s in current shardconfig", name)
}

func (s *ShardConfig) DelTable(name string, tablename string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	schemaConfigs := s.currentConf.SchemaConfigs
	for _, sch := range schemaConfigs {
		if sch.Name == name {
			for i, t := range sch.TableConfigs {
				if t.Name == tablename {
					sch.TableConfigs = append(sch.TableConfigs[:i], sch.TableConfigs[i+1:]...)
					glog.Glog.Infof("del schema %s table %s from current shardconfig", name, tablename)
					return
				}
			}
		}
	}
	glog.Glog.Errorf("can't find  schema %s table %s  in current shardconfig", name, tablename)
}

func (s *ShardConfig) UpdateRwSplit(schemaName string, rwSplit bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	schemaConfigs := s.currentConf.SchemaConfigs
	var index int
	for i, sch := range schemaConfigs {
		if sch.Name == schemaName {
			index = i + 1
			break
		}
	}
	if index == 0 {
		return
	}
	schemaConfigs[index-1].RwSplit = rwSplit
}

func (s *ShardConfig) BackupCurConfig() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.currentConf == nil {
		return errors.Errorf("current config is nil")
	}
	glog.Glog.Infoln("backup current config ", s.currentConf.Version)
	//c = *s.currentConf
	//s.oldConf = &c
	s.oldConf = &config.Config{}
	if err := deepCopy(s.oldConf, s.currentConf); err != nil {
		return err
	}
	return nil
}

//print shardconfig info for debug
func (s *ShardConfig) Print() {
	glog.Glog.Infoln("current config ")
	PrettyPrint(s.currentConf)
	glog.Glog.Infoln("backup config")
	PrettyPrint(s.oldConf)
}
func (s *ShardConfig) AddHostGroupCluster(cluster config.HostGroupCluster) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//check deplicate
	hostgroupclusters := s.currentConf.HostGroupClusters
	for _, h := range hostgroupclusters {
		if h.Name == cluster.Name {
			return
		}
	}
	s.currentConf.HostGroupClusters = append(s.currentConf.HostGroupClusters, cluster)
}

func (s *ShardConfig) DelHostGroupCluster(clusterName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	hostgroupclusters := s.currentConf.HostGroupClusters
	for i, h := range hostgroupclusters {
		if h.Name == clusterName {
			s.currentConf.HostGroupClusters = append(hostgroupclusters[:i], hostgroupclusters[i+1:]...)
			return
		}
	}
}

//func (s *ShardConfig)UpdateRwSplit(name string)
func PrettyPrint(v interface{}) (err error, s string) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		glog.Glog.Infof(string(b))
	}
	return
}

func deepCopy(old *config.Config, cur *config.Config) error {
	if old == nil {
		return errors.Errorf("old cannot be nil")
	}
	if cur == nil {
		return errors.Errorf("cur cannot be nil")
	}
	bytes, err := yaml.Marshal(cur)
	if err != nil {
		return errors.Errorf("Unable to marshal cur: %s", err)
	}
	err = yaml.Unmarshal(bytes, old)
	if err != nil {
		return errors.Errorf("Unable to unmarshal into old: %s", err)
	}
	return nil
}
