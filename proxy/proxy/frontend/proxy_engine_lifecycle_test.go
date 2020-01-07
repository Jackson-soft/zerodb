package frontend

import "testing"

func TestCloseProxyEngine(t *testing.T) {
	testProxyEngine.CloseProxyEngine()

	started := make(chan int)
	go testProxyEngine.ServeInBlocking(func() {
		started <- 1
	})
	<-started
}
