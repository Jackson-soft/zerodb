package etcdtool

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"

	"git.2dfire.net/zerodb/common/config"

	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
)

//Puts 批量放入多值
func (s *Store) Puts(datas map[string]string) error {
	if len(datas) == 0 {
		return parameterIsNil
	}

	ops := make([]clientv3.Op, 0)
	for key, value := range datas {
		ops = append(ops, clientv3.OpPut(key, value))
	}

	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	resp, err := s.cli.Txn(ctx).Then(ops...).Commit()
	if err != nil {
		return err
	}
	cancel()
	var keys []string
	for k, _ := range datas {
		keys = append(keys, k)
	}
	//backup lastModify
	return s.BackupLastModify(resp.Header.Revision, keys)
}

// Gets 批量获取值
func (s *Store) Gets(keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, parameterIsNil
	}

	ops := make([]clientv3.Op, 0)
	for i := range keys {
		ops = append(ops, clientv3.OpGet(keys[i]))
	}

	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	resp, err := s.cli.Txn(ctx).Then(ops...).Commit()
	cancel()
	if err != nil {
		return nil, err
	}

	retData := make(map[string]string)

	for i := range resp.Responses {
		retData[string(resp.Responses[i].GetResponseRange().Kvs[0].Key)] = string(resp.Responses[i].GetResponseRange().Kvs[0].Value)
	}

	return retData, nil
}

//Deletes 批量删除
func (s *Store) Deletes(keys []string) error {
	if len(keys) == 0 {
		return parameterIsNil
	}

	ops := make([]clientv3.Op, 0)
	for i := range keys {
		ops = append(ops, clientv3.OpDelete(keys[i]))
	}

	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	_, err := s.cli.Txn(ctx).Then(ops...).Commit()
	cancel()
	return err
}

