package etcdtool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"git.2dfire.net/zerodb/keeper/pkg/models"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/utils"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

/*
 * 这一块要说明一个问题是从etcd中获取数据，如果数据为空统一返回 nil+nil
 */

//Store 存储相关的封装
type Store struct {
	cli       *clientv3.Client
	endpoints []string
}

//最后一次对etcd修改的Revision
type lastModify struct {
	Revision int64
	Keys     []string
}

var (
	parameterIsNil = errors.New("parameter is nil.")
)

const (
	etcdTimeOut = 3 * time.Second
)

//NewStore new a store client
func NewStore(endpoints []string) (*Store, error) {
	if len(endpoints) == 0 {
		return nil, parameterIsNil
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: etcdTimeOut,
	})
	if err != nil {
		return nil, err
	}

	s := new(Store)
	s.cli = cli
	s.endpoints = endpoints
	return s, nil
}

//GetClient 获取客户端
func (s *Store) GetClient() *clientv3.Client {
	return s.cli
}

//Close 关闭
func (s *Store) Close() error {
	return s.cli.Close()
}

//PutData 把key与value存入存储端
func (s *Store) PutData(key, data string) error {
	if len(key) == 0 || len(data) == 0 {
		return parameterIsNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Put(ctx, key, data)
	cancel()
	if err != nil {
		return err
	}
	// backup lastModify
	return s.BackupLastModify(resp.Header.Revision, []string{key})
}

//GetData 从存储端获取指定key的数据
func (s *Store) GetData(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, parameterIsNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	return resp.Kvs[0].Value, nil
}

//GetDatas return slice
func (s *Store) GetDatas(prefix string) (map[string][]byte, error) {
	if len(prefix) == 0 {
		return nil, parameterIsNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	reData := make(map[string][]byte)

	for _, value := range resp.Kvs {
		reData[string(value.Key)] = value.Value
	}

	return reData, nil
}

//GetShardConf 获取分库分表的指定的快照的完整配置
// snapshot 如果为空则获取当前在用配置
func (s *Store) GetShardConf(clusterName, snapshot string) ([]byte, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}

	if len(snapshot) == 0 {
		snapshot = DefaultConfig
	}

	key := fmt.Sprintf(ShardPrefix, clusterName, snapshot)

	mData, err := s.GetDatas(key)
	if err != nil {
		return nil, err
	}

	if len(mData) == 0 {
		return nil, nil
	}

	reData := config.Config{}
	var index int
	var tKey string

	for k, v := range mData {
		index = strings.LastIndex(k, "/")
		tKey = k[index+1:]
		if tKey == Basic {
			reData.Basic = config.Basic{}
			if err = utils.LoadYaml(v, &reData.Basic); err != nil {
				break
			}
		} else if tKey == StopService {
			reData.StopService = config.StopService{}
			if err = utils.LoadYaml(v, &reData.StopService); err != nil {
				break
			}
		} else if tKey == Switchdb {
			reData.Switch = config.SwitchDB{}
			if err = utils.LoadYaml(v, &reData.Switch); err != nil {
				break
			}
		} else if tKey == HostGroups {
			reData.HostGroups = make([]config.HostGroup, 0)
			if err = utils.LoadYaml(v, &reData.HostGroups); err != nil {
				break
			}
		} else if tKey == HostGroupClusters {
			reData.HostGroupClusters = make([]config.HostGroupCluster, 0)
			if err = utils.LoadYaml(v, &reData.HostGroupClusters); err != nil {
				break
			}
		} else if tKey == SchemaConfigs {
			reData.SchemaConfigs = make([]config.SchemaConfig, 0)
			if err = utils.LoadYaml(v, &reData.SchemaConfigs); err != nil {
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return utils.UnLoadYaml(&reData)
}

//SetShardConf 存储分库分表
func (s *Store) SetShardConf(clusterName string, data []byte) error {
	if len(clusterName) == 0 || len(data) == 0 {
		return parameterIsNil
	}

	lists, err := s.GetClusterList()
	if err != nil {
		return err
	}
	if len(lists) == 0 {
		lists = []string{clusterName}
	} else {
		repetition := false
		for i := range lists {
			if lists[i] == clusterName {
				repetition = true
			}
		}
		if !repetition {
			lists = append(lists, clusterName)
		}
	}

	strList := strings.Join(lists, ",")

	conf := config.Config{}
	if err = utils.LoadYaml(data, &conf); err != nil {
		return err
	}

	confMap := make(map[string]interface{})
	confMap[Basic] = conf.Basic
	confMap[StopService] = conf.StopService
	confMap[Switchdb] = conf.Switch
	confMap[HostGroups] = conf.HostGroups
	confMap[HostGroupClusters] = conf.HostGroupClusters
	confMap[SchemaConfigs] = conf.SchemaConfigs

	sVersion, err := s.GetNextVersion(clusterName, true)
	if err != nil {
		return err
	}

	opCmds := make([]clientv3.Op, 0)
	// 写入集群列表
	opCmds = append(opCmds, clientv3.OpPut(ClusterList, strList))

	opCmds = append(opCmds, clientv3.OpPut(fmt.Sprintf(ConfigVersion, clusterName), sVersion))

	for key, value := range confMap {
		eKey := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, key)
		data, err = utils.UnLoadYaml(value)
		if err != nil {
			break
		}

		opCmds = append(opCmds, clientv3.OpPut(eKey, string(data)))
	}
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	_, err = s.cli.Txn(ctx).Then(opCmds...).Commit()
	cancel()

	return err
}

// GetClusterList
func (s *Store) GetClusterList() ([]string, error) {
	data, err := s.GetData(ClusterList)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	return strings.Split(string(data), ","), nil
}

//GetAllReadIP 获取Hostgroups中所有的read
func (s *Store) GetAllReadIP(clusterName string) ([]string, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}

	retStr := make([]string, 0)
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return retStr, err
	}

	t := make([]config.HostGroup, 0) //[]config.HostGroup{}

	if err = utils.LoadYaml(data, &t); err != nil {
		return retStr, err
	}

	for _, value := range t {
		if len(value.Read) != 0 {
			retStr = append(retStr, value.Read)
		}
	}

	return retStr, nil
}

//GetReadIPs 获取Hostgroups中的read
func (s *Store) GetReadIPs(clusterName, hostGroup string) (string, error) {
	if len(clusterName) == 0 || len(hostGroup) == 0 {
		return "", parameterIsNil
	}
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return "", err
	}

	t := make([]config.HostGroup, 0) //[]config.HostGroup{}
	if err = utils.LoadYaml(data, &t); err != nil {
		return "", err
	}

	for _, value := range t {
		if value.Name == hostGroup {
			return value.Read, nil
		}
	}

	return "", errors.New("can't find read ip")
}

