package utils

import (
	"time"
)

const (
	//默认最大重试次数
	DefaultMaxRetries int32 = 10
	//默认backoff 时间间隔单位
	RetryInterval uint64 = 500
)

// RunWithRetry 将会调用f函数按照backoff算法来进行重试
// retryCnt 表示最大重试次数
// backoff backoff 表示时间间隔单位
// f 表示重试函数 返回2个参数,参数一表示是否继续重试,参数二表示产生的错误
func RunWithRetry(retryCnt int, backoff uint64, f func() (bool, error)) (err error) {
	for i := 1; i <= retryCnt; i++ {
		var retryAble bool
		retryAble, err = f()
		if err == nil || !retryAble {
			return err
		}
		sleepTime := time.Duration(backoff*uint64(i)) * time.Millisecond
		time.Sleep(sleepTime)
	}
	return err
}
