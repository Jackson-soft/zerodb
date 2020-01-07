package frontend

import "git.2dfire.net/zerodb/proxy/proxy/mysql"

func (fc *FrontendConn) writeOK(r *mysql.Result) error {
	if r == nil {
		r = &mysql.Result{Status: fc.status}
	}
	//data := make([]byte, 4, 32)
	data := fc.handlerArena.AllocWithLen(4, 32)

	data = append(data, mysql.OK_HEADER)

	data = append(data, mysql.PutLengthEncodedInt(r.AffectedRows)...)
	data = append(data, mysql.PutLengthEncodedInt(r.InsertId)...)

	if fc.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, byte(r.Status), byte(r.Status>>8))
		data = append(data, 0, 0)
	}

	return fc.writePacket(data)
}

func (fc *FrontendConn) writeError(e error) error {
	var m *mysql.SqlError
	var ok bool
	if m, ok = e.(*mysql.SqlError); !ok {
		m = mysql.NewError(mysql.ER_UNKNOWN_ERROR, e.Error())
	}

	//data := make([]byte, 4, 16+len(m.Message))
	data := fc.handlerArena.AllocWithLen(4, 16+len(m.Message))

	data = append(data, mysql.ERR_HEADER)
	data = append(data, byte(m.Code), byte(m.Code>>8))

	if fc.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, '#')
		data = append(data, m.State...)
	}

	data = append(data, m.Message...)

	return fc.writePacket(data)
}

func (fc *FrontendConn) writeEOF(status uint16) error {
	//data := make([]byte, 4, 9)
	data := fc.handlerArena.AllocWithLen(4, 9)

	data = append(data, mysql.EOF_HEADER)
	if fc.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return fc.writePacket(data)
}

func (fc *FrontendConn) writeEOFBatch(total []byte, status uint16, direct bool) ([]byte, error) {
	//data := make([]byte, 4, 9)
	data := fc.handlerArena.AllocWithLen(4, 9)

	data = append(data, mysql.EOF_HEADER)
	if fc.capability&mysql.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, 0, 0)
		data = append(data, byte(status), byte(status>>8))
	}

	return fc.writePacketBatch(total, data, direct)
}
