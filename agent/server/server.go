package server

import (
	"database/sql"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"git.2dfire.net/zerodb/agent/glog"

	"git.2dfire.net/zerodb/agent/metrics"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/agent"
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mySQLDie = iota
	mySQLLive
)

//Server 服务类
type Server struct {
	mysqlConn *sql.DB             //agent与mysql的ipc
	keeperCli keeper.KeeperClient //与keeper通信的rpc
	conf      *Config
	ticker    *time.Ticker
	done      chan bool
}

//NewServer 创建服务
func NewServer(conf *Config) (*Server, error) {
	s := new(Server)
	var err error

	s.conf = conf

	s.mysqlConn, err = sql.Open("mysql", s.conf.MysqlConn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var conn *grpc.ClientConn

	conn, err = grpc.Dial(s.conf.KeeperAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.keeperCli = keeper.NewKeeperClient(conn)

	s.done = make(chan bool)
	s.ticker = time.NewTicker(HearbeatTime)

	return s, nil
}

func (s *Server) startRPC() {
	ln, err := net.Listen("tcp", s.conf.RPCServer)
	if err != nil {
		glog.Glog.Fatalln(err)
	}
	rpcServer := grpc.NewServer()
	agent.RegisterAgentServer(rpcServer, s)
	reflection.Register(rpcServer)

	if err := rpcServer.Serve(ln); err != nil {
		glog.Glog.Fatalf("agent rpc server failed: %v", err)
	}
}

//Run 服务主循环
func (s *Server) Run() error {
	go s.startRPC()

	go func() {
		if err := metrics.Run(s.conf.MetricsServer); err != nil {
			glog.Glog.Fatalln(err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":8088", nil); err != nil {
			glog.Glog.Fatalln(err)
		}
	}()

	go func() {
		for {
			select {
			case done := <-s.done:
				if done {
					s.ticker.Stop()
					return
				}
			case <-s.ticker.C:
				s.heartbeat()
			}
		}
	}()
	select {}
}

//Close 关闭服务
func (s *Server) Close() {
	s.mysqlConn.Close()
	s.done <- true
}

//mysqlInfo 收集mysql信息
func (s *Server) mysqlInfo() (*keeper.Mysql, error) {
	in := new(keeper.Mysql)

	sql := "select 1"
	if _, err := s.sqlQuery(sql); err != nil {
		in.Status = mySQLDie
		return in, nil
	}
	in.Status = mySQLLive

	sql = "show status like 'Threads_connected'"
	mData, err := s.sqlQuery(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if mData != nil {
		in.Connected, err = strconv.ParseInt(mData["Value"].(string), 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	sql = "show status like 'Memory_used'"
	mData, err = s.sqlQuery(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if mData != nil {
		in.Memory, err = strconv.ParseInt(mData["Value"].(string), 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return in, nil
}

func (s *Server) sqlQuery(query string, args ...interface{}) (map[string]interface{}, error) {
	stmt, err := s.mysqlConn.Prepare(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	values := make([]interface{}, len(cols))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, errors.WithStack(err)
		}
		result := make(map[string]interface{}, len(cols))
		for ii, key := range cols {
			if scanArgs[ii] == nil {
				continue
			}
			value := reflect.Indirect(reflect.ValueOf(scanArgs[ii]))
			if value.Elem().Kind() == reflect.Slice {
				result[key] = string(value.Interface().([]byte))
			} else {
				result[key] = value.Interface()
			}
		}
		return result, nil
	}
	return nil, nil
}
