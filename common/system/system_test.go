package system

import (
	"os/exec"
	"testing"
)

func TestCpuLoad(t *testing.T) {
	s, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", s)
}
