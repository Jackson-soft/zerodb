// 参考etcdv3的分布式锁

package concurrency

import (
	"context"

	pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/pkg/errors"
	v3 "go.etcd.io/etcd/clientv3"
)

// Mutex implements the sync Locker interface with etcd
type Mutex struct {
	s     *Session
	pfx   string
	myKey string
	myRev int64
	hdr   *pb.ResponseHeader
}

const (
	key = "lock"
)

var (
	ErrLockAllreadyHold = errors.New("lock already hold by other client")
)

func NewMutex(s *Session, pfx string) *Mutex {
	return &Mutex{s, pfx + "/", "", -1, nil}
}

// Lock 方法锁住mutex
func (m *Mutex) Lock(ctx context.Context) error {
	s := m.s
	client := m.s.Client()
	m.myKey = m.pfx + key
	//如果key不存在返回createRevision为0.
	cmp := v3.Compare(v3.CreateRevision(m.myKey), "=", 0)
	put := v3.OpPut(m.myKey, "", v3.WithLease(s.Lease()))
	//如果key不存在resp.Succeeded为true，否则为false
	resp, err := client.Txn(ctx).If(cmp).Then(put).Commit()

	if err != nil {
		return err
	}
	if !resp.Succeeded {
		//如果CreateRevision不为0表示etcd中已经有这个key存在直接返回ErrLockAllreadyHold
		return ErrLockAllreadyHold
	}
	return nil
}

// Unlock 解锁mutex,删除etcd中的key
func (m *Mutex) Unlock(ctx context.Context) error {
	client := m.s.Client()
	if _, err := client.Delete(ctx, m.myKey); err != nil {
		return err
	}
	m.myKey = "\x00"
	m.myRev = -1
	return nil
}

func (m *Mutex) IsOwner() v3.Cmp {
	return v3.Compare(v3.CreateRevision(m.myKey), "=", m.myRev)
}

func (m *Mutex) Key() string { return m.myKey }

// Header is the response header received from etcd on acquiring the lock.
func (m *Mutex) Header() *pb.ResponseHeader { return m.hdr }

type lockerMutex struct{ *Mutex }

func (lm *lockerMutex) Lock() error {
	client := lm.s.Client()
	if err := lm.Mutex.Lock(client.Ctx()); err != nil {
		return err
	}
	return nil
}
func (lm *lockerMutex) Unlock() error {
	client := lm.s.Client()
	if err := lm.Mutex.Unlock(client.Ctx()); err != nil {
		return err
	}
	return nil
}

// NewLocker creates a sync.Locker backed by an etcd mutex.
func NewLocker(s *Session, pfx string) *lockerMutex {
	return &lockerMutex{NewMutex(s, pfx)}
}
