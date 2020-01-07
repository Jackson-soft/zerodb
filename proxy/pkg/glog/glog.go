package glog

import (
	"git.2dfire.net/zerodb/common/zerolog"
)

var Glog *zerolog.Zlog

func CreateLogs(logFile, lvl string) error {
	var err error
	Glog, err = zerolog.NewLog(zerolog.CONSOLE, logFile, lvl, 0)
	if err != nil {
		return err
	}

	return nil
}
