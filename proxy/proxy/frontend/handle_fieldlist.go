package frontend

import (
	"bytes"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

func (fc *FrontendConn) handleFieldList(data []byte) error {
	index := bytes.IndexByte(data, 0x00)
	table := string(data[0:index])
	wildcard := string(data[index+1:])

	if len(fc.logicDb) == 0 {
		return mysql.NewDefaultError(mysql.ER_NO_DB_ERROR)
	}

	co, err := fc.getInfoSingleBackendConn("")
	defer fc.closeConn(co, false)
	if err != nil {
		return err
	}

	if err = co.UseDB(fc.logicDb); err != nil {
		//reset the database to null
		fc.logicDb = ""
		return err
	}

	if fs, err := co.FieldList(table, wildcard); err != nil {
		return err
	} else {
		return fc.writeFieldList(fc.status, fs)
	}
	return nil
}

func (fc *FrontendConn) writeFieldList(status uint16, fs []*mysql.Field) error {
	fc.affectedRows = int64(-1)
	var err error
	total := make([]byte, 0, 1024)
	data := make([]byte, 4, 512)

	for _, v := range fs {
		data = data[0:4]
		data = append(data, v.Dump()...)
		total, err = fc.writePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	_, err = fc.writeEOFBatch(total, status, true)
	return err
}
