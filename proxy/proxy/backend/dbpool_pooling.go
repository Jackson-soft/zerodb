package backend

// 向连接池「还」连接。是外部归还连接池的唯一方式
func (pool *DBPool) ReturnConn(b *BackendConn) {
	if b != nil && b.MysqlConn != nil {
		b.MysqlConn.recordTime()

		if b.MysqlConn.pkgErr != nil {
			b.dbPool.returnConnToDeadChan(b.MysqlConn)
		} else {
			b.dbPool.returnConnToLiveChan(b.MysqlConn, nil)
		}
		id := b.MysqlConn.ID
		b.MysqlConn = nil

		b.dbPool.ActiveConns.Delete(id)
	}
}

// 向连接池「借」连接。是外部获取连接的唯一方式
// TODO nanxing 如果MySQL执行慢了，并且Proxy并发量很大，那么，借出去的连接迟迟不能归还，会让Proxy大量创建MySQL连接，更加拖慢MySQL
// TODO nanxing 需要设计一个testOnBorrow，因为一个大SQL需要验证每个连接，不然就其它的连接就白执行了
func (pool *DBPool) BorrowConn() (*BackendConn, error) {
	c, err := pool.popConn()
	if err != nil {
		return nil, err
	}
	pool.ActiveConns.Store(c.ID, c)
	c.recordTime()

	return &BackendConn{c, pool}, nil
}

func (pool *DBPool) RemainSize() int {
	return len(pool.idleConns) + len(pool.deadConns) // 活的连接 + 死的连接
}

func (pool *DBPool) RemainIdleSize() int {
	return len(pool.idleConns)
}

func (pool *DBPool) RemainDeadSize() int {
	return len(pool.deadConns)
}

func (pool *DBPool) ActiveSize() int {
	return pool.maxConnNum - pool.RemainSize()
}
