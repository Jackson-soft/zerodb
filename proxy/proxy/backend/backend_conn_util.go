package backend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"time"
)

func (c *MysqlConn) readPacket() ([]byte, error) {
	d, err := c.pkg.ReadPacket()
	c.pkgErr = err
	return d, err
}

func (c *MysqlConn) writePacket(data []byte) error {
	err := c.pkg.WritePacket(data)
	c.pkgErr = err
	return err
}

func (c *MysqlConn) readInitialHandshake() error {
	data, err := c.readPacket()
	if err != nil {
		return err
	}

	if data[0] == mysql.ERR_HEADER {
		return fmt.Errorf("conn [%s] read initial handshake failed, please check [show VARIABLES like 'max_connections'] on mysql end", c.addr)
	}

	if data[0] < mysql.MinProtocolVersion {
		return fmt.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	//skip mysql version and connection id
	//mysql version end with 0x00
	//connection id length is 4
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4

	c.salt = append(c.salt, data[pos:pos+8]...)

	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	c.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos++

		c.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		c.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | c.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		c.salt = append(c.salt, data[pos:pos+12]...)
	}

	return nil
}

func (c *MysqlConn) writeAuthHandshake() error {
	// Adjust client capability flags based on server support
	capability := mysql.CLIENT_PROTOCOL_41 | mysql.CLIENT_SECURE_CONNECTION |
		mysql.CLIENT_LONG_PASSWORD | mysql.CLIENT_TRANSACTIONS | mysql.CLIENT_LONG_FLAG

	capability &= c.capability

	//packet length
	//capbility 4
	//max-packet size 4
	//charset 1
	//reserved all[0] 23
	length := 4 + 4 + 1 + 23

	//username
	length += len(c.user) + 1

	//we only support secure connection
	auth := mysql.CalcPassword(c.salt, []byte(c.password))

	length += 1 + len(auth)

	if len(c.db) > 0 {
		capability |= mysql.CLIENT_CONNECT_WITH_DB

		length += len(c.db) + 1
	}

	c.capability = capability

	data := make([]byte, length+4)

	//capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//MaxPacketSize [32 bit] (none)
	//data[8] = 0x00
	//data[9] = 0x00
	//data[10] = 0x00
	//data[11] = 0x00

	//Charset [1 byte]
	data[12] = byte(c.collation)

	//Filler [23 bytes] (all 0x00)
	pos := 13 + 23

	//User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	//data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		//data[pos] = 0x00
	}

	return c.writePacket(data)
}

func (c *MysqlConn) writeCommand(command byte) error {
	c.pkg.Sequence = 0

	return c.writePacket([]byte{
		0x01, //1 bytes long
		0x00,
		0x00,
		0x00, //sequence
		command,
	})
}

func (c *MysqlConn) writeCommandWithDeadline(command byte) error {
	//add read/write deadline
	c.conn.SetDeadline(time.Now().Add(time.Second))
	c.pkg.Sequence = 0

	return c.writePacket([]byte{
		0x01, //1 bytes long
		0x00,
		0x00,
		0x00, //sequence
		command,
	})
}

func (c *MysqlConn) writeCommandBuf(command byte, arg []byte) error {
	c.pkg.Sequence = 0

	length := len(arg) + 1

	data := make([]byte, length+4)

	data[4] = command

	copy(data[5:], arg)

	return c.writePacket(data)
}

func (c *MysqlConn) writeCommandStr(command byte, arg string) error {
	c.pkg.Sequence = 0

	length := len(arg) + 1

	data := make([]byte, length+4)

	data[4] = command

	copy(data[5:], arg)

	return c.writePacket(data)
}

func (c *MysqlConn) writeCommandStrWithAllocator(alloc *util.FastYetSafeAllocator, command byte, arg string) error {
	c.pkg.Sequence = 0

	length := len(arg) + 1

	//data := make([]byte, length+4)
	data := alloc.AllocWithLen(length+4, length+4)

	data[4] = command

	copy(data[5:], arg)

	return c.writePacket(data)
}

