// Copyright 2017 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package concurrency_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"git.2dfire.net/zerodb/keeper/pkg/concurrency"
	"go.etcd.io/etcd/clientv3"
)

var (
	wa sync.WaitGroup
)

func ExampleMutex_Lock() {
	for i := 0; i < 3; i++ {
		wa.Add(1)
		go testLock(i)
	}
	wa.Wait()
	// Output:
	// goroutine0 get lock...
	// goroutine1 can't get lock
	// goroutine2 can't get lock
}

func testLock(g int) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 50,
	})
	if err != nil {
		log.Fatalf("create etcd client %v", err)
	}
	defer cli.Close()
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatalf("create etcd session %v", err)
	}
	defer s.Close()
	m := concurrency.NewMutex(s, "/zerodb/keeper/hostgroup/hostgroup1")
	switch g {
	case 1:
		time.Sleep(time.Second * 3)
	case 2:
		time.Sleep(time.Second * 6)

	}
	err = m.Lock(context.TODO())
	if err != nil {
		if err == concurrency.ErrLockAllreadyHold {
			fmt.Printf("goroutine%d can't get lock\n", g)
		} else {
			log.Fatalln(err)
		}
		wa.Done()
		return
	}
	fmt.Printf("goroutine%d get lock...\n", g)
	time.Sleep(time.Second * 15)
	err = m.Unlock(context.TODO())
	if err != nil {
		log.Println(err)
	}
	wa.Done()
}
