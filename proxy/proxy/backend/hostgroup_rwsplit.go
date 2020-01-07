package backend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"strconv"
	"strings"
)

// 读写分离
func (node *HostGroupNode) GetWriteConn() (*BackendConn, error) {
	dbPool := node.Write[node.activedWriteIndex]
	if dbPool == nil {
		return nil, errcode.ErrNoWriteDBPool
	}

	return dbPool.BorrowConn()
}

func (node *HostGroupNode) GetReadConn() (*BackendConn, error) {
	node.Lock()
	dbPool, err := node.GetNextRead()
	node.Unlock()
	if err != nil {
		return nil, err
	}

	if dbPool == nil {
		return nil, errcode.ErrNoReadDBPool
	}

	return dbPool.BorrowConn()
}

func (node *HostGroupNode) AddRead(addr string) error {
	var dbPool *DBPool
	var weight int
	var err error
	if len(addr) == 0 {
		return errcode.ErrAddressNull
	}
	node.Lock()
	defer node.Unlock()
	for _, v := range node.Read {
		if strings.Split(v.addr, WeightSplit)[0] == strings.Split(addr, WeightSplit)[0] {
			return errcode.ErrReadExist
		}
	}
	addrAndWeight := strings.Split(addr, WeightSplit)
	if len(addrAndWeight) == 2 {
		weight, err = strconv.Atoi(addrAndWeight[1])
		if err != nil {
			return err
		}
	} else {
		weight = 1
	}

	if dbPool, err = node.InitPool(addrAndWeight[0], "R"); err != nil {
		return err
	} else {
		node.ReadWeights = append(node.ReadWeights, weight)
		node.Read = append(node.Read, dbPool)
		node.InitBalancer()
		return nil
	}
}

func (node *HostGroupNode) DeleteRead(addr string) error {
	var i int
	node.Lock()
	defer node.Unlock()
	readCount := len(node.Read)
	if readCount == 0 {
		return errcode.ErrNoReadDBPool
	}
	for i = 0; i < readCount; i++ {
		if node.Read[i].addr == addr {
			break
		}
	}
	if i == readCount {
		return errcode.ErrReadConfNotExist
	}
	if readCount == 1 {
		node.Read = nil
		node.ReadWeights = nil
		node.RoundRobinQ = nil
		return nil
	}

	s := make([]*DBPool, 0, readCount-1)
	sw := make([]int, 0, readCount-1)
	for i = 0; i < readCount; i++ {
		if node.Read[i].addr != addr {
			s = append(s, node.Read[i])
			sw = append(sw, node.ReadWeights[i])
		}
	}

	node.Read = s
	node.ReadWeights = sw
	node.InitBalancer()
	return nil
}

func (node *HostGroupNode) InitWrite(writeStr string, initPool bool) error {
	var dbPool *DBPool
	var err error

	if len(writeStr) == 0 {
		return fmt.Errorf("writeStr should not be empty")
	}
	writeStr = strings.Trim(writeStr, HostSplit)
	writeArray := strings.Split(writeStr, HostSplit)
	count := len(writeArray)
	node.Write = make([]*DBPool, 0, count)

	//parse addr
	for i := 0; i < count; i++ {
		if initPool {
			if dbPool, err = node.InitPool(writeArray[i], "W"); err != nil {
				return err
			}
		} else {
			glog.Glog.Warnf("initPool for [%v] is false. this is test only", writeArray[i])
		}
		node.Write = append(node.Write, dbPool)
	}
	return nil
}

//readStr(127.0.0.1:3306@2,192.168.0.12:3306@3)
func (node *HostGroupNode) InitRead(readStr string, initPool bool) error {
	var dbPool *DBPool
	var weight int
	var err error

	if len(readStr) == 0 {
		//glog.ProxyLog.Warnf("readStr is empty for hostGroup:[%s]", node.Cfg.Name)
		return nil
	}
	readStr = strings.Trim(readStr, HostSplit)
	readArray := strings.Split(readStr, HostSplit)
	count := len(readArray)
	node.Read = make([]*DBPool, 0, count)
	node.ReadWeights = make([]int, 0, count)

	//parse addr and weight
	for i := 0; i < count; i++ {
		addrAndWeight := strings.Split(readArray[i], WeightSplit)
		if len(addrAndWeight) == 2 {
			weight, err = strconv.Atoi(addrAndWeight[1])
			if err != nil {
				return err
			}
		} else {
			weight = 1
		}
		node.ReadWeights = append(node.ReadWeights, weight)
		if initPool {
			if dbPool, err = node.InitPool(addrAndWeight[0], "R"); err != nil {
				return err
			}
		} else {
			glog.Glog.Warnf("initPool for [%v] is false. this is test only", addrAndWeight[0])
		}
		node.Read = append(node.Read, dbPool)
	}
	node.InitBalancer()
	return nil
}
