package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLoop(t *testing.T) {
	test := "127.0.0.1:8080"
	want := true
	want2 := "8080"

	v, _, v2 := IsLoopback(test)
	if v != want || v2 != want2 {
		t.Error("test1")
	}

	test = "0.0.0.0:8080"
	v, _, v2 = IsLoopback(test)

	if v != want || v2 != want2 {
		t.Error("test2")
	}

	test = ":8080"
	v, _, v2 = IsLoopback(test)
	if v != want || v2 != want2 {
		t.Error("test3")
	}
}

func Test_validIP(t *testing.T) {
	type args struct {
		ips []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case1",
			args: args{ips: []string{"10.0.10.1", "192.168.2.1"}},
			want: true,
		},
		{
			name: "case2",
			args: args{ips: []string{"10.1.10.1,'hello", "300.178.123.11"}},
			want: false,
		},
		{
			name: "case3",
			args: args{ips: []string{""}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidIP(tt.args.ips); got != tt.want {
				t.Errorf("validIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	var a Accuracy = func() float64 {
		return 0.00001
	}
	t.Log(a.Equal(0.11111222, 0.11111222233333))
	t.Log(a.Greater(0.01, 0.00))
	t.Log(0.01 > 0)
}

func TestCheckTime(t *testing.T) {
	case1 := "2013-02-03 19:54:00"
	tt, err := ParseTime(case1)
	assert.Nil(t, err)
	t.Log(tt)
}
