package backend

import (
	"net"
	"time"

	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

var (
	pingPeriod  = int64(time.Second * 16)
	idGenerator = new(util.IdGenerator)
)

//proxy <-> mysql server
type MysqlConn struct {
	ID       int64
	conn     net.Conn
	pkg      *mysql.PacketIO
	addr     string
	user     string
	password string
	db       string

	capability uint32
	status     uint16

	collation mysql.CollationId
	charset   string
	salt      []byte

	pkgErr error

	pushTimestamp       int64
	lastActiveTimestamp int64
}

type BackendConn struct {
	*MysqlConn
	dbPool *DBPool
}

func (bc *BackendConn) GetDBPool() *DBPool {
	return bc.dbPool
}

func NewConnCache() *MysqlConn {
	mysqlConn := new(MysqlConn)
	mysqlConn.ID = idGenerator.NextId()
	return mysqlConn
}

// 如果连接是从连接池里获取的，Close了别忘记把MySqlConn的壳还给DBPool.deadConns
func (c *MysqlConn) Close() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.salt = nil
		c.pkgErr = nil
	}

	return nil
}
