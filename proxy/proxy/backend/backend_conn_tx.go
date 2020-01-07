package backend

import "git.2dfire.net/zerodb/proxy/proxy/mysql"

func (c *MysqlConn) Begin() error {
	_, err := c.exec("begin")
	return err
}

func (c *MysqlConn) Commit() error {
	_, err := c.exec("commit")
	return err
}

func (c *MysqlConn) Rollback() error {
	_, err := c.exec("rollback")
	return err
}

func (c *MysqlConn) SetAutoCommit(n uint8) error {
	if n == 0 {
		if _, err := c.exec("set autocommit = 0"); err != nil {
			c.conn.Close()
			return err
		}
	} else {
		if _, err := c.exec("set autocommit = 1"); err != nil {
			c.conn.Close()
			return err
		}
	}
	return nil
}

func (c *MysqlConn) IsAutoCommit() bool {
	return c.status&mysql.SERVER_STATUS_AUTOCOMMIT > 0
}

func (c *MysqlConn) IsInTransaction() bool {
	return c.status&mysql.SERVER_STATUS_IN_TRANS > 0
}
