package load

import (
	"os/exec"
	"strconv"
	"strings"
)

func LoadAvg() (*Avg, error) {
	data, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		return nil, err
	}

	v := strings.Replace(string(data), "{ ", "", 1)
	v = strings.Replace(v, " }", "", 1)

	L := strings.Fields(v)
	load1, err := strconv.ParseFloat(L[0], 64)
	if err != nil {
		return nil, err
	}
	load5, err := strconv.ParseFloat(L[1], 64)
	if err != nil {
		return nil, err
	}
	load15, err := strconv.ParseFloat(L[2], 64)
	if err != nil {
		return nil, err
	}

	return &Avg{
		Avg1min:  load1,
		Avg5min:  load5,
		Avg15min: load15,
	}, nil
}