//GetReadIPs 获取Hostgroups中的指定的hostgroup的write
func (s *Store) GetWriteIPs(clusterName, hostGroup string) (string, error) {
	if len(clusterName) == 0 || len(hostGroup) == 0 {
		return "", parameterIsNil
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return "", err
	}

	t := []config.HostGroup{}
	if err = utils.LoadYaml(data, &t); err != nil {
		return "", err
	}

	for _, value := range t {
		if value.Name == hostGroup {
			return value.Write, nil
		}
	}

	return "", errors.New("can't find write ip")
}

//GetAllWriteIP 获取Hostgroups中的所有write
func (s *Store) GetAllWriteIP(clusterName string) ([]string, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}
	retStr := make([]string, 0)
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return retStr, err
	}

	t := []config.HostGroup{}

	if err = utils.LoadYaml(data, &t); err != nil {
		return retStr, err
	}

	for _, value := range t {
		if value.Write != "" {
			retStr = append(retStr, value.Write)
		}
	}

	return retStr, nil
}

//CreateSnapshot 创建配置快照
func (s *Store) CreateSnapshot(clusterName, name string) error {
	if len(clusterName) == 0 || len(name) == 0 {
		return parameterIsNil
	}

	lists, err := s.GetSnapshotList(clusterName)
	if err != nil {
		return err
	}

	if len(lists) > 0 {
		for _, v := range lists {
			if v == name {
				return errors.New("snapshot name is exist.")
			}
		}
	}

	key := fmt.Sprintf(ShardPrefix, clusterName, DefaultConfig)
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}

	if resp.Count == 0 {
		return errors.New("default config is not exist.")
	}

	ops := make([]clientv3.Op, 0)
	for _, value := range resp.Kvs {
		ops = append(ops, clientv3.OpPut(strings.Replace(string(value.Key), DefaultConfig, name, 1), string(value.Value)))
	}

	lists = append(lists, name)

	listValue := strings.Join(lists, ",")

	ops = append(ops, clientv3.OpPut(fmt.Sprintf(ShardList, clusterName), listValue))

	ctx, cancel = context.WithTimeout(context.Background(), etcdTimeOut)
	_, err = s.cli.Txn(ctx).Then(ops...).Commit()
	cancel()

	return err
}

