package load

import (
	"strconv"
	"strings"

	"git.2dfire.net/zerodb/common/system"
)

func LoadAvg() (*Avg, error) {
	data, err := system.ToTrimString("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	L := strings.Fields(data)
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
