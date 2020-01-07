package frontend

import "git.2dfire.net/zerodb/proxy/proxy/backend"

type HostGroupInfo struct {
	activedWriteIndex int
	remainIdleSize    int
}

func (e *ProxyEngine) GetHostGroupStat() map[string]*HostGroupInfo {
	hostGroupStat := make(map[string]*HostGroupInfo)

	e.HostGroupNodes.Range(func(key, value interface{}) bool {
		if v, ok := value.(*backend.HostGroupNode); ok {
			if w := v.Write[v.GetActivedWriteIndex()]; w != nil {
				hostGroupStat[v.Cfg.Name] = &HostGroupInfo{
					activedWriteIndex: v.GetActivedWriteIndex(),
					remainIdleSize:    w.RemainIdleSize(),
				}
			}
		}
		return true
	})
	return hostGroupStat
}
