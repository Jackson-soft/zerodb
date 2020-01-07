package utils

import (
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"git.2dfire.net/zerodb/common/statics"

	"errors"

	"gopkg.in/yaml.v3"

	jsoniter "github.com/json-iterator/go"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	LoadYaml   = yaml.Unmarshal
	UnLoadYaml = yaml.Marshal
	LoadJSON   = json.Unmarshal
	UnLoadJSON = json.Marshal

	parameterIsNil = errors.New("parameter is nil.")
)

//GetLocalIP return loacl ip
func GetLocalIP() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

//IsLoopback check ip is 127.0.0.1 or 0.0.0.0 and host,port
func IsLoopback(addr string) (bool, string, string) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return true, host, port
	}
	if len(host) == 0 {
		return true, host, port
	}
	ip := net.ParseIP(host)
	if ip.IsLoopback() || ip.Equal(net.IPv4zero) {
		return true, host, port
	}
	return false, host, port
}

//WaitTimeout 带超时的waitgroup
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

//检查一组ip地址是否是有效的ipv4或者ipv6地址,如果有任何一个不是返回false
func ValidIP(ips []string) bool {
	if len(ips) == 0 {
		return false
	}
	for _, ip := range ips {
		i := strings.Split(ip, ":")
		if net.ParseIP(i[0]) == nil {
			return false
		}
	}
	return true
}

//check slice contains string
func Contains(s []string, str string) bool {
	for _, a := range s {
		if a == str {
			return true
		}
	}
	return false
}

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

//ParseTime 校验时间的正确性
func ParseTime(strTime string) (time.Time, error) {
	if len(strTime) == 0 {
		return time.Time{}, parameterIsNil
	}
	return time.ParseInLocation(statics.LayoutTime, strTime, time.Local)
}
