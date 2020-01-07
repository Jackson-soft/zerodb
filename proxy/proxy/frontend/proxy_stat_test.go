package frontend

import "testing"

func TestGetHostGroupStat(t *testing.T) {
	hostGroupStatMap := testProxyEngine.GetHostGroupStat()
	println(hostGroupStatMap)
}
