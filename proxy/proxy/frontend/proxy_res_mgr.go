package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
)

func (e *ProxyEngine) AddHostGroup(hostGroupCfg config.HostGroup) (string, error) {
	if e.GetHostGroupNode(hostGroupCfg.Name) != nil {
		return fmt.Sprintf("HostGroup[%s] has been already added", hostGroupCfg.Name), nil
	}

	glog.Glog.Infof("parsing new hostGroup[%s]", hostGroupCfg.Name)
	n, err := e.parseHostGroupNode(hostGroupCfg, true)
	if err != nil {
		return "ok", err
	}

	e.HostGroupNodes.Store(hostGroupCfg.Name, n)

	glog.Glog.Infof("hostGroup[%s] has been added.", hostGroupCfg.Name)
	return fmt.Sprintf("AddHostGroup [%s] success", hostGroupCfg.Name), nil
}

func (e *ProxyEngine) DeleteHostGroup(hostGroupName string) (string, error) {
	if e.GetHostGroupNode(hostGroupName) == nil {
		return fmt.Sprintf("HostGroup[%s] has been already deleted", hostGroupName), nil
	}

	hostGroup := e.GetHostGroupNode(hostGroupName)
	if used, schema := hostGroupBeingUsed(e, hostGroup); used {
		return "", fmt.Errorf("unable to delete hostGroup[%s] since it is being used by hostGroupCluster[%s]", hostGroupName, schema)
	}

	for _, wPool := range hostGroup.Write {
		wPool.CloseResource()
		wPool.CloseActiveBackendConns()
	}

	for _, rPool := range hostGroup.Read {
		rPool.CloseResource()
		rPool.CloseActiveBackendConns()
	}

	hostGroup = nil

	e.HostGroupNodes.Delete(hostGroupName)

	glog.Glog.Infof("hostGroup[%s] has been deleted.", hostGroupName)
	return fmt.Sprintf("DeleteHostGroup [%s] success", hostGroupName), nil
}

func (e *ProxyEngine) AddHostGroupCluster(clusterCfg config.HostGroupCluster) (string, error) {
	if e.GetHostGroupCluster(clusterCfg.Name) != nil {
		return fmt.Sprintf("HostGroupCluster[%s] has been already added", clusterCfg.Name), nil
	}

	if len(clusterCfg.NonshardingHostGroup) != 0 && e.GetHostGroupNode(clusterCfg.NonshardingHostGroup) == nil {
		return "", fmt.Errorf("HostGroupCluster[%s]: NonshardingHostGroup[%s] does not exist", clusterCfg.Name, clusterCfg.NonshardingHostGroup)
	}

	existCluster := make(map[string]bool)
	for _, hcg := range clusterCfg.ShardingHostGroups {
		if len(hcg) == 0 || e.GetHostGroupNode(hcg) == nil {
			return "", fmt.Errorf("HostGroupCluster[%s]: ShardingHostGroup[%s] does not exist", clusterCfg.Name, hcg)
		}

		if existCluster[hcg] {
			return "", fmt.Errorf("HostGroupCluster[%s]: ShardingHostGroup[%s] duplicates", clusterCfg.Name, hcg)
		}

		existCluster[hcg] = true
	}

	e.HostGroupClusters.Store(clusterCfg.Name, &clusterCfg)

	glog.Glog.Infof("HostGroupCluster[%s] has been added.", clusterCfg.Name)
	return fmt.Sprintf("AddHostGroupCluster [%s] success", clusterCfg.Name), nil
}

func (e *ProxyEngine) UpdateHostGroupCluster(clusterCfg config.HostGroupCluster) (string, error) {
	if e.GetHostGroupCluster(clusterCfg.Name) == nil {
		return fmt.Sprintf("HostGroupCluster[%s] does not exist", clusterCfg.Name), nil
	}

	var err error

	e.schemaNodes.Range(func(key, value interface{}) bool {
		if v, ok := value.(*backend.SchemaNode); ok {
			if v.HostGroupCluster == clusterCfg.Name {
				// 卸载原来的规则
				err = e.router.UnloadSchemaRule(v.Name)
				if err != nil {
					return false
				}
				// 重新安装规则
				_, err = e.router.LoadSchemaRule(v.Conf, &clusterCfg)
				if err != nil {
					return false
				}
			}
		}

		return true
	})

	if err != nil {
		return "", err
	}

	e.HostGroupClusters.Delete(clusterCfg.Name)
	e.HostGroupClusters.Store(clusterCfg.Name, &clusterCfg)

	glog.Glog.Infof("HostGroupCluster[%s] has been updated.", clusterCfg.Name)
	return fmt.Sprintf("UpdateHostGroupCluster [%s] success", clusterCfg.Name), err
}

