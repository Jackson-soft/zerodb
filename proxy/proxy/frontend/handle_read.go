package frontend

import (
	"fmt"
	"strconv"
	"strings"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
	"git.2dfire.net/zerodb/proxy/pkg/util"
	"git.2dfire.net/zerodb/proxy/proxy/monitor"
	"git.2dfire.net/zerodb/proxy/proxy/mysql"
	"git.2dfire.net/zerodb/proxy/proxy/sqlparser"
	"time"
)

const (
	MasterComment    = "/*master*/"
	SumFunc          = "sum"
	CountFunc        = "count"
	MaxFunc          = "max"
	MinFunc          = "min"
	LastInsertIdFunc = "last_insert_id"
	FUNC_EXIST       = 1
)

var funcNameMap = map[string]int{
	"sum":            FUNC_EXIST,
	"count":          FUNC_EXIST,
	"max":            FUNC_EXIST,
	"min":            FUNC_EXIST,
	"last_insert_id": FUNC_EXIST,
}

// 处理select语句
// 不支持任何跨库JOIN
// 不支持UNION
func (fc *FrontendConn) handleReadSQL(sql string, stmt *sqlparser.Select, args []interface{}) error {
	if len(fc.logicDb) == 0 {
		return errcode.BuildError(errcode.NoDBUsed)
	}
	sn := fc.proxy.GetSchemaNode(fc.logicDb)
	if sn == nil {
		return errcode.BuildError(errcode.DBNotExist, fc.logicDb)
	}
	startTime := time.Now().UnixNano()
	var rwSplit = sn.RwSplit
	plan, err := fc.proxy.GetRouter().BuildPlan(fc.logicDb, fc.GetRemoteIP(), stmt, sql, fc.permitMultiRoute(fc.status), sn.MultiRoutePermitted, fc.handlerArena.AllocWithLen(0, 256))
	if err != nil {
		return err
	}
	buildPlanTime := time.Now().UnixNano()

	conns, err := fc.getShardingBackendConns(rwSplit, true, plan)
	if err != nil {
		return err
	}
	if conns == nil {
		r := fc.newEmptyResultset(stmt)
		return fc.writeResultset(fc.status, r)
	}
	gotConnsTime := time.Now().UnixNano()

	var rs []*mysql.Result
	rs, err = fc.executeInMultiNodes(conns, plan.ShardingNodes, args, fc.proxy.LogSQL)
	monitor.Monitor.IncrClientQPS(fc.logicDb)

	executeTime := time.Now().UnixNano()

	defer fc.closeShardConns(conns, false)

	if err != nil {
		return err
	}

	err = fc.mergeSelectResult(rs, stmt)
	if err != nil {
		return err
	}

	mergeTime := time.Now().UnixNano()

	logSlowSql(mergeTime, startTime, buildPlanTime, gotConnsTime, executeTime, fc.proxy.slowLogTime, plan.ShardingNodes, sql)

	return err
}

//only merge result with aggregate function in group by opt
func (fc *FrontendConn) mergeGroupByWithFunc(rs []*mysql.Result, groupByIndexs []int,
	funcExprs map[int]string) (*mysql.Result, error) {
	r := rs[0]
	//load rs into a map, in order to make group
	resultMap, err := fc.loadResultWithFuncIntoMap(rs, groupByIndexs, funcExprs)
	if err != nil {
		return nil, err
	}

	//set status
	status := fc.status
	for i := 0; i < len(rs); i++ {
		status = status | rs[i].Status
	}

	//change map into Resultset
	r.Values = nil
	r.RowDatas = nil
	for _, v := range resultMap {
		r.Values = append(r.Values, v.Value)
		r.RowDatas = append(r.RowDatas, v.RowData)
	}
	r.Status = status

	return r, nil
}

//only merge result without aggregate function in group by opt
func (fc *FrontendConn) mergeGroupByWithoutFunc(rs []*mysql.Result,
	groupByIndexs []int) (*mysql.Result, error) {
	r := rs[0]
	//load rs into a map
	resultMap, err := fc.loadResultIntoMap(rs, groupByIndexs)
	if err != nil {
		return nil, err
	}

	//set status
	status := fc.status
	for i := 0; i < len(rs); i++ {
		status = status | rs[i].Status
	}

	//load map into Resultset
	r.Values = nil
	r.RowDatas = nil
	for _, v := range resultMap {
		r.Values = append(r.Values, v.Value)
		r.RowDatas = append(r.RowDatas, v.RowData)
	}
	r.Status = status

	return r, nil
}

