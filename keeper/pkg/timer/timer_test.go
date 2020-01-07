package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	tt := NewTimer(time.Second, 5)
	id, err := tt.Register(Single, 3*time.Second, func(args interface{}) {
		t.Log(args)
	}, "fsdfsdfsa")
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
}

func TestMy(t *testing.T) {
	tt := "2018-09-28 18:59:12"
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
	assert.Nil(t, err)
	t.Log(t1)
	bb := time.Now().After(t1)
	t.Log(bb)

	t.Log(t1.Sub(time.Now()))
}