func (e *ProxyEngine) DeleteHostGroupCluster(clusterName string) (string, error) {
	if e.GetHostGroupCluster(clusterName) == nil {
		return fmt.Sprintf("HostGroupClusterName[%s] has been already deleted", clusterName), nil
	}
	// if it's being used by schema
	if used, schema := hostGroupClusterBeingUsed(e, clusterName); used {
		return "", fmt.Errorf("unable to delete hostGroupCluster[%s] since it is being used by schema[%s]", clusterName, schema)
	}

	e.HostGroupClusters.Delete(clusterName)

	glog.Glog.Infof("hostGroupCluster[%s] has been deleted.", clusterName)
	return fmt.Sprintf("DeleteHostGroupCluster [%s] success", clusterName), nil
}

func (e *ProxyEngine) AddSchema(schemaCfg config.SchemaConfig) (string, error) {
	if len(schemaCfg.Name) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "schemaCfg.Name")
	}

	if e.GetSchemaNode(schemaCfg.Name) != nil {
		return fmt.Sprintf("Schema[%s] has been already added", schemaCfg.Name), nil
	}

	hostGroupClusterConf, ok := e.HostGroupClusters.Load(schemaCfg.HostGroupCluster)
	if !ok {
		return "", fmt.Errorf("hostGroupCluster [%s] config does not exist", schemaCfg.HostGroupCluster)
	}

	if schemaCfg.Custody {
		schemaRule, err := e.router.LoadCustodySchemaRule(&schemaCfg, hostGroupClusterConf.(*config.HostGroupCluster))
		if err != nil {
			return "", err
		}
		err = e.parseCustodySchema(&schemaCfg, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
		if err != nil {
			return "", err
		}
		glog.Glog.Infof("Schema[%s] [custody type] has been parsed.", schemaCfg.Name)
	} else {
		schemaRule, err := e.router.LoadShardingSchemaRule(&schemaCfg, hostGroupClusterConf.(*config.HostGroupCluster))
		if err != nil {
			return "", err
		}
		err = e.parseSchema(&schemaCfg, schemaRule, hostGroupClusterConf.(*config.HostGroupCluster))
		if err != nil {
			return "", err
		}
		glog.Glog.Infof("Schema[%s] [sharding type %v * %v] has been parsed.", schemaCfg.Name, schemaCfg.SchemaSharding, schemaCfg.TableSharding)
	}

	glog.Glog.Infof("schema[%s] has been added.", schemaCfg.Name)
	return fmt.Sprintf("AddSchema [%s] success", schemaCfg.Name), nil
}

func (e *ProxyEngine) DeleteSchema(schemaName string) (string, error) {
	if len(schemaName) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "schemaName")
	}

	// delete schema route
	e.router.DeleteSchemaRule(schemaName)

	if e.GetSchemaNode(schemaName) == nil {
		return fmt.Sprintf("Schema[%s] has been already deleted", schemaName), nil
	}
	// delete schema
	e.schemaNodes.Delete(schemaName)

	glog.Glog.Infof("schema[%s] has been deleted.", schemaName)
	return fmt.Sprintf("DeleteSchema [%s] success", schemaName), nil
}

