package load

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadOSX(t *testing.T) {
	avg, err := LoadAvg()
	assert.Nil(t, err)
	assert.NotNil(t, avg)
	t.Log(avg.Avg1min, avg.Avg5min, avg.Avg15min)
	data, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	assert.Nil(t, err)
	assert.NotNil(t, data)
	t.Logf("%s", data)
}