func (c *MysqlConn) writeCommandUint32(command byte, arg uint32) error {
	c.pkg.Sequence = 0

	return c.writePacket([]byte{
		0x05, //5 bytes long
		0x00,
		0x00,
		0x00, //sequence

		command,

		byte(arg),
		byte(arg >> 8),
		byte(arg >> 16),
		byte(arg >> 24),
	})
}

func (c *MysqlConn) writeCommandStrStr(command byte, arg1 string, arg2 string) error {
	c.pkg.Sequence = 0

	data := make([]byte, 4, 6+len(arg1)+len(arg2))

	data = append(data, command)
	data = append(data, arg1...)
	data = append(data, 0)
	data = append(data, arg2...)

	return c.writePacket(data)
}

func (c *MysqlConn) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	if err := c.writeCommandStrStr(mysql.COM_FIELD_LIST, table, wildcard); err != nil {
		return nil, err
	}

	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	fs := make([]*mysql.Field, 0, 4)
	var f *mysql.Field
	if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	}

	for {
		if data, err = c.readPacket(); err != nil {
			return nil, err
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			return fs, nil
		}

		if f, err = mysql.FieldData(data).Parse(); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}

	return nil, fmt.Errorf("field list error")
}

func (c *MysqlConn) readResultset(data []byte, binary bool) (*mysql.Result, error) {
	result := &mysql.Result{
		Status:       0,
		InsertId:     0,
		AffectedRows: 0,

		Resultset: &mysql.Resultset{},
	}

	// column count
	count, _, n := mysql.LengthEncodedInt(data)

	if n-len(data) != 0 {
		return nil, mysql.ErrMalformPacket
	}

	result.Fields = make([]*mysql.Field, count)
	result.FieldNames = make(map[string]int, count)

	if err := c.readResultColumns(result); err != nil {
		return nil, err
	}

	if err := c.readResultRows(result, binary); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *MysqlConn) readResultColumns(result *mysql.Result) (err error) {
	var i int
	var data []byte

	for {
		data, err = c.readPacket()
		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			if i != len(result.Fields) {
				err = mysql.ErrMalformPacket
			}

			return
		}

		result.Fields[i], err = mysql.FieldData(data).Parse()
		if err != nil {
			return
		}

		result.FieldNames[string(result.Fields[i].Name)] = i

		i++
	}
}

func (c *MysqlConn) readResultRows(result *mysql.Result, isBinary bool) (err error) {
	var data []byte

	for {
		data, err = c.readPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		result.RowDatas = append(result.RowDatas, data)
	}

	result.Values = make([][]interface{}, len(result.RowDatas))

	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MysqlConn) readUntilEOF() (err error) {
	var data []byte

	for {
		data, err = c.readPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			return
		}
	}
}

func (c *MysqlConn) isEOFPacket(data []byte) bool {
	return data[0] == mysql.EOF_HEADER && len(data) <= 5
}

func (c *MysqlConn) handleOKPacket(data []byte) (*mysql.Result, error) {
	var n int
	pos := 1

	r := new(mysql.Result)

	r.AffectedRows, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = mysql.LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&mysql.CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2
	}

	//info
	return r, nil
}

func (c *MysqlConn) handleErrorPacket(data []byte) error {
	e := new(mysql.SqlError)

	pos := 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		e.State = string(data[pos : pos+5])
		pos += 5
	}

	e.Message = string(data[pos:])

	return e
}

func (c *MysqlConn) readOK() (*mysql.Result, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == mysql.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		return nil, errors.New("invalid ok packet")
	}
}

func (c *MysqlConn) readResult(binary bool) (*mysql.Result, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == mysql.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == mysql.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else if data[0] == mysql.LocalInFile_HEADER {
		return nil, mysql.ErrMalformPacket
	}

	return c.readResultset(data, binary)
}

func (c *MysqlConn) GetCharset() string {
	return c.charset
}

func (c *MysqlConn) recordTime() {
	c.lastActiveTimestamp = time.Now().Unix()
}
