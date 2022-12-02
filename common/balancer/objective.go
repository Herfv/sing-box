package balancer

// Objective is the interface for balancer objectives
type Objective interface {
	// Filter filters nodes from all.
	//
	// Conventions：
	//  1. it keeps slice `all` unchanged, because the number and order matters for consistent hash strategy.
	//  2. takes care of the fallback logic, so it will never return a empty slice if `all` is not empty.
	Filter(all []*Node) []*Node
	// Sort sorts nodes according to the objective, better nodes are in front.
	Sort([]*Node)
}