//Delete 删除key
func (s *Store) Delete(key string) error {
	if len(key) == 0 {
		return parameterIsNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	_, err := s.cli.Delete(ctx, key)
	cancel()

	return err
}

//DeleteKeys 批量删除带有前缀的key
func (s *Store) DeleteKeys(prefix string) error {
	if len(prefix) == 0 {
		return parameterIsNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	_, err := s.cli.Delete(ctx, prefix, clientv3.WithPrefix())
	cancel()

	return err
}

func (s *Store) AddHostGroup(clusterName string, groups []config.HostGroup) error {
	if len(clusterName) == 0 || groups == nil {
		return parameterIsNil
	}
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	hosts := make([]config.HostGroup, 0)

	if err = utils.LoadYaml(data, &hosts); err != nil {
		return err
	}

	for _, group := range groups {
		//添加前做简单检查
		if len(group.Name) == 0 || group.InitConn > group.MaxConn {
			continue
		}

		//是否有重复名称
		isRepetition := false
		for _, host := range hosts {
			if host.Name == group.Name {
				isRepetition = true
				continue
			}
		}
		if !isRepetition {
			hosts = append(hosts, group)
		}
	}

	data, err = utils.UnLoadYaml(&hosts)
	if err != nil {
		return err
	}
	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("add hostGroup successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

func (s *Store) DeleteGroup(clusterName string, names []string) error {
	if len(clusterName) == 0 || len(names) == 0 {
		return parameterIsNil
	}
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	hosts := make([]config.HostGroup, 0)

	if err = utils.LoadYaml(data, &hosts); err != nil {
		return err
	}

	for i := 0; i < len(names); i++ {
		for j := 0; j < len(hosts); j++ {
			if hosts[j].Name == names[i] {
				hosts = append(hosts[:j], hosts[j+1:]...)
				break
			}
		}
	}

	data, err = utils.UnLoadYaml(&hosts)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("delete hostGroup successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

func (s *Store) UpdateStopSve(clusterName string, stopSvr config.StopService) error {
	if len(clusterName) == 0 {
		return parameterIsNil
	}

	data, err := utils.UnLoadYaml(&stopSvr)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, StopService)

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("update stopService successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

// AddClusters add host group cluster
func (s *Store) AddClusters(clusterName string, clusters []config.HostGroupCluster) error {
	if len(clusterName) == 0 || len(clusters) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroupClusters)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	myResults := make([]config.HostGroupCluster, 0)
	if err = utils.LoadYaml(data, &myResults); err != nil {
		return err
	}

	for j := range clusters {
		isRepetition := false
		for i := range myResults {
			if clusters[j].Name == myResults[i].Name {
				isRepetition = true
				continue
			}
		}
		if !isRepetition {
			myResults = append(myResults, clusters[j])
		}
	}

	data, err = utils.UnLoadYaml(&myResults)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("add cluster successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

//DeleteClusters delete cluster
func (s *Store) DeleteClusters(clusterName string, names []string) error {
	if len(clusterName) == 0 || len(names) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroupClusters)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	myResults := make([]config.HostGroupCluster, 0)
	if err = utils.LoadYaml(data, &myResults); err != nil {
		return err
	}

	for i := 0; i < len(names); i++ {
		for j := 0; j < len(myResults); j++ {
			if myResults[j].Name == names[i] {
				myResults = append(myResults[:j], myResults[j+1:]...)
				break
			}
		}
	}

	data, err = utils.UnLoadYaml(&myResults)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)
	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("delete cluster successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

//AddSchema add schema
func (s *Store) AddSchema(clusterName string, schema []config.SchemaConfig) error {
	if len(clusterName) == 0 || len(schema) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, SchemaConfigs)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	schemas := make([]config.SchemaConfig, 0)

	if err = utils.LoadYaml(data, &schemas); err != nil {
		return err
	}

	for i := range schema {
		isRepetition := false
		for j := range schemas {
			if schema[i].Name == schemas[j].Name {
				isRepetition = true
				break
			}
		}
		if !isRepetition {
			schemas = append(schemas, schema[i])
		}
	}

	data, err = utils.UnLoadYaml(&schemas)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("add schema successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

func (s *Store) DeleteSchema(clusterName string, names []string) error {
	if len(clusterName) == 0 || len(names) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, SchemaConfigs)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	schemas := make([]config.SchemaConfig, 0)
	if err = utils.LoadYaml(data, &schemas); err != nil {
		return err
	}

	for i := 0; i < len(names); i++ {
		for j := 0; j < len(schemas); j++ {
			if schemas[j].Name == names[i] {
				schemas = append(schemas[:j], schemas[j+1:]...)
				break
			}
		}
	}

	data, err = utils.UnLoadYaml(&schemas)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("delete schema successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

func (s *Store) AddTable(clusterName, schema string, tables []config.TableConfig) error {
	if len(clusterName) == 0 || len(schema) == 0 || len(tables) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, SchemaConfigs)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	schemas := make([]config.SchemaConfig, 0)

	if err = utils.LoadYaml(data, &schemas); err != nil {
		return err
	}

	for i := range schemas {
		if schemas[i].Name == schema {
			for j := range tables {
				isRepetition := false
				for n := range schemas[i].TableConfigs {
					if schemas[i].TableConfigs[n].Name == tables[j].Name {
						isRepetition = true
						break
					}
				}
				if !isRepetition {
					schemas[i].TableConfigs = append(schemas[i].TableConfigs, tables[j])
				}
			}
		}
	}

	data, err = utils.UnLoadYaml(&schemas)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("add table successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

func (s *Store) DeleteTable(clusterName, schema string, tableName string) error {
	if len(clusterName) == 0 || len(schema) == 0 || len(tableName) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, SchemaConfigs)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	schemas := make([]config.SchemaConfig, 0)

	if err = utils.LoadYaml(data, &schemas); err != nil {
		return err
	}

	for i := 0; i < len(schemas); i++ {
		if schemas[i].Name == schema {
			for j := 0; j < len(schemas[i].TableConfigs); j++ {
				if schemas[i].TableConfigs[j].Name == tableName {
					schemas[i].TableConfigs = append(schemas[i].TableConfigs[:j], schemas[i].TableConfigs[j+1:]...)
					break
				}
			}
		}
	}

	data, err = utils.UnLoadYaml(&schemas)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("delete table successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

//UpdateShardRW 更新读写分离配置
func (s *Store) UpdateShardRW(clusterName string, rwMap map[string]bool) error {
	if len(clusterName) == 0 || len(rwMap) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, SchemaConfigs)
	data, err := s.GetData(key)
	if err != nil {
		return err
	}

	schemas := make([]config.SchemaConfig, 0)

	if err = utils.LoadYaml(data, &schemas); err != nil {
		return err
	}

	for name, rw := range rwMap {
		for i := range schemas {
			if schemas[i].Name == name {
				schemas[i].RwSplit = rw
			}
		}
	}

	data, err = utils.UnLoadYaml(&schemas)
	if err != nil {
		return err
	}

	sVersion, err := s.GetNextVersion(clusterName, false)
	if err != nil {
		return err
	}
	sKey := fmt.Sprintf(ConfigVersion, clusterName)

	ops := map[string]string{
		key:  string(data),
		sKey: sVersion,
	}

	glog.GLog.Infof("update read/write split successfully, node:[%s], value:[%s]", key, string(data))

	return s.Puts(ops)
}

/* GetNextVersion  获取配置文件的版本号
 * 返回+1之后的版本号
 */
func (s *Store) GetNextVersion(clusterName string, bigVersion bool) (string, error) {
	if len(clusterName) == 0 {
		return "", parameterIsNil
	}

	key := fmt.Sprintf(ConfigVersion, clusterName)
	data, err := s.GetData(key)
	if err != nil {
		return "", err
	}
	var sVersion string
	if len(data) == 0 {
		sVersion = "1.0.0"
	} else {
		myWant := strings.Split(string(data), ".")
		if len(myWant) != 3 {
			return "", errors.New("version layout error.")
		}

		first, err := strconv.Atoi(myWant[0])
		if err != nil {
			return "", err
		}

		// 大版本直接写入
		if bigVersion {
			first++
			sVersion = fmt.Sprintf("%d.0.0", first)
		} else {
			second, err := strconv.Atoi(myWant[1])
			if err != nil {
				return "", err
			}
			third, err := strconv.Atoi(myWant[2])
			if err != nil {
				return "", err
			}
			if third >= 9 {
				third = 0
				if second >= 9 {
					first++
					second = 0
				} else {
					second++
				}
			} else {
				third++
			}
			sVersion = fmt.Sprintf("%d.%d.%d", first, second, third)
		}
	}
	//return s.PutData(key, sVersion)
	return sVersion, nil
}

// 写入配置文件版本号
func (s *Store) SetVersion(clusterName string, bigVersion bool) error {
	sVersion, err := s.GetNextVersion(clusterName, bigVersion)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(ConfigVersion, clusterName)
	return s.PutData(key, sVersion)
}

//写入带版本号的数据
func (s *Store) PutDataWithVersion(clusterName, key, data, version string) error {
	sKey := fmt.Sprintf(ConfigVersion, clusterName)
	ops := map[string]string{
		key:  data,
		sKey: version,
	}

	return s.Puts(ops)
}
