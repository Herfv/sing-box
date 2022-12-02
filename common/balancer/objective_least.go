package balancer

import (
	"sort"

	"github.com/sagernet/sing-box/common/healthcheck"
	"github.com/sagernet/sing-box/option"
)

var _ Objective = (*LeastObjective)(nil)

// LeastObjective is the least load / least ping balancing objective
type LeastObjective struct {
	*QualifiedObjective
	expected  uint
	baselines []healthcheck.RTT
	rttFunc   func(node *Node) healthcheck.RTT
}

// NewLeastObjective returns a new LeastObjective
func NewLeastObjective(sampling uint, options option.LoadBalancePickOptions, rttFunc func(node *Node) healthcheck.RTT) *LeastObjective {
	return &LeastObjective{
		QualifiedObjective: NewQualifiedObjective(),
		expected:           options.Expected,
		baselines:          healthcheck.RTTsOf(options.Baselines),
		rttFunc:            rttFunc,
	}
}

// Filter implements Objective.
// NOTICE: be aware of the coding convention of this function
func (o *LeastObjective) Filter(all []*Node) []*Node {
	// nodes are either qualified, alive or all nodes
	nodes := o.QualifiedObjective.Filter(all)
	o.Sort(nodes)
	// LeastNodes will always select at least one node
	return LeastNodes(
		nodes,
		o.expected, o.baselines,
		o.rttFunc,
	)
}

// Sort implements Objective.
func (o *LeastObjective) Sort(all []*Node) {
	SortByLeast(all, o.rttFunc)
}

// LeastNodes filters ordered nodes according to Baselines and Expected Count.
//
// The strategy always improves network response speed, not matter which mode below is configurated.
// But they can still have different priorities.
//
// 1. Bandwidth priority: no Baseline + Expected Count > 0.: selects `Expected Count` of nodes.
// (one if Expected Count <= 0)
//
// 2. Bandwidth priority advanced: Baselines + Expected Count > 0.
// Select `Expected Count` amount of nodes, and also those near them according to baselines.
// In other words, it selects according to different Baselines, until one of them matches
// the Expected Count, if no Baseline matches, Expected Count applied.
//
// 3. Speed priority: Baselines + `Expected Count <= 0`.
// go through all baselines until find selects, if not, select none. Used in combination
// with 'balancer.fallbackTag', it means: selects qualified nodes or use the fallback.
func LeastNodes(
	nodes []*Node, expected uint, baselines []healthcheck.RTT,
	rttFunc func(node *Node) healthcheck.RTT,
) []*Node {
	if len(nodes) == 0 {
		// s.logger.Debug("no qualified nodes")
		return nil
	}
	expected2 := int(expected)
	availableCount := len(nodes)
	if expected2 > availableCount {
		return nodes
	}

	if expected2 <= 0 {
		expected2 = 1
	}
	if len(baselines) == 0 {
		return nodes[:expected2]
	}

	count := 0
	// go through all base line until find expected selects
	for _, baseline := range baselines {
		for i := count; i < availableCount; i++ {
			if rttFunc(nodes[i]) >= baseline {
				break
			}
			count = i + 1
		}
		// don't continue if find expected selects
		if count >= expected2 {
			break
		}
	}
	if expected > 0 && count < expected2 {
		count = expected2
	}
	return nodes[:count]
}

// SortByLeast sorts nodes by least value from rttFunc and more.
func SortByLeast(nodes []*Node, rttFunc func(*Node) healthcheck.RTT) {
	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]
		if left.Status != right.Status {
			return left.Status > right.Status
		}
		leftRTT, rightRTT := rttFunc(left), rttFunc(right)
		if leftRTT != rightRTT {
			// 0, 100
			if leftRTT == healthcheck.Failed {
				return false
			}
			// 100, 0
			if rightRTT == healthcheck.Failed {
				return true
			}
			// 100, 200
			return leftRTT < rightRTT
		}
		if left.Fail != right.Fail {
			return left.Fail < right.Fail
		}
		return left.All > right.All
	})
}
