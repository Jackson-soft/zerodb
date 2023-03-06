//codes package 定义了keeper所有可能返回的error codes,参考了grpc-go的code定义,所有由keeper产生的错误由2000-3000来表示

package codes

type Code = uint32

const (
	Failed Code = iota
	Succeed
)

const (
	//OK 成功
	OK Code = 0

	//NetConnectFailed 网络连接失败
	NetConnectFailed Code = 1

	//ParameterIsNil api参数为空
	ParameterIsNil Code = 2

	//YamlParseFail yaml文件解析失败
	YamlParseFail = 3

	//DataReadFail 数据读取失败
	DataReadFail = 4

	//CommonError  一 个错误
	CommonError = 5

	//ClusterStatusError 代理集群状态不对
	ClusterStatusError = 6

	//ConfigCheckFailed 配置检测失败
	ConfigCheckFailed = 7

	//FailedConnectEtcd 表示无法连接到Etcd集群
	FailedConnectEtcd Code = 2001

	//FailedCreateLockSession 表示无法创建获取host锁的session
	FailedCreateLockSession Code = 2002

	//FailedGetLock 表示无法获取host锁
	FailedGetLock Code = 2003

	//LockAllReadyHold 表示锁已经被其他Client获得
	LockAllReadyHold Code = 2004

	//EtcdGetFailed etcd get失败
	EtcdGetFailed Code = 2005

	//EtcdPutFailed etcd put失败
	EtcdPutFailed Code = 2006

	//FailedGetVote 表示keeper在获取投票的过程中,zeroproxy返回的basicResp为nil.
	FailedGetVote Code = 2007

	//SystemLoadHigh Mysql负载高
	SystemLoadHigh Code = 2008

	//binlog 不同步
	BinlogNotConsistent Code = 2019

	//SystemLoadFailed 获取负载失败
	SystemLoadFailed Code = 2009

	//MysqlQueryFailed mysql查询错误
	MysqlQueryFailed Code = 2010

	//FailedOperateEtcd 表示操作keeper操作etcd出错
	FailedOperateEtcd Code = 2011

	//LessHalfVote 投票未过半
	LessHalfVote Code = 2012

	//FailedSetDB 表示setDB出错
	FailedSetDB Code = 2013

	//SwitchDBFrequently 表示切换太频繁
	SwitchDBFrequently Code = 2014

	//FailedUnLock 解锁失败
	FailedUnLock Code = 2015

	//EtcdKeyNoExist 表示获取etcd key时候key不存在
	EtcdKeyNoExist Code = 2016
	//InvalidSwitch 表示一个非法的切换,例如要切换到的数据库不存在
	InvalidSwitch Code = 2020
)

var codeText = map[Code]string{
	OK:                 "Ok",
	FailedGetVote:      "Failed Get Proxy Vote",
	SwitchDBFrequently: "SwitchDB too Frequently",
	InvalidSwitch:      "InvalidSwitch",
	LessHalfVote:       "less than need vote",
	LockAllReadyHold:   "lock allready hold by other proxy",
	SystemLoadHigh:     "SystemLoadHigh",
}

func CodeText(code Code) string {
	return codeText[code]
}
