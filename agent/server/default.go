package server

import "time"

const (
	//TimeOut 网络超时
	TimeOut = 3 * time.Second

	//HearbeatTime 心跳频率
	HearbeatTime = 5 * time.Second
)
