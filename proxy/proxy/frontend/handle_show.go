package frontend

import (
	"fmt"
	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
)

func (fc *FrontendConn) handleShow(s *sqlparser.Show) error {
	var Column = 1
	var rows [][]string
	var names []string
	var values [][]interface{}

	if s.Type == sqlparser.AST_SHOW_TYPE_DATABASES {
		names = []string{"DATABASES"}

		fc.proxy.schemaNodes.Range(func(key, value interface{}) bool {
			if k, ok := key.(string); ok {
				rows = append(rows, []string{k})
			}

			return true
		})

		values = make([][]interface{}, len(rows))
		for i := range rows {
			values[i] = make([]interface{}, Column)
			for j := range rows[i] {
				values[i][j] = rows[i][j]
			}
		}
	} else if s.Type == sqlparser.AST_SHOW_TYPE_TABLES {
		// TODO nanxing need include nonsharding tables later on
		names = []string{"TABLES"}

		if len(fc.logicDb) == 0 {
			return errcode.BuildError(errcode.NoDBUsed)
		}

		fc.proxy.schemaNodes.Range(func(key, value interface{}) bool {
			if v, ok := value.(*backend.SchemaNode); ok {
				if v.Name == fc.logicDb {

					v.TableNodes.Range(func(key, value interface{}) bool {
						if k, ok := key.(string); ok {
							rows = append(rows, []string{k})
						}
						return true
					})
				}
			}

			return true
		})

		values = make([][]interface{}, len(rows))
		for i := range rows {
			values[i] = make([]interface{}, Column)
			for j := range rows[i] {
				values[i][j] = rows[i][j]
			}
		}
	} else if s.Type == sqlparser.AST_SHOW_TYPE_COBAR_CLUSTER {
		names = []string{"HOST", "WEIGHT"}

		for _, v := range(fc.proxy.proxyClusters) {
			rows = append(rows, []string{v.host, v.weight})
		}

		values = make([][]interface{}, len(rows))
		for i := range rows {
			values[i] = make([]interface{}, len(names))
			for j := range rows[i] {
				values[i][j] = rows[i][j]
			}
		}

	}

	if len(rows) == 0 {
		return fmt.Errorf("no data for show statement")
	}

	resultset, err := fc.buildResultset(nil, names, values)
	if err != nil {
		return err
	}

	return fc.writeResultset(fc.status, resultset)
}
