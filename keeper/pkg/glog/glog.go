package glog

import (
	"git.2dfire.net/zerodb/common/zerolog"
)

var GLog *zerolog.Zlog

func CreateLog(file, lvl string) error {
	var err error
	GLog, err = zerolog.NewLog(zerolog.CONSOLE, file, lvl, 0)
	if err != nil {
		return err
	}
	return nil
}