type ResultRow struct {
	Value   []interface{}
	RowData mysql.RowData
}

func (fc *FrontendConn) generateMapKey(groupColumns []interface{}) (string, error) {
	bk := make([]byte, 0, 8)
	separatorBuf, err := formatValue("+")
	if err != nil {
		return "", err
	}

	for _, v := range groupColumns {
		b, err := formatValue(v)
		if err != nil {
			return "", err
		}
		bk = append(bk, b...)
		bk = append(bk, separatorBuf...)
	}

	return string(bk), nil
}

func (fc *FrontendConn) loadResultIntoMap(rs []*mysql.Result,
	groupByIndexs []int) (map[string]*ResultRow, error) {
	//load Result into map
	resultMap := make(map[string]*ResultRow)
	for _, r := range rs {
		for i := 0; i < len(r.Values); i++ {
			keySlice := r.Values[i][groupByIndexs[0]:]
			mk, err := fc.generateMapKey(keySlice)
			if err != nil {
				return nil, err
			}

			resultMap[mk] = &ResultRow{
				Value:   r.Values[i],
				RowData: r.RowDatas[i],
			}
		}
	}

	return resultMap, nil
}

func (fc *FrontendConn) loadResultWithFuncIntoMap(rs []*mysql.Result,
	groupByIndexs []int, funcExprs map[int]string) (map[string]*ResultRow, error) {

	resultMap := make(map[string]*ResultRow)
	rt := new(mysql.Result)
	rt.Resultset = new(mysql.Resultset)
	rt.Fields = rs[0].Fields

	//change Result into map
	for _, r := range rs {
		for i := 0; i < len(r.Values); i++ {
			keySlice := r.Values[i][groupByIndexs[0]:]
			mk, err := fc.generateMapKey(keySlice)
			if err != nil {
				return nil, err
			}

			if v, ok := resultMap[mk]; ok {
				//init rt
				rt.Values = nil
				rt.RowDatas = nil

				//append v and r into rt, and calculate the function value
				rt.Values = append(rt.Values, r.Values[i], v.Value)
				rt.RowDatas = append(rt.RowDatas, r.RowDatas[i], v.RowData)
				resultTmp := []*mysql.Result{rt}

				for funcIndex, funcName := range funcExprs {
					funcValue, err := fc.calFuncExprValue(funcName, resultTmp, funcIndex)
					if err != nil {
						return nil, err
					}
					//set the function value in group by
					resultMap[mk].Value[funcIndex] = funcValue
				}
			} else { //key is not exist
				resultMap[mk] = &ResultRow{
					Value:   r.Values[i],
					RowData: r.RowDatas[i],
				}
			}
		}
	}

	return resultMap, nil
}

//build select result without group by opt
func (fc *FrontendConn) buildSelectOnlyResult(rs []*mysql.Result,
	stmt *sqlparser.Select) (*mysql.Result, error) {
	var err error
	r := rs[0].Resultset
	status := fc.status | rs[0].Status

	funcExprs := fc.getFuncExprs(stmt)
	if len(funcExprs) == 0 {
		for i := 1; i < len(rs); i++ {
			status |= rs[i].Status
			for j := range rs[i].Values {
				r.Values = append(r.Values, rs[i].Values[j])
				r.RowDatas = append(r.RowDatas, rs[i].RowDatas[j])
			}
		}
	} else {
		//result only one row, status doesn't need set
		r, err = fc.buildFuncExprResult(stmt, rs, funcExprs)
		if err != nil {
			return nil, err
		}
	}
	return &mysql.Result{
		Status:    status,
		Resultset: r,
	}, nil
}

func (fc *FrontendConn) buildOnlyResult(rs []*mysql.Result,
	stmt *sqlparser.Select) (*mysql.Result, error) {
	r := rs[0].Resultset
	status := fc.status | rs[0].Status

	//r.Fields = append(r.Fields, )

	for i := 1; i < len(rs); i++ {
		status |= rs[i].Status
		for j := range rs[i].Values {
			r.Values = append(r.Values, rs[i].Values[j])
			r.RowDatas = append(r.RowDatas, rs[i].RowDatas[j])
		}
	}

	return &mysql.Result{
		Status:    status,
		Resultset: r,
	}, nil
}

