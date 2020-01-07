package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"github.com/pkg/errors"
	"sync"
)

//schema配置
type ProxyEngineConfig struct {
	HostGroupClusters sync.Map
	HostGroupNodes    sync.Map
	SchemaNodes       sync.Map
	Router            *Router
	Version           int64 //用于表识当前使用的schema config的版本
}

//期望的配置结构
type ProxyEngineConfigBak struct {
	lock              sync.Mutex
	HostGroupClusters []config.HostGroupCluster
	hostGroupNodes    map[string]*backend.HostGroupNode
	SchemaNodes       map[string]*backend.SchemaNode
	Router            *Router
	SlowLogTime       int
	Version           string
}

type schemhgclusters struct {
	SchemaConfigs     []config.SchemaConfig
	HostGroupClusters []config.HostGroupCluster
}

func (p *ProxyEngineConfig) getHostGroupCluster(name string) (*config.HostGroupCluster, error) {
	value, ok := p.HostGroupClusters.Load(name)
	if ok {
		return value.(*config.HostGroupCluster), nil
	}
	return nil, errors.Errorf("hostgroupcluster %s not exist", name)
}

//解析hostgroup配置数据到hostgroupNode中
func (p *ProxyEngineConfig) ParseHostGroupNode(cfg config.HostGroup, initPool bool) (*backend.HostGroupNode, error) {
	var err error
	n := new(backend.HostGroupNode)
	n.Cfg = cfg
	n.InitActivedWriteIndex(cfg.ActiveWrite)
	n.EnableWriting()

	if cfg.MaxConn < 1024 {
		return nil, fmt.Errorf("the MaxConn:[%d] of hostGroup:[%s] should bigger than 1024", cfg.MaxConn, cfg.Name)
	}

	if cfg.InitConn > cfg.MaxConn {
		return nil, fmt.Errorf("the MaxConn:[%d] of hostGroup:[%s] should bigger than InitConn:[%d]",
			cfg.MaxConn, cfg.Name, cfg.InitConn)
	}

	if len(cfg.User) == 0 {
		return nil, errcode.BuildError(errcode.ObjectEmptyErr, "cfg.User")
	}

	// TODO nanxing 异步化创建
	err = n.InitWrite(cfg.Write, initPool)
	if err != nil {
		return nil, err
	}

	err = n.InitRead(cfg.Read, initPool)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (p *ProxyEngineConfig) ParseHostGroupClusters(cfg *config.Config) {
	for _, gc := range cfg.HostGroupClusters {
		if len(gc.Name) == 0 {
			return
		}
		h := gc
		p.HostGroupClusters.Store(gc.Name, &h)
		glog.Glog.Infof("HostGroupCluster:[%s] has been parsed.", gc.Name)
	}

}

// ParSeConfig 负责字节数据解析到ProxyEngineConfig中
func (p *ProxyEngineConfig) ParseConfig(cfg *config.Config) error {
	//解析HostGroupClusters
	p.ParseHostGroupClusters(cfg)
	//解析hostgruop配置
	glog.Glog.Infof("Parsing hostGroups")
	for _, g := range cfg.HostGroups {
		if len(g.Name) == 0 {
			return errcode.BuildError(errcode.ObjectEmptyErr, "the name of hostGroup")
		}

		if p.GetHostGroupNode(g.Name) != nil {
			return fmt.Errorf("duplicate hostGroup [%s]", g.Name)
		}

		n, err := p.ParseHostGroupNode(g, true)
		if err != nil {
			return err
		}

		p.HostGroupNodes.Store(g.Name, n)

		glog.Glog.Infof("HostGroup:[%s] WriteIPs[%s] has been parsed.", g.Name, g.Write)

	}
	glog.Glog.Infof("Parsed hostGroups")

	// 构建整个schema和hostGroup和table的聚合逻辑
	glog.Glog.Infof("Parsing schemas and tables")
	if err := p.parseSchemas(cfg); err != nil {
		return err
	}
	glog.Glog.Infof("Parsed schemas and tables")
	return nil
}

func (p *ProxyEngineConfig) GetHostGroupNode(name string) *backend.HostGroupNode {
	result, ok := p.HostGroupNodes.Load(name)
	if ok {
		return result.(*backend.HostGroupNode)
	} else {
		return nil
	}
}

func (p *ProxyEngineConfig) ParseCustodySchema(schemaCfg *config.SchemaConfig, schemaRouteRule *SchemaRouteRule) (*backend.SchemaNode, error) {

	sn := new(backend.SchemaNode)
	sn.Name = schemaCfg.Name
	sn.NonshardingHostGroupNode = p.GetHostGroupNode(schemaCfg.HostGroupCluster)
	sn.RwSplit = schemaCfg.RwSplit
	sn.MultiRoutePermitted = schemaCfg.MultiRoute
	sn.InitConnMultiRoutePermitted = schemaCfg.InitConnMultiRoute
	return sn, nil
}

func (p *ProxyEngineConfig) parseSchemas(cfg *config.Config) error {
	schemaCfgs := cfg.SchemaConfigs
	rt, err := NewRouter()
	p.Router = rt
	for _, v := range schemaCfgs {
		sc := v
		if len(sc.Name) == 0 {
			return fmt.Errorf("empty schema is not allowed")
		}
		if p.GetSchemaNode(sc.Name) != nil {
			return fmt.Errorf("duplicated schema [%s]", sc.Name)
		}

		if len(sc.HostGroupCluster) == 0 {
			return fmt.Errorf("schema:[%s], empty hostGroupCluster is not allowed", sc.Name)
		}

		hostGroupClusterConf, ok := p.HostGroupClusters.Load(sc.HostGroupCluster)
		if !ok {
			return fmt.Errorf("hostGroupCluster [%s] config does not exist", sc.HostGroupCluster)
		}

		// 是否仅仅是托管
		// 防止分支逻辑太杂乱，就用不同的函数来分离逻辑
		if sc.Custody {
			schemaRule, err := rt.LoadCustodySchemaRule(&sc, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			err = p.parseCustodySchema(&sc, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			glog.Glog.Infof("Schema:[%s] [custody type] has been parsed.", sc.Name)
		} else {
			schemaRule, err := rt.LoadShardingSchemaRule(&sc, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			err = p.parseSchema(&sc, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			glog.Glog.Infof("Schema:[%s] [sharding type %v * %v] has been parsed.", sc.Name, sc.SchemaSharding, sc.TableSharding)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *ProxyEngineConfig) GetSchemaNode(name string) *backend.SchemaNode {
	result, ok := p.SchemaNodes.Load(name)
	if ok {
		return result.(*backend.SchemaNode)
	} else {
		return nil
	}
}

func (p *ProxyEngineConfig) parseCustodySchema(schemaCfg *config.SchemaConfig, schemaRouteRule *SchemaRouteRule, hostGroupClusterConf *config.HostGroupCluster) error {
	if p.GetSchemaNode(schemaCfg.Name) != nil {
		return fmt.Errorf("duplicated schema: [%s]", schemaCfg.Name)
	}

	// Custody
	if len(hostGroupClusterConf.NonshardingHostGroup) != 0 && p.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup) == nil {
		return fmt.Errorf("hostGroup [%s] config does not exist", hostGroupClusterConf.NonshardingHostGroup)
	}

	sn := new(backend.SchemaNode)
	sn.Name = schemaCfg.Name
	sn.NonshardingHostGroupNode = p.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup)
	sn.RwSplit = schemaCfg.RwSplit
	sn.MultiRoutePermitted = schemaCfg.MultiRoute
	sn.InitConnMultiRoutePermitted = schemaCfg.InitConnMultiRoute
	sn.Conf = schemaCfg
	sn.HostGroupCluster = hostGroupClusterConf.Name

	p.SchemaNodes.Store(schemaCfg.Name, sn)
	return nil
}

func (p *ProxyEngineConfig) parseSchema(schemaCfg *config.SchemaConfig, schemaRouteRule *SchemaRouteRule, hostGroupClusterConf *config.HostGroupCluster) error {
	if p.GetSchemaNode(schemaCfg.Name) != nil {
		return fmt.Errorf("duplicated schema: [%s]", schemaCfg.Name)
	}

	sn := new(backend.SchemaNode)
	sn.Name = schemaCfg.Name
	sn.ShardingHostGroupNodes = make(map[string]*backend.HostGroupNode)
	sn.RwSplit = schemaCfg.RwSplit
	sn.MultiRoutePermitted = schemaCfg.MultiRoute
	sn.InitConnMultiRoutePermitted = schemaCfg.InitConnMultiRoute
	sn.Conf = schemaCfg
	sn.HostGroupCluster = hostGroupClusterConf.Name

	for _, groupName := range hostGroupClusterConf.ShardingHostGroups {
		if p.GetHostGroupNode(groupName) == nil {
			return fmt.Errorf("host group [%s] config doesn't exist", groupName)
		}

		sn.ShardingHostGroupNodes[groupName] = p.GetHostGroupNode(groupName)
	}

	// nonsharding
	if len(hostGroupClusterConf.NonshardingHostGroup) != 0 && p.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup) == nil {
		return fmt.Errorf("hostGroup [%s] config does not exist", hostGroupClusterConf.NonshardingHostGroup)
	}
	sn.NonshardingHostGroupNode = p.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup)

	for _, tc := range schemaCfg.TableConfigs {
		if len(tc.Name) == 0 {
			return fmt.Errorf("empty table is not allowed")
		}
		if sn.GetTableNode(tc.Name) != nil {
			return fmt.Errorf("depulicated table node: [%s]", tc.Name)
		}

		if err := schemaRouteRule.AddTableRouteRule(tc); err != nil {
			return err
		}

		if sn.GetTableNode(tc.Name) == nil {
			sn.TableNodes.Store(tc.Name, &backend.TableNode{
				Name:        tc.Name,
				ShardingKey: tc.ShardingKey,
			})
		}

		glog.Glog.Infof("Table:[%s] has been parsed.", tc.Name)
	}

	p.SchemaNodes.Store(schemaCfg.Name, sn)

	return nil
}
