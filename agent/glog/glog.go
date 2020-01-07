package glog

import (
	"errors"

	"git.2dfire.net/zerodb/common/zerolog"
)

var Glog *zerolog.Zlog

//CreateLogs 初始化各文件日志
func CreateLogs(logPath, lvl string) error {
	if logPath == "" || lvl == "" {
		return errors.New("path or level is nil")
	}

	var err error
	Glog, err = zerolog.NewLog(zerolog.CONSOLE, logPath, lvl, 0)
	if err != nil {
		return err
	}

	return nil
}
