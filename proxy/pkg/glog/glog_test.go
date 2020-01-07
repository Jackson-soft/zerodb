package glog

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlog(t *testing.T) {
	u := "rpc-keeper-middleware-daily.app.2dfire-daily.com"
	conn, err := net.Dial("tcp", u)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}
