package frontend

import (
	"fmt"
	"strings"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/backend"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
)

type ExecuteDB struct {
	ExecNode *backend.ShardingNode
	RwSplit  bool
	sql      string
}

//preprocessing sql before parse sql
func (fc *FrontendConn) preHandleShard(sql string) (bool, error) {
	var rs []*mysql.Result
	var err error

	if len(sql) == 0 {
		return false, errcode.ErrSQLNotSupported
	}

	tokens := strings.FieldsFunc(sql, util.IsSqlSep)

	if len(tokens) == 0 {
		return false, errcode.ErrSQLNotSupported
	}

	conn, err := fc.getInfoSingleBackendConn(sql)

	defer fc.closeConn(conn, false)
	if err != nil {
		return false, err
	}
	//execute.sql may be rewritten in getShowExecDB
	rs, err = fc.executeInNode(conn, sql, nil)
	if err != nil {
		return false, err
	}

	if len(rs) == 0 {
		msg := fmt.Sprintf("result is empty")
		return false, mysql.NewError(mysql.ER_UNKNOWN_ERROR, msg)
	}

	fc.lastInsertId = int64(rs[0].InsertId)
	fc.affectedRows = int64(rs[0].AffectedRows)

	if rs[0].Resultset != nil {
		err = fc.writeResultset(fc.status, rs[0].Resultset)
	} else {
		err = fc.writeOK(rs[0])
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (fc *FrontendConn) GetTransExecDB(tokens []string, sql string) (*ExecuteDB, error) {
	var err error
	//tokensLen := len(tokens)
	executeDB := new(ExecuteDB)
	executeDB.sql = sql

	//transaction execute in master db
	executeDB.RwSplit = false

	//if 2 <= tokensLen {
	//    if tokens[0][0] == mysql.COMMENT_PREFIX {
	//        nodeName := strings.Trim(tokens[0], mysql.COMMENT_STRING)
	//        if c.schema.nodes[nodeName] != nil {
	//            executeDB.ExecNode = c.schema.nodes[nodeName]
	//        }
	//    }
	//}

	if executeDB.ExecNode == nil {
		executeDB, err = fc.GetExecDB(tokens, sql)
		if err != nil {
			return nil, err
		}
		if executeDB == nil {
			return nil, nil
		}
		return executeDB, nil
	}
	if len(fc.txConns) == 1 && fc.txConns[executeDB.ExecNode.SchemaIndex] == nil {
		return nil, errcode.ErrMultiNodeTranNotSupport
	}
	return executeDB, nil
}

//if sql need shard return nil, else return the unshard db
func (fc *FrontendConn) GetExecDB(tokens []string, sql string) (*ExecuteDB, error) {
	tokensLen := len(tokens)
	if 0 < tokensLen {
		tokenId, ok := mysql.PARSE_TOKEN_MAP[strings.ToLower(tokens[0])]
		if ok == true {
			switch tokenId {
			case mysql.TK_ID_SELECT:
				return fc.getSelectExecDB(sql, tokens, tokensLen)
				//case mysql.TK_ID_DELETE:
				//    return c.getDeleteExecDB(sql, tokens, tokensLen)
				//case mysql.TK_ID_INSERT, mysql.TK_ID_REPLACE:
				//    return c.getInsertOrReplaceExecDB(sql, tokens, tokensLen)
				//case mysql.TK_ID_UPDATE:
				//    return c.getUpdateExecDB(sql, tokens, tokensLen)
			case mysql.TK_ID_SET:
				return fc.getSetExecDB(sql, tokens, tokensLen)
			case mysql.TK_ID_SHOW:
				return fc.getShowExecDB(sql, tokens, tokensLen)
				//case mysql.TK_ID_TRUNCATE:
				//    return c.getTruncateExecDB(sql, tokens, tokensLen)
			default:
				return nil, nil
			}
		}
	}
	executeDB := new(ExecuteDB)
	executeDB.sql = sql
	err := fc.SetExecuteNodeForDefault(executeDB)
	if err != nil {
		return nil, err
	}
	return executeDB, nil
}

func (fc *FrontendConn) SetExecuteNodeForDefault(executeDB *ExecuteDB) error {
	if executeDB.ExecNode == nil {
		// get an arbitrary value from HostGroup map
		var arbitraryNode *backend.HostGroupNode

		fc.proxy.HostGroupNodes.Range(func(key, value interface{}) bool {
			if v, ok := value.(*backend.HostGroupNode); ok {
				//modified by huajia
				if v.Write[v.GetActivedWriteIndex()].IsActive() {
					arbitraryNode = v
					return false
				}
				return true
				//end modified
			}
			return false
		})


		hostGroupName := arbitraryNode.Cfg.Name
		executeDB.ExecNode = backend.NewShardingNode(0, 0, backend.ShardingTypeSchema, executeDB.sql, hostGroupName)
	}

	return nil
}

//get the execute database for select sql
func (fc *FrontendConn) getSelectExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
	executeDB := new(ExecuteDB)
	executeDB.sql = sql
	executeDB.RwSplit = false

	//if len(rules) != 0 {
	//    for i := 1; i < tokensLen; i++ {
	//        if strings.ToLower(tokens[i]) == mysql.TK_STR_FROM {
	//            if i+1 < tokensLen {
	//                DBName, tableName := sqlparser.GetDBTable(tokens[i+1])
	//                //if the token[i+1] like this:kingshard.test_shard_hash
	//                if DBName != "" {
	//                    ruleDB = DBName
	//                } else {
	//                    ruleDB = c.logicDb
	//                }
	//                if router.GetSchemaRule(ruleDB, tableName) != nil {
	//                    return nil, nil
	//                } else {
	//                    //if the table is not shard table,send the sql
	//                    //to default db
	//                    break
	//                }
	//            }
	//        }
	//
	//        if strings.ToLower(tokens[i]) == mysql.TK_STR_LAST_INSERT_ID {
	//            return nil, nil
	//        }
	//    }
	//}

	//if send to master
	//if 2 < tokensLen {
	//    if strings.ToLower(tokens[1]) == mysql.TK_STR_MASTER_HINT {
	//        executeDB.IsSlave = false
	//    }
	//}
	err := fc.SetExecuteNodeForDefault(executeDB)
	if err != nil {
		return nil, err
	}
	return executeDB, nil
}

//get the execute database for delete sql
//func (c *FrontendConn) getDeleteExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
//    var ruleDB string
//    executeDB := new(ExecuteDB)
//    executeDB.sql = sql
//    router := c.proxy.router
//    rules := router.SchemaRules
//
//    if len(rules) != 0 {
//        for i := 1; i < tokensLen; i++ {
//            if strings.ToLower(tokens[i]) == mysql.TK_STR_FROM {
//                if i+1 < tokensLen {
//                    DBName, tableName := sqlparser.GetDBTable(tokens[i+1])
//                    //if the token[i+1] like this:kingshard.test_shard_hash
//                    if DBName != "" {
//                        ruleDB = DBName
//                    } else {
//                        ruleDB = c.logicDb
//                    }
//                    if router.GetSchemaRule(ruleDB, tableName) != nil {
//                        return nil, nil
//                    } else {
//                        break
//                    }
//                }
//            }
//        }
//    }
//
//    err := c.SetExecuteNodeForDefault(tokens, tokensLen, executeDB)
//    if err != nil {
//        return nil, err
//    }
//
//    return executeDB, nil
//}

//get the execute database for insert or replace sql
//func (c *FrontendConn) getInsertOrReplaceExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
//    var ruleDB string
//    executeDB := new(ExecuteDB)
//    executeDB.sql = sql
//    router := c.proxy.router
//    rules := router.SchemaRules
//
//    if len(rules) != 0 {
//        for i := 0; i < tokensLen; i++ {
//            if strings.ToLower(tokens[i]) == mysql.TK_STR_INTO {
//                if i+1 < tokensLen {
//                    DBName, tableName := sqlparser.GetInsertDBTable(tokens[i+1])
//                    //if the token[i+1] like this:kingshard.test_shard_hash
//                    if DBName != "" {
//                        ruleDB = DBName
//                    } else {
//                        ruleDB = c.logicDb
//                    }
//                    if router.GetSchemaRule(ruleDB, tableName) != nil {
//                        return nil, nil
//                    } else {
//                        break
//                    }
//                }
//            }
//        }
//    }
//
//    err := c.SetExecuteNodeForDefault(tokens, tokensLen, executeDB)
//    if err != nil {
//        return nil, err
//    }
//
//    return executeDB, nil
//}

//get the execute database for update sql
//func (c *FrontendConn) getUpdateExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
//    var ruleDB string
//    executeDB := new(ExecuteDB)
//    executeDB.sql = sql
//    router := c.proxy.router
//    rules := router.SchemaRules
//
//    if len(rules) != 0 {
//        for i := 0; i < tokensLen; i++ {
//            if strings.ToLower(tokens[i]) == mysql.TK_STR_SET {
//                DBName, tableName := sqlparser.GetDBTable(tokens[i-1])
//                //if the token[i+1] like this:kingshard.test_shard_hash
//                if DBName != "" {
//                    ruleDB = DBName
//                } else {
//                    ruleDB = c.logicDb
//                }
//                if router.GetSchemaRule(ruleDB, tableName) != nil {
//                    return nil, nil
//                } else {
//                    break
//                }
//            }
//        }
//    }
//
//    err := c.SetExecuteNodeForDefault(tokens, tokensLen, executeDB)
//    if err != nil {
//        return nil, err
//    }
//
//    return executeDB, nil
//}

//get the execute database for set sql
func (fc *FrontendConn) getSetExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
	executeDB := new(ExecuteDB)
	executeDB.sql = sql

	//handle three styles:
	//set autocommit= 0
	//set autocommit = 0
	//set autocommit=0
	if 2 <= len(tokens) {
		before := strings.Split(sql, "=")
		//uncleanWorld is 'autocommit' or 'autocommit '
		uncleanWord := strings.Split(before[0], " ")
		secondWord := strings.ToLower(uncleanWord[1])
		if _, ok := mysql.SET_KEY_WORDS[secondWord]; ok {
			return nil, nil
		}

		//SET [gobal/session] TRANSACTION ISOLATION LEVEL SERIALIZABLE
		//ignore this sql
		if 3 <= len(uncleanWord) {
			if strings.ToLower(uncleanWord[1]) == mysql.TK_STR_TRANSACTION ||
				strings.ToLower(uncleanWord[2]) == mysql.TK_STR_TRANSACTION {
				return nil, errcode.ErrSQLIgnored
			}
		}
	}

	err := fc.SetExecuteNodeForDefault(executeDB)
	if err != nil {
		return nil, err
	}

	return executeDB, nil
}

