package system

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Mem struct {
	Buffers   uint64
	Cached    uint64
	MemTotal  uint64
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

func (this *Mem) String() string {
	return fmt.Sprintf("<MemTotal:%d, MemFree:%d, Buffers:%d, Cached:%d...>", this.MemTotal, this.MemFree, this.Buffers, this.Cached)
}

var Multi uint64 = 1024

var WANT = map[string]struct{}{
	"Buffers:":   struct{}{},
	"Cached:":    struct{}{},
	"MemTotal:":  struct{}{},
	"MemFree:":   struct{}{},
	"SwapTotal:": struct{}{},
	"SwapFree:":  struct{}{},
}

func MemInfo() (*Mem, error) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	memInfo := &Mem{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, err := ReadLine(reader)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(string(line))
		fieldName := fields[0]

		_, ok := WANT[fieldName]
		if ok && len(fields) == 3 {
			val, numerr := strconv.ParseUint(fields[1], 10, 64)
			if numerr != nil {
				continue
			}
			switch fieldName {
			case "Buffers:":
				memInfo.Buffers = val * Multi
			case "Cached:":
				memInfo.Cached = val * Multi
			case "MemTotal:":
				memInfo.MemTotal = val * Multi
			case "MemFree:":
				memInfo.MemFree = val * Multi
			case "SwapTotal:":
				memInfo.SwapTotal = val * Multi
			case "SwapFree:":
				memInfo.SwapFree = val * Multi
			}
		}
	}

	memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree

	return memInfo, nil
}

func MemLoad() (float64, error) {
	m, err := MemInfo()
	if err != nil {
		return 0, err
	}
	mf := m.MemTotal - m.MemFree - m.Buffers - m.Cached
	l := float64(mf) / float64(m.MemTotal)
	return round(l, 4), nil
}