func (e *ProxyEngine) AddTable(tableCfg config.TableConfig, schemaName string) (string, error) {
	if len(schemaName) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "schemaName")
	}
	sn := e.GetSchemaNode(schemaName)

	if sn == nil {
		return "", fmt.Errorf("schema node [%s] doesn't exist", schemaName)
	}

	if len(tableCfg.Name) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "tableCfg.Name")
	}

	if len(tableCfg.ShardingKey) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "tableCfg.ShardingKey")
	}

	if sn.GetTableNode(tableCfg.Name) != nil {
		return fmt.Sprintf("Table[%s] has been already added", tableCfg.Name), nil
	}

	srr := e.router.GetSchemaRule(schemaName)

	if srr == nil {
		return "", fmt.Errorf("schema route [%s] doesn't exist", schemaName)
	}

	if err := srr.AddTableRouteRule(tableCfg); err != nil {
		return "", err
	}

	if sn.GetTableNode(tableCfg.Name) == nil {
		sn.TableNodes.Store(tableCfg.Name, &backend.TableNode{
			Name:        tableCfg.Name,
			ShardingKey: tableCfg.ShardingKey,
		})
	}

	glog.Glog.Infof("table[%s] has been added.", tableCfg.Name)
	return fmt.Sprintf("AddTable [%s] in schema [%s] success", tableCfg.Name, schemaName), nil
}

func (e *ProxyEngine) DeleteTable(tableName, schemaName string) (string, error) {
	if len(tableName) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "tableName")
	}
	if len(schemaName) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "schemaName")
	}

	sn := e.GetSchemaNode(schemaName)

	if sn == nil {
		return "", fmt.Errorf("schema node [%s] doesn't exist", schemaName)
	}

	if sn.GetTableNode(tableName) == nil {
		return fmt.Sprintf("Table[%s] has been already deleted", schemaName), nil
	}

	srr := e.router.GetSchemaRule(schemaName)

	// delete table route
	srr.TableRouteRules.Delete(tableName)
	// delete table
	sn.TableNodes.Delete(tableName)

	glog.Glog.Infof("table[%s] has been deleted.", tableName)
	return fmt.Sprintf("DeleteTable [%s] in [%s] success", tableName, schemaName), nil
}

func (e *ProxyEngine) UpdateSchemaRWSplit(schema string, rwSplit bool) (string, error) {
	if len(schema) == 0 {
		return "", errcode.BuildError(errcode.ObjectEmptyErr, "schema")
	}
	sn := e.GetSchemaNode(schema)
	if sn == nil {
		return "", fmt.Errorf("invalid schema [%s]", schema)
	}
	from := sn.RwSplit
	sn.RwSplit = rwSplit

	glog.Glog.Infof("schema[%s]'s RWSplit has been updated from: %v -> to: %v", schema, from, rwSplit)
	return fmt.Sprintf("UpdateSchemaRWSplit [%s] success", schema), nil
}

func (e *ProxyEngine) ApplyLatestConfig(newCfg config.Config) error {
	e.confApplyLock.Lock()
	defer e.confApplyLock.Unlock()
	glog.Glog.Infof("applying latest config")

	glog.Glog.Infof("close listener")
	e.CloseProxyEngine()

	var err error

	// 清理schema
	e.schemaNodes.Range(func(key, value interface{}) bool {
		if v, ok := value.(*backend.SchemaNode); ok {

			v.TableNodes.Range(func(key, value interface{}) bool {
				if k, ok := key.(string); ok {
					_, err = e.DeleteTable(k, v.Name)
					if err != nil {
						panic(err)
					}
				}
				return true
			})

			_, err = e.DeleteSchema(v.Name)
			value = nil
			if err != nil {
				panic(err)
			}
		}
		return true
	})

	// 清理hostGroupCluster
	e.HostGroupClusters.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			_, err = e.DeleteHostGroupCluster(k)
			value = nil
			if err != nil {
				panic(err)
			}
		}

		return true
	})

	// 清理hostGroup
	e.HostGroupNodes.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			_, err = e.DeleteHostGroup(k)
			value = nil
			if err != nil {
				panic(err)
			}
		}

		return true
	})

	err = e.ParserConfig(&newCfg)

	if err != nil {
		return err
	}

	go e.ServeInBlocking(func() {
	})

	glog.Glog.Infof("latest config has been applied.")

	return nil
}