//get the execute database for show sql
//choose slave preferentially
//tokens[0] is show
func (fc *FrontendConn) getShowExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
	executeDB := new(ExecuteDB)
	executeDB.RwSplit = false
	executeDB.sql = sql

	//handle show columns/fields
	err := fc.handleShowColumns(sql, tokens, tokensLen, executeDB)
	if err != nil {
		return nil, err
	}

	err = fc.SetExecuteNodeForDefault(executeDB)
	if err != nil {
		return nil, err
	}

	return executeDB, nil
}

//handle show columns/fields
func (fc *FrontendConn) handleShowColumns(sql string, tokens []string,
	tokensLen int, executeDB *ExecuteDB) error {
	/*var ruleDB string
	  for i := 0; i < tokensLen; i++ {
	      tokens[i] = strings.ToLower(tokens[i])
	      //handle SQL:
	      //SHOW [FULL] COLUMNS FROM tbl_name [FROM db_name] [like_or_where]
	      if (strings.ToLower(tokens[i]) == mysql.TK_STR_FIELDS ||
	          strings.ToLower(tokens[i]) == mysql.TK_STR_COLUMNS) &&
	          i+2 < tokensLen {
	          if strings.ToLower(tokens[i+1]) == mysql.TK_STR_FROM {
	              tableName := strings.Trim(tokens[i+2], "`")
	              //get the ruleDB
	              if i+4 < tokensLen && strings.ToLower(tokens[i+1]) == mysql.TK_STR_FROM {
	                  ruleDB = strings.Trim(tokens[i+4], "`")
	              } else {
	                  ruleDB = c.db
	              }
	              showRouter := c.schema.rule
	              showRule := showRouter.GetSchemaRule(ruleDB, tableName)
	              //this SHOW is sharding SQL
	              if showRule.Type != router.DefaultRuleType {
	                  if 0 < len(showRule.FullTableIndexes) {
	                      tableIndex := showRule.FullTableIndexes[0]
	                      nodeIndex := showRule.TableToNode[tableIndex]
	                      nodeName := showRule.Nodes[nodeIndex]
	                      tokens[i+2] = fmt.Sprintf("%s_%04d", tableName, tableIndex)
	                      executeDB.sql = strings.Join(tokens, " ")
	                      executeDB.ExecNode = c.schema.nodes[nodeName]
	                      return nil
	                  }
	              }
	          }
	      }
	  }
	  return nil*/
	return nil
}

//get the execute database for truncate sql
//sql: TRUNCATE [TABLE] tbl_name
//func (c *FrontendConn) getTruncateExecDB(sql string, tokens []string, tokensLen int) (*ExecuteDB, error) {
//    var ruleDB string
//    executeDB := new(ExecuteDB)
//    executeDB.sql = sql
//    router := c.proxy.router
//    rules := router.SchemaRules
//    if len(rules) != 0 && tokensLen >= 2 {
//        DBName, tableName := sqlparser.GetDBTable(tokens[tokensLen-1])
//        //if the token[i+1] like this:kingshard.test_shard_hash
//        if DBName != "" {
//            ruleDB = DBName
//        } else {
//            ruleDB = c.logicDb
//        }
//        if router.GetSchemaRule(ruleDB, tableName) != nil {
//            return nil, nil
//        }
//
//    }
//
//    err := c.SetExecuteNodeForDefault(tokens, tokensLen, executeDB)
//    if err != nil {
//        return nil, err
//    }
//
//    return executeDB, nil
//}