func (fc *FrontendConn) sortSelectResult(r *mysql.Resultset, stmt *sqlparser.Select) error {
	if stmt.OrderBy == nil {
		return nil
	}

	sk := make([]mysql.SortKey, len(stmt.OrderBy))

	for i, o := range stmt.OrderBy {
		sk[i].Name = nstring(o.Expr)
		sk[i].Direction = o.Direction
	}

	return r.Sort(sk)
}

func (fc *FrontendConn) limitSelectResult(r *mysql.Resultset, stmt *sqlparser.Select) error {
	if stmt.Limit == nil {
		return nil
	}

	var offset, count int64
	var err error
	if stmt.Limit.Offset == nil {
		offset = 0
	} else {
		if o, ok := stmt.Limit.Offset.(sqlparser.NumVal); !ok {
			return fmt.Errorf("invalid select limit %s", nstring(stmt.Limit))
		} else {
			if offset, err = strconv.ParseInt(util.String([]byte(o)), 10, 64); err != nil {
				return err
			}
		}
	}

	if o, ok := stmt.Limit.Rowcount.(sqlparser.NumVal); !ok {
		return fmt.Errorf("invalid limit %s", nstring(stmt.Limit))
	} else {
		if count, err = strconv.ParseInt(util.String([]byte(o)), 10, 64); err != nil {
			return err
		} else if count < 0 {
			return fmt.Errorf("invalid limit %s", nstring(stmt.Limit))
		}
	}
	if offset > int64(len(r.Values)) {
		r.Values = nil
		r.RowDatas = nil
		return nil
	}

	if offset+count > int64(len(r.Values)) {
		count = int64(len(r.Values)) - offset
	}

	r.Values = r.Values[offset : offset+count]
	r.RowDatas = r.RowDatas[offset : offset+count]

	return nil
}

func (fc *FrontendConn) buildFuncExprResult(stmt *sqlparser.Select,
	rs []*mysql.Result, funcExprs map[int]string) (*mysql.Resultset, error) {

	var names []string
	var err error
	r := rs[0].Resultset
	funcExprValues := make(map[int]interface{})

	for index, funcName := range funcExprs {
		funcExprValue, err := fc.calFuncExprValue(
			funcName,
			rs,
			index,
		)
		if err != nil {
			return nil, err
		}
		funcExprValues[index] = funcExprValue
	}

	r.Values, err = fc.buildFuncExprValues(rs, funcExprValues)

	if 0 < len(r.Values) {
		for _, field := range rs[0].Fields {
			names = append(names, string(field.Name))
		}
		r, err = fc.buildResultset(rs[0].Fields, names, r.Values)
		if err != nil {
			return nil, err
		}
	} else {
		r = fc.newEmptyResultset(stmt)
	}

	return r, nil
}

//get the index of funcExpr, the value is function name
func (fc *FrontendConn) getFuncExprs(stmt *sqlparser.Select) map[int]string {
	var f *sqlparser.FuncExpr
	funcExprs := make(map[int]string)

	for i, expr := range stmt.SelectExprs {
		nonStarExpr, ok := expr.(*sqlparser.NonStarExpr)
		if !ok {
			continue
		}

		f, ok = nonStarExpr.Expr.(*sqlparser.FuncExpr)
		if !ok {
			continue
		} else {
			f = nonStarExpr.Expr.(*sqlparser.FuncExpr)
			funcName := strings.ToLower(string(f.Name))
			switch funcNameMap[funcName] {
			case FUNC_EXIST:
				funcExprs[i] = funcName
			}
		}
	}
	return funcExprs
}

func (fc *FrontendConn) getSumFuncExprValue(rs []*mysql.Result,
	index int) (interface{}, error) {
	var sumf float64
	var sumi int64
	var IsInt bool
	var err error
	var result interface{}

	for _, r := range rs {
		for k := range r.Values {
			result, err = r.GetValue(k, index)
			if err != nil {
				return nil, err
			}
			if result == nil {
				continue
			}

			switch v := result.(type) {
			case int:
				sumi = sumi + int64(v)
				IsInt = true
			case int32:
				sumi = sumi + int64(v)
				IsInt = true
			case int64:
				sumi = sumi + v
				IsInt = true
			case float32:
				sumf = sumf + float64(v)
			case float64:
				sumf = sumf + v
			case []byte:
				tmp, err := strconv.ParseFloat(string(v), 64)
				if err != nil {
					return nil, err
				}

				sumf = sumf + tmp
			default:
				return nil, errcode.ErrSumColumnType
			}
		}
	}
	if IsInt {
		return sumi, nil
	} else {
		return sumf, nil
	}
}

