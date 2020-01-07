package backend

import (
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"time"
)

// 高于水位时 告警（WatermarkCheckTask） + 驱逐（IdleConnEvictTask）
// 低于水位时 补充（IdleConnSupplyTask） + 保活（IdleConnKeepAliveTask）

// 高水位监控
func (pool *DBPool) WatermarkCheckTask() {
	defer func() {
		err := recover()
		if err != nil {
			glog.Glog.Errorf("WatermarkCheckTask panic occurs: Pool status: %d err: %v", pool.status, err)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	//glog.WatermarkLog.Infof("DbPool %s WatermarkCheckTask starts normally.", pool.GetName())
	for {

		select {
		case <-ticker.C:
			if pool.status == Cleaned || pool.status == Cleaning {
				// 需要初始化
				continue
			}
			if pool.status == Down {
				break
			}

			remainSize := pool.RemainSize()
			warnSize := pool.maxConnNum / 4
			if remainSize < warnSize {
				//glog.WatermarkLog.Errorf("DBPool watermark warning. only %d conn(s) in pool %s, active: %d, idle: %d, dead: %d",
				//	remainSize,
				//	pool.GetName(),
				//	pool.ActiveSize(),
				//	pool.RemainIdleSize(),
				//	pool.RemainDeadSize())
			}
		}
	}

	glog.Glog.Infof("DbPool %s WatermarkCheckTask stops normally.", pool.GetName())
}

func (pool *DBPool) IdleConnEvictTask() {
	defer func() {
		err := recover()
		if err != nil {
			glog.Glog.Errorf("IdleConnEvictTask panic occurs: Pool status: %d, err: %v", pool.status, err)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	//glog.WatermarkLog.Infof("DbPool %s IdleConnEvictTask starts normally.", pool.GetName())
	for {

		select {
		case <-ticker.C:
			co := <-pool.idleConns

			if pool.status == Cleaned || pool.status == Cleaning {
				// 需要初始化
				continue
			}
			if pool.status == Down {
				break
			}

			if time.Now().Unix()-co.lastActiveTimestamp >= pool.maxIdleLiveTime && pool.RemainIdleSize() > pool.InitConnNum {
				co.Close()
				//glog.WatermarkLog.Infof("[Evict]. pool: %s. active: %d, idle: %d, dead: %d",
				//	pool.GetName(),
				//	pool.ActiveSize(),
				//	pool.RemainIdleSize(),
				//	pool.RemainDeadSize())
				pool.deadConns <- co
			} else {
				pool.idleConns <- co
			}
		}
	}

	glog.Glog.Infof("DbPool %s IdleConnEvictTask stops normally.", pool.GetName())
}

func (pool *DBPool) IdleConnSupplyTask() {
	defer func() {
		err := recover()
		if err != nil {
			glog.Glog.Errorf("IdleConnSupplyTask panic occurs: Pool status: %d, err: %v", pool.status, err)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	//glog.WatermarkLog.Infof("DbPool %s IdleConnSupplyTask starts normally.", pool.GetName())
	for {

		select {
		case <-ticker.C:
			if pool.status == Cleaned || pool.status == Cleaning {
				// 需要初始化
				continue
			}
			if pool.status == Down {
				break
			}

			// 由于有多个任务同时在工作，有时会多补充几个连接
			for pool.RemainIdleSize() < pool.InitConnNum {
				conn, err := pool.getConnFromDead()
				if err != nil {
					//glog.WatermarkLog.Errorf("DbPool %s IdleConnSupplyTask: Can't activate mysql conn from dead. active: %d, idle: %d, dead: %d",
					//	pool.GetName(),
					//	pool.ActiveSize(),
					//	pool.RemainIdleSize(),
					//	pool.RemainDeadSize())
					break
				}
				conn.pushTimestamp = time.Now().Unix()
				pool.idleConns <- conn
				//glog.WatermarkLog.Infof("[Supply]. pool: %s. active: %d, idle: %d, dead: %d",
				//	pool.GetName(),
				//	pool.ActiveSize(),
				//	pool.RemainIdleSize(),
				//	pool.RemainDeadSize())
			}
		}
	}

	glog.Glog.Infof("DbPool %s IdleConnSupplyTask stops normally.", pool.GetName())
}

// 由于MySQL有「interactive_timeout」的机制，只要连接太久没有活动，MySQL会自动清理掉该连接，那Proxy保持的那个idle其实不再空闲，是CLOSE-WAIT的状态
// 所以这个定时器可以定时清理CLOSE-WAIT状态的连接，频率最好比「Supply」高点
func (pool *DBPool) IdleConnKeepAliveTask() {
	defer func() {
		err := recover()
		if err != nil {
			glog.Glog.Errorf("IdleConnKeepAliveTask panic occurs: Pool status: %d, err: %v", pool.status, err)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	//glog.WatermarkLog.Infof("DbPool %s IdleConnKeepAliveTask starts normally.", pool.GetName())
	for {

		select {
		case <-ticker.C:
			if pool.status == Cleaned || pool.status == Cleaning {
				// 需要初始化
				continue
			}
			if pool.status == Down {
				break
			}

			if pool.RemainIdleSize() <= pool.InitConnNum {
				co := <-pool.idleConns

				if time.Now().Unix()-co.lastActiveTimestamp >= pool.maxIdleLiveTime {
					err := co.Ping()
					if err != nil {
						pool.returnConnToDeadChan(co)
						//glog.WatermarkLog.Infof("[KeepAlive] failed. pool: %s. active: %d, idle: %d, dead: %d",
						//	pool.GetName(),
						//	pool.ActiveSize(),
						//	pool.RemainIdleSize(),
						//	pool.RemainDeadSize())
					} else {
						pool.idleConns <- co
						//glog.WatermarkLog.Infof("[KeepAlive] pool: %s. active: %d, idle: %d, dead: %d",
						//	pool.GetName(),
						//	pool.ActiveSize(),
						//	pool.RemainIdleSize(),
						//	pool.RemainDeadSize())
					}
				} else {
					// 本来就健康，不需要保活
					pool.idleConns <- co
				}
			}
		}
	}

	glog.Glog.Infof("DbPool %s IdleConnKeepAliveTask stops normally.", pool.GetName())
}
