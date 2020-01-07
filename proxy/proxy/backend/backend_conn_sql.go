package backend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"github.com/pkg/errors"
	"net"
	"strings"
	"time"
)

func (c *MysqlConn) ConnectMysql(addr string, user string, password string, db string) error {
	c.addr = addr
	c.user = user
	c.password = password
	c.db = db

	//use utf8
	c.collation = mysql.DEFAULT_COLLATION_ID
	c.charset = mysql.DEFAULT_CHARSET

	return c.doConnectMysql()
}

func (c *MysqlConn) doConnectMysql() error {
	if c.conn != nil {
		c.conn.Close()
	}

	n := "tcp"
	if strings.Contains(c.addr, "/") {
		n = "unix"
	}

	netConn, err := net.DialTimeout(n, c.addr, time.Second*3)
	if err != nil {
		return errors.WithMessage(err, "addr:"+c.addr)
	}

	tcpConn := netConn.(*net.TCPConn)

	tcpConn.SetNoDelay(false)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetWriteBuffer(16 * 1024)
	tcpConn.SetReadBuffer(16 * 1024)
	c.conn = tcpConn
	c.pkg = mysql.NewPacketIO(tcpConn)

	if err := c.readInitialHandshake(); err != nil {
		c.conn.Close()
		return err
	}

	if err := c.writeAuthHandshake(); err != nil {
		c.conn.Close()
		return err
	}

	if _, err := c.readOK(); err != nil {
		c.conn.Close()
		return err
	}

	//we must always use autocommit
	if !c.IsAutoCommit() {
		if _, err := c.exec("set autocommit = 1"); err != nil {
			c.conn.Close()
			return err
		}
	}

	c.recordTime()

	return nil
}

func (c *MysqlConn) Ping() error {
	if err := c.writeCommand(mysql.COM_PING); err != nil {
		return err
	}

	if _, err := c.readOK(); err != nil {
		return err
	}

	c.recordTime()

	return nil
}

func (c *MysqlConn) PingTimeout() error {
	if err := c.writeCommandWithDeadline(mysql.COM_PING); err != nil {
		return err
	}

	if _, err := c.readOK(); err != nil {
		return err
	}

	c.recordTime()

	return nil
}

func (c *MysqlConn) UseDB(dbName string) error {
	if c.db == dbName || len(dbName) == 0 {
		return nil
	}

	if err := c.writeCommandStr(mysql.COM_INIT_DB, dbName); err != nil {
		return err
	}

	if _, err := c.readOK(); err != nil {
		return err
	}

	c.db = dbName
	return nil
}

func (c *MysqlConn) GetDB() string {
	return c.db
}

func (c *MysqlConn) GetAddr() string {
	return c.addr
}

func (c *MysqlConn) Execute(command string, args ...interface{}) (*mysql.Result, error) {
	if len(args) == 0 {
		return c.exec(command)
	}

	s, err := c.Prepare(command)
	if err != nil {
		return nil, err
	}

	var r *mysql.Result
	r, err = s.Execute(args...)
	s.Close()
	return r, err
}

func (c *MysqlConn) ExecuteWithAlloc(alloc *util.FastYetSafeAllocator, command string, args ...interface{}) (*mysql.Result, error) {
	if len(args) == 0 {
		return c.execWithAlloc(alloc, command)
	}

	s, err := c.Prepare(command)
	if err != nil {
		return nil, err
	}

	var r *mysql.Result
	r, err = s.Execute(args...)
	s.Close()
	return r, err
}

func (c *MysqlConn) exec(query string) (*mysql.Result, error) {
	if err := c.writeCommandStr(mysql.COM_QUERY, query); err != nil {
		return nil, err
	}

	return c.readResult(false)
}

func (c *MysqlConn) execWithAlloc(alloc *util.FastYetSafeAllocator, query string) (*mysql.Result, error) {
	if err := c.writeCommandStrWithAllocator(alloc, mysql.COM_QUERY, query); err != nil {
		return nil, err
	}

	return c.readResult(false)
}

func (c *MysqlConn) ClosePrepare(id uint32) error {
	return c.writeCommandUint32(mysql.COM_STMT_CLOSE, id)
}

func (c *MysqlConn) SetCharset(charset string, collation mysql.CollationId) error {
	charset = strings.Trim(charset, "\"'`")

	if collation == 0 {
		collation = mysql.CollationNames[mysql.Charsets[charset]]
	}

	if c.charset == charset && c.collation == collation {
		return nil
	}

	_, ok := mysql.CharsetIds[charset]
	if !ok {
		return fmt.Errorf("invalid charset %s", charset)
	}

	_, ok = mysql.Collations[collation]
	if !ok {
		return fmt.Errorf("invalid collation %v", collation)
	}

	if _, err := c.exec(fmt.Sprintf("SET NAMES %s COLLATE %s", charset, mysql.Collations[collation])); err != nil {
		return err
	}
	c.collation = collation
	c.charset = charset
	return nil
}
