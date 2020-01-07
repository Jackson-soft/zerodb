package frontend

import (
	"git.2dfire.net/zerodb/common/zeroproto/pkg/keeper"
	"github.com/pkg/errors"
	"strconv"
)

type ProxyCluster struct {
	host      string
	port      string
	weight    string
}


//added by huajia,for cobar
func (e *ProxyEngine) SetProxyCluster(data interface{}) error {

	val, _ := data.([]*keeper.ClusterInfo)
	length := len(val)
	if length == 0 {
		return errors.New("get proxy clusters error")
	}

	e.proxyClusters = make([]ProxyCluster, length)
	for i, v := range val {
		e.proxyClusters[i].host = v.Host
		e.proxyClusters[i].port = strconv.FormatUint(v.Port, 10)
		e.proxyClusters[i].weight = strconv.FormatUint(v.Weight, 10)
	}

	return nil
}
//end added