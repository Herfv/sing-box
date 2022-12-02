package balancer

import (
	"sync"

	"github.com/sagernet/sing-box/adapter"
)

var _ Strategy = (*FallbackStrategy)(nil)

// FallbackStrategy is the fallback strategy
type FallbackStrategy struct {
	RandomStrategy

	sync.Mutex
	currentTag string
}

// NewFallbackStrategy returns a new FallbackStrategy
func NewFallbackStrategy() *FallbackStrategy {
	return &FallbackStrategy{}
}

// Pick implements Strategy
func (s *FallbackStrategy) Pick(_, filtered []*Node, _ *adapter.InboundContext) *Node {
	if len(filtered) == 0 {
		return nil
	}
	s.Lock()
	defer s.Unlock()
	if s.currentTag != "" {
		for _, node := range filtered {
			if node.Tag() == s.currentTag {
				return node
			}
		}
	}
	// fallback to best node
	s.currentTag = filtered[0].Tag()
	return filtered[0]
}
