package frontend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/glog"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"net"
	"runtime"
)

var defaultCapability = mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_LONG_FLAG |
	mysql.CLIENT_CONNECT_WITH_DB | mysql.CLIENT_PROTOCOL_41 |
	mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_SECURE_CONNECTION

func (e *ProxyEngine) newConn(c net.Conn) {
	conn := NewFrontendConn(c, e)
	defer func() {
		err := recover()
		if err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)] //获得当前goroutine的stacktrace
			glog.Glog.Errorf("newConn error. remoteAddr:%v, stack:%s", c.RemoteAddr(), string(buf))
		}

		conn.Close(false)
	}()

	if err := conn.Handshake(); err != nil {
		glog.Glog.Errorf("Handshake failed. IP:[%s], err:%+v", conn.GetRemoteIP(), err)
		conn.writeError(err)
		conn.Close(false)
		return
	}

	conn.Start()
}

func (fc *FrontendConn) Handshake() error {
	if err := fc.writeInitialHandshake(); err != nil {
		return err
	}

	if err := fc.readHandshakeResponse(); err != nil {
		return err
	}

	if err := fc.writeOK(nil); err != nil {
		return err
	}

	fc.pkg.Sequence = 0
	return nil
}

func (fc *FrontendConn) writeInitialHandshake() error {
	//data := make([]byte, 4, 128)
	data := fc.handlerArena.AllocWithLen(4, 128)

	//min version 10
	data = append(data, 10)

	//server version[00]
	data = append(data, mysql.ServerVersion...)
	data = append(data, 0)

	//connection id
	data = append(data, byte(fc.connectionId), byte(fc.connectionId>>8), byte(fc.connectionId>>16), byte(fc.connectionId>>24))

	//auth-plugin-data-part-1
	data = append(data, fc.salt[0:8]...)

	//filter [00]
	data = append(data, 0)

	//capability flag lower 2 bytes, using default capability here
	data = append(data, byte(defaultCapability), byte(defaultCapability>>8))

	//charset, utf-8 default
	data = append(data, uint8(mysql.DEFAULT_COLLATION_ID))

	//status
	data = append(data, byte(fc.status), byte(fc.status>>8))

	//below 13 byte may not be used
	//capability flag upper 2 bytes, using default capability here
	data = append(data, byte(defaultCapability>>16), byte(defaultCapability>>24))

	//filter [0x15], for wireshark dump, value is 0x15
	data = append(data, 0x15)

	//reserved 10 [00]
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	//auth-plugin-data-part-2
	data = append(data, fc.salt[8:]...)

	//filter [00]
	data = append(data, 0)

	return fc.writePacket(data)
}

func (fc *FrontendConn) readHandshakeResponse() error {
	data, err := fc.readPacket()

	if err != nil {
		return err
	}

	pos := 0

	//capability
	fc.capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4

	//skip max packet size
	pos += 4

	//charset, skip, if you want to use another charset, use set names
	//c.collation = CollationId(data[pos])
	pos++

	//skip reserved 23[00]
	pos += 23

	//user name
	fc.user = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])

	pos += len(fc.user) + 1

	//auth length and auth
	authLen := int(data[pos])
	pos++
	auth := data[pos : pos+authLen]

	checkAuth := mysql.CalcPassword(fc.salt, []byte(fc.proxy.password))
	if fc.user != fc.proxy.user || !bytes.Equal(auth, checkAuth) {
		return mysql.NewDefaultError(mysql.ER_ACCESS_DENIED_ERROR, fc.user, fc.c.RemoteAddr().String(), "Yes")
	}

	pos += authLen

	var db string
	if fc.capability&mysql.CLIENT_CONNECT_WITH_DB > 0 {
		if len(data[pos:]) == 0 {
			return nil
		}

		db = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(fc.logicDb) + 1

	}
	if len(db) != 0 && fc.proxy.GetSchemaNode(db) == nil {
		return fmt.Errorf("unknown database '%s'", db)
	}

	fc.logicDb = db

	return nil
}
