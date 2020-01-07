package etcdtool

import (
	"fmt"

	"git.2dfire.net/zerodb/common/utils"
)

// 集群列表是通过心跳机制来实现，目前所有的权重写死为1
type ClusterInfo struct {
	Host   string `json:"host,omitempty"`
	Weight uint64 `json:"weight,omitempty"`
	Port   uint64 `json:"port,omitempty"`
}

// 添加
func (s *Store) AddCluster(clusterName, host string, port, weight uint64) error {
	if len(clusterName) == 0 {
		return parameterIsNil
	}

	key := fmt.Sprintf(ProxyList, clusterName)
	datas, err := s.GetClusters(key)
	if err != nil {
		return err
	}
	if datas == nil {
		datas = make([]ClusterInfo, 0)
	}
	datas = append(datas, ClusterInfo{
		Host:   host,
		Port:   port,
		Weight: weight,
	})
	data, err := utils.UnLoadJSON(&datas)
	if err != nil {
		return err
	}

	return s.PutData(key, string(data))
}

// 删除
func (s *Store) DelCluster(clusterName, host string) error {
	key := fmt.Sprintf(ProxyList, clusterName)
	v, err := s.GetData(key)
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return nil
	}
	datas := make([]ClusterInfo, 0)
	if err = utils.LoadJSON(v, &datas); err != nil {
		return err
	}

	for k, vv := range datas {
		if vv.Host == host {
			datas = append(datas[:k], datas[k+1:]...)
			break
		}
	}

	data, err := utils.UnLoadJSON(&datas)
	if err != nil {
		return err
	}
	return s.PutData(key, string(data))
}

func (s *Store) GetClusters(clusterName string) ([]ClusterInfo, error) {
	if len(clusterName) == 0 {
		return nil, parameterIsNil
	}

	key := fmt.Sprintf(ProxyList, clusterName)
	values, err := s.GetData(key)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	data := make([]ClusterInfo, 0)
	if err = utils.LoadJSON([]byte(values), &data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}