//应用上一个版本的schema配置
func (e *ProxyEngine) ApplyBackConfig(cfg *config.Config) error {

	glog.Glog.Infof("rollback config")
	peconf := &ProxyEngineConfig{}
	if err := peconf.ParseConfig(cfg); err != nil {
		return err
	}
	//清理前后端链接
	e.cleanUpBackendConnections()
	e.HostGroupNodes = peconf.HostGroupNodes
	e.schemaNodes = peconf.SchemaNodes
	e.router = peconf.Router
	e.HostGroupClusters = peconf.HostGroupClusters
	glog.Glog.Infof("rollback schema config %v", peconf)
	glog.Glog.Infof("success rollback config.")
	return nil
}

func (e *ProxyEngine) GetHostGroupNode(name string) *backend.HostGroupNode {
	result, ok := e.HostGroupNodes.Load(name)
	if ok {
		return result.(*backend.HostGroupNode)
	} else {
		return nil
	}
}

func (e *ProxyEngine) GetHostGroupCluster(name string) *config.HostGroupCluster {
	result, ok := e.HostGroupClusters.Load(name)
	if ok {
		return result.(*config.HostGroupCluster)
	} else {
		return nil
	}
}

func (e *ProxyEngine) GetSchemaNode(name string) *backend.SchemaNode {
	result, ok := e.schemaNodes.Load(name)
	if ok {
		return result.(*backend.SchemaNode)
	} else {
		return nil
	}
}

func (e *ProxyEngine) GetRouter() *Router {
	return e.router
}

func (e *ProxyEngine) GetSlowLogTime() int {
	return e.slowLogTime
}

func (e *ProxyEngine) StopWritingAbility(hostGroup string) error {
	e.updateHostGroupWritable(hostGroup, false)
	glog.Glog.Infof("hostGroup[%s] stop writing ability.", hostGroup)
	return nil
}

func (e *ProxyEngine) RecoverWritingAbility(hostGroup string) error {
	e.updateHostGroupWritable(hostGroup, true)
	glog.Glog.Infof("hostGroup[%s] recovered writing ability.", hostGroup)
	return nil
}

func (e *ProxyEngine) updateHostGroupWritable(hostGroup string, writable bool) error {
	if len(hostGroup) == 0 {
		return errcode.BuildError(errcode.ObjectEmptyErr, "hostGroup")
	}
	n := e.GetHostGroupNode(hostGroup)
	if n == nil {
		return errcode.BuildError(errcode.HostGroupNotExist, hostGroup)
	}

	if writable {
		n.EnableWriting()
	} else {
		n.StopWriting()
	}
	return nil
}

func hostGroupBeingUsed(s *ProxyEngine, hostGroup *backend.HostGroupNode) (bool, string) {
	var used bool
	var hostGroupCluster string

	s.HostGroupClusters.Range(func(key, value interface{}) bool {
		if v, ok := value.(*config.HostGroupCluster); ok {
			if v.NonshardingHostGroup == hostGroup.Cfg.Name {
				used = true
				hostGroupCluster = v.Name
				return false
			}

			for _, n := range v.ShardingHostGroups {
				if n == hostGroup.Cfg.Name {
					used = true
					hostGroupCluster = v.Name
					return false
				}
			}
		}

		return true
	})

	return used, hostGroupCluster
}

func hostGroupClusterBeingUsed(s *ProxyEngine, hostGroupClusterName string) (bool, string) {
	var used bool
	var schemaName string

	s.schemaNodes.Range(func(key, value interface{}) bool {
		if v, ok := value.(*backend.SchemaNode); ok {
			if v.HostGroupCluster == hostGroupClusterName {
				used = true
				schemaName = v.Name
				return false
			}
		}

		return true
	})

	return used, schemaName
}

func (e *ProxyEngine) DeleteRead(node, readAddr string) error {
	n := e.GetHostGroupNode(node)
	if n == nil {
		return fmt.Errorf("invalid node [%s]", node)
	}
	return n.DeleteRead(readAddr)
}

func (e *ProxyEngine) AddRead(node, readAddr string) error {
	n := e.GetHostGroupNode(node)
	if n == nil {
		return fmt.Errorf("invalid node [%s]", node)
	}
	return n.DeleteRead(readAddr)
}



