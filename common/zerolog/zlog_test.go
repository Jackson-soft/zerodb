package zerolog

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestZlogNew(t *testing.T) {
	z, err := NewLog("console", "zlog/log", "info", 0)
	assert.Nil(t, err)
	z.Infoln("this is a message")
	z.Infof("this %d", 45)
	z.WithField("dkey", "dsfs").Infof("ss%d", 445)
	z.WithField("fsdfsd", 3434).Infoln("args ...interface{}")
}

func BenchmarkZlog(t *testing.B) {
	z, err := NewLog("console", "zlog/log", "info", 0)
	assert.Nil(t, err)
	for i := 0; i < t.N; i++ {
		z.Infoln("this is a message")
	}
}

func TestZap(t *testing.T) {
	logerr, err := zap.NewDevelopment()
	assert.Nil(t, err)
	logerr.Sugar().Info("fsdfasdf")
	logerr.Sugar().Infof("dfsdfas %d", 3434)
	logerr.Sugar().With("ddfs", "dfsdf").Info("fsdfasd")
}