func (fc *FrontendConn) getMaxFuncExprValue(rs []*mysql.Result,
	index int) (interface{}, error) {
	var max interface{}
	var findNotNull bool
	if len(rs) == 0 {
		return nil, nil
	} else {
		for _, r := range rs {
			for k := range r.Values {
				result, err := r.GetValue(k, index)
				if err != nil {
					return nil, err
				}
				if result != nil {
					max = result
					findNotNull = true
					break
				}
			}
			if findNotNull {
				break
			}
		}
	}
	for _, r := range rs {
		for k := range r.Values {
			result, err := r.GetValue(k, index)
			if err != nil {
				return nil, err
			}
			if result == nil {
				continue
			}
			switch result.(type) {
			case int64:
				if max.(int64) < result.(int64) {
					max = result
				}
			case uint64:
				if max.(uint64) < result.(uint64) {
					max = result
				}
			case float64:
				if max.(float64) < result.(float64) {
					max = result
				}
			case string:
				if max.(string) < result.(string) {
					max = result
				}
			}
		}
	}
	return max, nil
}

func (fc *FrontendConn) getMinFuncExprValue(
	rs []*mysql.Result, index int) (interface{}, error) {
	var min interface{}
	var findNotNull bool
	if len(rs) == 0 {
		return nil, nil
	} else {
		for _, r := range rs {
			for k := range r.Values {
				result, err := r.GetValue(k, index)
				if err != nil {
					return nil, err
				}
				if result != nil {
					min = result
					findNotNull = true
					break
				}
			}
			if findNotNull {
				break
			}
		}
	}
	for _, r := range rs {
		for k := range r.Values {
			result, err := r.GetValue(k, index)
			if err != nil {
				return nil, err
			}
			if result == nil {
				continue
			}
			switch result.(type) {
			case int64:
				if min.(int64) > result.(int64) {
					min = result
				}
			case uint64:
				if min.(uint64) > result.(uint64) {
					min = result
				}
			case float64:
				if min.(float64) > result.(float64) {
					min = result
				}
			case string:
				if min.(string) > result.(string) {
					min = result
				}
			}
		}
	}
	return min, nil
}

//calculate the the value funcExpr(sum or count)
func (fc *FrontendConn) calFuncExprValue(funcName string,
	rs []*mysql.Result, index int) (interface{}, error) {

	var num int64
	switch strings.ToLower(funcName) {
	case CountFunc:
		if len(rs) == 0 {
			return nil, nil
		}
		for _, r := range rs {
			if r != nil {
				for k := range r.Values {
					result, err := r.GetInt(k, index)
					if err != nil {
						return nil, err
					}
					num += result
				}
			}
		}
		return num, nil
	case SumFunc:
		return fc.getSumFuncExprValue(rs, index)
	case MaxFunc:
		return fc.getMaxFuncExprValue(rs, index)
	case MinFunc:
		return fc.getMinFuncExprValue(rs, index)
	case LastInsertIdFunc:
		return fc.lastInsertId, nil
	default:
		if len(rs) == 0 {
			return nil, nil
		}
		//get a non-null value of funcExpr
		for _, r := range rs {
			for k := range r.Values {
				result, err := r.GetValue(k, index)
				if err != nil {
					return nil, err
				}
				if result != nil {
					return result, nil
				}
			}
		}
	}

	return nil, nil
}

//build values of resultset,only build one row
func (fc *FrontendConn) buildFuncExprValues(rs []*mysql.Result,
	funcExprValues map[int]interface{}) ([][]interface{}, error) {
	values := make([][]interface{}, 0, 1)
	//build a row in one result
	for i := range rs {
		for j := range rs[i].Values {
			for k := range funcExprValues {
				rs[i].Values[j][k] = funcExprValues[k]
			}
			values = append(values, rs[i].Values[j])
			if len(values) == 1 {
				break
			}
		}
		break
	}

	//generate one row just for sum or count
	if len(values) == 0 {
		value := make([]interface{}, len(rs[0].Fields))
		for k := range funcExprValues {
			value[k] = funcExprValues[k]
		}
		values = append(values, value)
	}

	return values, nil
}
