package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
)

func (e *ProxyEngine) ParserConfig(cfg *config.Config) error {
	if len(cfg.Basic.User) == 0 {
		return errcode.BuildError(errcode.ObjectEmptyErr, "cfg.Basic.User")
	}

	if len(cfg.Basic.Password) == 0 {
		return errcode.BuildError(errcode.ObjectEmptyErr, "cfg.Basic.Password")
	}

	if cfg.Basic.SlowLogTime == 0 {
		return errcode.BuildError(errcode.ObjectEmptyErr, "cfg.Basic.SlowLogTime")
	}

	e.user = cfg.Basic.User
	e.password = cfg.Basic.Password
	e.slowLogTime = cfg.Basic.SlowLogTime

	glog.Glog.Infof("Parsing hostGroups")
	if err := e.ParseHostGroupNodes(cfg, true); err != nil {
		return err
	}
	glog.Glog.Infof("Parsed hostGroups")

	if err := e.ParseHostGroupNodeClusters(cfg); err != nil {
		return err
	}

	// 构建整个schema和hostGroup和table的聚合逻辑
	glog.Glog.Infof("Parsing schemas and tables")
	if err := e.ParseSchemas(cfg); err != nil {
		return err
	}
	glog.Glog.Infof("Parsed schemas and tables")

	return nil
}

func (e *ProxyEngine) ReconfigBasic(basic *config.Basic) {
	e.slowLogTime = basic.SlowLogTime
	e.user = basic.User
	e.password = basic.Password
}

func (e *ProxyEngine) parseHostGroupNode(cfg config.HostGroup, initPool bool) (*backend.HostGroupNode, error) {
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

func (e *ProxyEngine) ParseHostGroupNodes(cfg *config.Config, initConn bool) error {
	for _, g := range cfg.HostGroups {
		if len(g.Name) == 0 {
			return errcode.BuildError(errcode.ObjectEmptyErr, "the name of hostGroup")
		}

		if e.GetHostGroupNode(g.Name) != nil {
			return fmt.Errorf("duplicate hostGroup [%s]", g.Name)
		}

		n, err := e.parseHostGroupNode(g, initConn)
		if err != nil {
			return err
		}

		e.HostGroupNodes.Store(g.Name, n)

		glog.Glog.Infof("HostGroup:[%s] WriteIPs[%s] has been parsed.", g.Name, g.Write)
	}

	return nil
}

func (e *ProxyEngine) ParseHostGroupNodeClusters(cfg *config.Config) error {
	for _, gc := range cfg.HostGroupClusters {
		if len(gc.Name) == 0 {
			return errcode.BuildError(errcode.ObjectEmptyErr, "the name of hostGroupCluster")
		}
		p := gc
		e.HostGroupClusters.Store(gc.Name, &p)
		glog.Glog.Infof("HostGroupCluster:[%s] has been parsed.", gc.Name)
	}

	return nil
}

func (e *ProxyEngine) ParseSchemas(cfg *config.Config) error {
	schemaCfgs := cfg.SchemaConfigs
	rt, err := NewRouter()
	e.router = rt
	for _, v := range schemaCfgs {
		sc := v
		if len(sc.Name) == 0 {
			return fmt.Errorf("empty schema is not allowed")
		}
		if e.GetSchemaNode(sc.Name) != nil {
			return fmt.Errorf("duplicated schema [%s]", sc.Name)
		}

		if len(sc.HostGroupCluster) == 0 {
			return fmt.Errorf("schema:[%s], empty hostGroupCluster is not allowed", sc.Name)
		}

		hostGroupClusterConf, ok := e.HostGroupClusters.Load(sc.HostGroupCluster)
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
			err = e.parseCustodySchema(&sc, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			glog.Glog.Infof("Schema:[%s] [custody type] has been parsed.", sc.Name)
		} else {
			schemaRule, err := rt.LoadShardingSchemaRule(&sc, hostGroupClusterConf.(*config.HostGroupCluster))
			if err != nil {
				return err
			}
			err = e.parseSchema(&sc, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
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

func (e *ProxyEngine) parseCustodySchema(schemaCfg *config.SchemaConfig, schemaRouteRule *SchemaRouteRule, hostGroupClusterConf *config.HostGroupCluster) error {
	if e.GetSchemaNode(schemaCfg.Name) != nil {
		return fmt.Errorf("duplicated schema: [%s]", schemaCfg.Name)
	}

	// Custody
	if len(hostGroupClusterConf.NonshardingHostGroup) != 0 && e.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup) == nil {
		return fmt.Errorf("hostGroup [%s] config does not exist", hostGroupClusterConf.NonshardingHostGroup)
	}

	sn := new(backend.SchemaNode)
	sn.Name = schemaCfg.Name
	sn.NonshardingHostGroupNode = e.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup)
	sn.RwSplit = schemaCfg.RwSplit
	sn.MultiRoutePermitted = schemaCfg.MultiRoute
	sn.InitConnMultiRoutePermitted = schemaCfg.InitConnMultiRoute
	sn.Conf = schemaCfg
	sn.HostGroupCluster = hostGroupClusterConf.Name

	e.schemaNodes.Store(schemaCfg.Name, sn)
	return nil
}

func (e *ProxyEngine) parseSchema(schemaCfg *config.SchemaConfig, schemaRouteRule *SchemaRouteRule, hostGroupClusterConf *config.HostGroupCluster) error {
	if e.GetSchemaNode(schemaCfg.Name) != nil {
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
		if e.GetHostGroupNode(groupName) == nil {
			return fmt.Errorf("host group [%s] config doesn't exist", groupName)
		}

		sn.ShardingHostGroupNodes[groupName] = e.GetHostGroupNode(groupName)
	}

	// nonsharding
	if len(hostGroupClusterConf.NonshardingHostGroup) != 0 && e.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup) == nil {
		return fmt.Errorf("hostGroup [%s] config does not exist", hostGroupClusterConf.NonshardingHostGroup)
	}
	sn.NonshardingHostGroupNode = e.GetHostGroupNode(hostGroupClusterConf.NonshardingHostGroup)

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

	e.schemaNodes.Store(schemaCfg.Name, sn)

	return nil
}
