package etcdtool

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.2dfire.net/zerodb/common/utils"
	"go.etcd.io/etcd/clientv3"
)

func (s *Store) BackupLastModify(revsion int64, keys []string) error {
	//只备份 zerodb/shardconfg/clusterName/default/basic、hostgroups、schemaconfigs、stopservice、switch、hostgroupcluters 下的数据
	l := lastModify{}
	for _, k := range keys {
		if IsShardConfKey(k) {
			l.Keys = append(l.Keys, k)
			l.Revision = revsion
		}
	}
	clusterName := GetClusterNameFromKey(keys[0])
	if len(clusterName) == 0 {
		return nil
	}
	//backup  lastModify to etcd
	b, err := utils.UnLoadYaml(&l)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut)
	_, err = s.cli.Put(ctx, fmt.Sprintf(LastModifyPrefix, clusterName), string(b))
	cancel()
	return err
}

func (s *Store) RollBackLastModify(clusterName string) error {
	//如果内存里面不为空
	var l = lastModify{}
	v, err := s.GetData(fmt.Sprintf(LastModifyPrefix, clusterName))
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("etcd backup is nil ")
	}
	if err = utils.LoadYaml(v, &l); err != nil {
		return err
	}

	if l.Revision != 0 {
		for _, key := range l.Keys {
			ctx, cancel := context.WithTimeout(context.Background(), etcdTimeOut*2)
			value, err := s.cli.Get(ctx, key, clientv3.WithRev(l.Revision-1))
			cancel()
			if err != nil {
				return err
			}
			v := string(value.Kvs[0].Value)
			_, err = s.cli.Put(ctx, key, v)
			if err != nil {
				return err
			}
		}

	}
	b, err := utils.UnLoadYaml(lastModify{})
	if err != nil {
		return err
	}
	return s.PutData(fmt.Sprintf(LastModifyPrefix, clusterName), string(b))
}

//TODO 使用正则匹配是最好的方式
func IsShardConfKey(key string) bool {
	var shardConfKey = []string{"basic", "hostgroups", "schemaconfigs", "stopservice", "switch", "hostgroupclusters"}
	for _, k := range shardConfKey {
		if strings.HasSuffix(key, DefaultConfig+"/"+k) {
			return true
		}
	}
	return false
}

//根据key获取clusterName
func GetClusterNameFromKey(key string) string {
	data := strings.Split(key, "/")
	if len(data) < 2 {
		return ""
	}
	return data[2]
}