//GetSnapshotList 获取分库分表的指定的快照列表
func (s *Store) GetSnapshotList(clusterName string) ([]string, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}

	retData := make([]string, 0)

	key := fmt.Sprintf(ShardList, clusterName)
	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, key)
	cancel()
	if err != nil {
		return retData, err
	}
	if resp.Count == 0 {
		return retData, nil
	}

	retData = strings.Split(string(resp.Kvs[0].Value), ",")
	return retData, nil
}

//GetClusterStatus 获取整个proxy集群的状态
func (s *Store) GetClusterStatus(clusterName string) (string, error) {
	if len(clusterName) == 0 {
		return "", parameterIsNil
	}

	key := fmt.Sprintf(StatusPrefix, clusterName)
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return "", err
	}

	retStr := models.StatusUp
	//如果没有proxy集群注册到keeper则认为整个集群是正常up状态
	if resp.Count == 0 {
		return retStr, nil
	}

	var infor models.StatusInfo
	for _, value := range resp.Kvs {
		infor = models.StatusInfo{}
		if err = utils.LoadJSON(value.Value, &infor); err != nil {
			return "", err
		}
		if infor.Status != models.StatusUp {
			return "", errors.New(fmt.Sprintf("%s:%s", value.Key, value.Value))
		}
	}

	return retStr, nil
}

//GetAllStatus 获取所有Proxy的状态信息
func (s *Store) GetAllStatus(clusterName string) (map[string]string, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}

	key := fmt.Sprintf(StatusPrefix, clusterName)
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	resp, err := s.cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	retMap := make(map[string]string)
	infor := models.StatusInfo{}
	for _, value := range resp.Kvs {
		key := string(value.Key)
		index := strings.LastIndex(key, "/")
		if index > 0 {
			key = key[index+1:]
			infor = models.StatusInfo{}
			if err = utils.LoadJSON(value.Value, &infor); err != nil {
				continue
			}
			retMap[key] = infor.Status
		}
	}

	return retMap, nil
}

//getProxyAddrs 根据clusterName获取所有的Proxy的addr
func (s *Store) GetProxyAddrs(clusterName string) ([]string, error) {
	addrs := []string{}
	cli := s.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	defer cancel()
	resp, err := cli.Get(ctx, fmt.Sprintf(StatusPrefix, clusterName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}
	for _, kv := range resp.Kvs {
		addr := strings.TrimPrefix(string(kv.Key), fmt.Sprintf(StatusPrefix, clusterName))
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

//返回指定的hostgroup
func (s *Store) GetHostGroupWithName(clusterName, hostGroupName string) (hostGroup *config.HostGroup, err error) {
	if len(clusterName) == 0 {
		err = parameterIsNil
		return
	}
	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return
	}
	hostgroups := []config.HostGroup{}
	err = utils.LoadYaml(data, &hostgroups)
	if err != nil {
		return
	}
	for _, hg := range hostgroups {
		if hg.Name == hostGroupName {
			hostGroup = &hg
			return
		}
	}
	return
}

// Set hostGroup activeWrite
func (s *Store) SetHostGroupActiveWrite(clusterName, hostGroupName string, activeWrite int) (err error) {
	if len(clusterName) == 0 || len(hostGroupName) == 0 {
		err = parameterIsNil
		return
	}

	key := fmt.Sprintf(ShardConfig, clusterName, DefaultConfig, HostGroups)
	data, err := s.GetData(key)
	if err != nil {
		return
	}
	hostGroups := []config.HostGroup{}
	err = utils.LoadYaml(data, &hostGroups)
	if err != nil {
		return
	}
	var index int
	for i, hostGroup := range hostGroups {
		if hostGroup.Name == hostGroupName {
			index = i + 1
			break
		}

	}
	if index == 0 {
		err = errors.Errorf("hostGroup %s not found ", hostGroupName)
		return
	}
	hostGroups[index-1].ActiveWrite = activeWrite
	data, err = yaml.Marshal(&hostGroups)
	if err != nil {
		return
	}
	err = s.PutData(key, string(data))
	if err != nil {
		return
	}
	return
}
