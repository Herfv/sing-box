package balancer_test

import (
	"testing"

	"github.com/sagernet/sing-box/common/balancer"
	"github.com/sagernet/sing-box/common/healthcheck"
)

func TestNodeStatus(t *testing.T) {
	var maxRTT healthcheck.RTT = healthcheck.Second
	var maxFailRate float32 = 0.2
	tests := []struct {
		name   string
		status balancer.Status
		stats  healthcheck.Stats
	}{
		{
			"nil RTTStorage", balancer.StatusUnknown, healthcheck.Stats{
				All: 0, Fail: 0, Latest: 0, Average: 0,
			},
		},
		{
			"untested", balancer.StatusUnknown, healthcheck.Stats{
				All: 0, Fail: 0, Latest: 0, Average: 0,
			},
		},
		{
			"@max_rtt", balancer.StatusQualified, healthcheck.Stats{
				All: 10, Fail: 0, Latest: healthcheck.Second, Average: healthcheck.Second,
			},
		},
		{
			"@max_fail", balancer.StatusQualified, healthcheck.Stats{
				All: 10, Fail: 2, Latest: healthcheck.Second, Average: healthcheck.Second,
			},
		},
		{
			"@max_fail_2", balancer.StatusQualified, healthcheck.Stats{
				All: 5, Fail: 1, Latest: healthcheck.Second, Average: healthcheck.Second,
			},
		},
		{
			"latest_fail", balancer.StatusDead, healthcheck.Stats{
				All: 10, Fail: 1, Latest: healthcheck.Failed, Average: healthcheck.Second,
			},
		},
		{
			"over max_fail", balancer.StatusAlive, healthcheck.Stats{
				All: 5, Fail: 2, Latest: healthcheck.Second, Average: healthcheck.Second,
			},
		},
		{
			"over max_rtt", balancer.StatusAlive, healthcheck.Stats{
				All: 10, Fail: 0, Latest: healthcheck.Second, Average: 2 * healthcheck.Second,
			},
		},
	}
	for i, tt := range tests {
		node := &balancer.Node{Stats: tt.stats}
		node.CalcStatus(maxRTT, maxFailRate)
		if node.Status != tt.status {
			t.Errorf("IsAlive(#%d %s) = (%s), want (%s)", i, tt.name, node.Status, tt.status)
		}
	}
}
