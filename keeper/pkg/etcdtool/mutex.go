package etcdtool

import (
	"context"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

const (
	defaultTTL = 10
	lockValue  = "lock"
)

var (
	ErrLockAllreadyHold = errors.New("lock already hold by other client.")
)

//Mutex 分布式锁
type Mutex struct {
	cli *clientv3.Client
	ttl int64
	key string
}

//NewMutex new mutex
func (s *Store) NewMutex(key string, ttl int64) (*Mutex, error) {
	if len(key) == 0 {
		return nil, parameterIsNil
	}
	m := new(Mutex)
	if ttl == 0 {
		m.ttl = defaultTTL
	} else {
		m.ttl = ttl
	}
	m.cli = s.cli
	m.key = key
	return m, nil
}

func (m *Mutex) Lock() error {
	resp, err := m.cli.Grant(context.TODO(), m.ttl)
	if err != nil {
		return err
	}

	//如果key不存在返回createRevision为0.
	cmp := clientv3.Compare(clientv3.CreateRevision(m.key), "=", 0)
	put := clientv3.OpPut(m.key, lockValue, clientv3.WithLease(resp.ID))
	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	//如果key不存在resp.Succeeded为true，否则为false
	txResp, err := m.cli.Txn(ctx).If(cmp).Then(put).Commit()
	cancel()
	if err != nil {
		return err
	}
	if !txResp.Succeeded {
		//如果CreateRevision不为0表示etcd中已经有这个key存在直接返回ErrLockAllreadyHold
		return ErrLockAllreadyHold
	}
	return nil
}

func (m *Mutex) UnLock() error {
	ctx, cancel := context.WithTimeout(context.TODO(), etcdTimeOut)
	_, err := m.cli.Delete(ctx, m.key)
	cancel()
	return err
}
