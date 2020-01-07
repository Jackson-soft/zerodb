package metrics_test

import (
	"testing"

	"git.2dfire.net/zerodb/proxy/metrics"
)

func TestManager(t *testing.T) {
	metrics.Manager.Load("order", 4.5)
}
