package frontend

import (
	"fmt"
	"strconv"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
)

func formatValue(value interface{}) ([]byte, error) {
	if value == nil {
		return util.Slice("NULL"), nil
	}
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return util.Slice(v), nil
	default:
		return nil, fmt.Errorf("invalid type %T", value)
	}
}

func (fc *FrontendConn) mergeSelectResult(rs []*mysql.Result, stmt *sqlparser.Select) error {
	var r *mysql.Result
	var err error

	if len(stmt.GroupBy) == 0 {
		r, err = fc.buildSelectOnlyResult(rs, stmt)
	} else {
		//group by
		r, err = fc.buildSelectGroupByResult(rs, stmt)
	}
	if err != nil {
		return err
	}

	fc.sortSelectResult(r.Resultset, stmt)
	//to do, add log here, sort may error because order by key not exist in resultset fields

	if err := fc.limitSelectResult(r.Resultset, stmt); err != nil {
		return err
	}

	return fc.writeResultset(r.Status, r.Resultset)
}

func (fc *FrontendConn) mergeResults(rs []*mysql.Result, stmt *sqlparser.Select) (*mysql.Result, error) {
	var r *mysql.Result
	var err error

	r, err = fc.buildOnlyResult(rs, stmt)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func formatField(field *mysql.Field, value interface{}) error {
	switch value.(type) {
	case int8, int16, int32, int64, int:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_LONGLONG
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG
	case uint8, uint16, uint32, uint64, uint:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_LONGLONG
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG | mysql.UNSIGNED_FLAG
	case float32, float64:
		field.Charset = 63
		field.Type = mysql.MYSQL_TYPE_DOUBLE
		field.Flag = mysql.BINARY_FLAG | mysql.NOT_NULL_FLAG
	case string, []byte:
		field.Charset = 33
		field.Type = mysql.MYSQL_TYPE_VAR_STRING
	default:
		return fmt.Errorf("unsupport type %T for resultset", value)
	}
	return nil
}

func (fc *FrontendConn) buildResultset(fields []*mysql.Field, names []string, values [][]interface{}) (*mysql.Resultset, error) {
	var ExistFields bool
	r := new(mysql.Resultset)

	r.Fields = make([]*mysql.Field, len(names))
	r.FieldNames = make(map[string]int, len(names))

	//use the field def that get from true database
	if len(fields) != 0 {
		if len(r.Fields) == len(fields) {
			ExistFields = true
		} else {
			return nil, errcode.ErrInvalidArgument
		}
	}

	var b []byte
	var err error

	for i, vs := range values {
		if len(vs) != len(r.Fields) {
			return nil, fmt.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		for j, value := range vs {
			//列的定义
			if i == 0 {
				if ExistFields {
					r.Fields[j] = fields[j]
					r.FieldNames[string(r.Fields[j].Name)] = j
				} else {
					field := &mysql.Field{}
					r.Fields[j] = field
					r.FieldNames[string(r.Fields[j].Name)] = j
					field.Name = util.Slice(names[j])
					if err = formatField(field, value); err != nil {
						return nil, err
					}
				}

			}
			b, err = formatValue(value)
			if err != nil {
				return nil, err
			}

			row = append(row, mysql.PutLengthEncodedString(b)...)
		}

		r.RowDatas = append(r.RowDatas, row)
	}
	//assign the values to the result
	r.Values = values

	return r, nil
}

func (fc *FrontendConn) mergeExecResult(rs []*mysql.Result) error {
	r := new(mysql.Result)
	for _, v := range rs {
		r.Status |= v.Status
		r.AffectedRows += v.AffectedRows
		if r.InsertId == 0 {
			r.InsertId = v.InsertId
		} else if r.InsertId > v.InsertId {
			//last insert id is first gen id for multi row inserted
			//see http://dev.mysql.com/doc/refman/5.6/en/information-functions.html#function_last-insert-id
			r.InsertId = v.InsertId
		}
	}

	if r.InsertId > 0 {
		fc.lastInsertId = int64(r.InsertId)
	}
	fc.affectedRows = int64(r.AffectedRows)

	return fc.writeOK(r)
}

func (fc *FrontendConn) writeResultset(status uint16, r *mysql.Resultset) error {
	fc.affectedRows = int64(-1)
	//total := make([]byte, 0, 4096)
	total := fc.handlerArena.AllocWithLen(0, 4096)
	//data := make([]byte, 4, 512)
	data := fc.handlerArena.AllocWithLen(4, 512)

	var err error

	columnLen := mysql.PutLengthEncodedInt(uint64(len(r.Fields)))

	data = append(data, columnLen...)
	total, err = fc.writePacketBatch(total, data, false)
	if err != nil {
		return err
	}

	for _, v := range r.Fields {
		data = data[0:4]
		data = append(data, v.Dump()...)
		total, err = fc.writePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	total, err = fc.writeEOFBatch(total, status, false)
	if err != nil {
		return err
	}

	for _, v := range r.RowDatas {
		data = data[0:4]
		data = append(data, v...)
		total, err = fc.writePacketBatch(total, data, false)
		if err != nil {
			return err
		}
	}

	total, err = fc.writeEOFBatch(total, status, true)
	total = nil
	if err != nil {
		return err
	}

	return nil
}

func (fc *FrontendConn) newEmptyResultset(stmt *sqlparser.Select) *mysql.Resultset {
	r := new(mysql.Resultset)
	r.Fields = make([]*mysql.Field, len(stmt.SelectExprs))

	for i, expr := range stmt.SelectExprs {
		r.Fields[i] = &mysql.Field{}
		switch e := expr.(type) {
		case *sqlparser.StarExpr:
			r.Fields[i].Name = []byte("*")
		case *sqlparser.NonStarExpr:
			if e.As != nil {
				r.Fields[i].Name = e.As
				r.Fields[i].OrgName = util.Slice(nstring(e.Expr))
			} else {
				r.Fields[i].Name = util.Slice(nstring(e.Expr))
			}
		default:
			r.Fields[i].Name = util.Slice(nstring(e))
		}
	}

	r.Values = make([][]interface{}, 0)
	r.RowDatas = make([]mysql.RowData, 0)

	return r
}
