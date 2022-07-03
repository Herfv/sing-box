package route

import (
	"strings"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/domain"
	"github.com/sagernet/sing/common"
)

var _ RuleItem = (*DomainItem)(nil)

type DomainItem struct {
	description string
	matcher     *domain.Matcher
}

func NewDomainItem(domains []string, domainSuffixes []string) *DomainItem {
	domains = common.Uniq(domains)
	domainSuffixes = common.Uniq(domainSuffixes)
	var description string
	if dLen := len(domains); dLen > 0 {
		if dLen == 1 {
			description = "domain=" + domains[0]
		} else if dLen > 3 {
			description = "domain=[" + strings.Join(domains[:3], " ") + "...]"
		} else {
			description = "domain=[" + strings.Join(domains, " ") + "]"
		}
	}
	if dsLen := len(domainSuffixes); dsLen > 0 {
		if len(description) > 0 {
			description += " "
		}
		if dsLen == 1 {
			description += "domainSuffix=" + domainSuffixes[0]
		} else if dsLen > 3 {
			description += "domainSuffix=[" + strings.Join(domainSuffixes[:3], " ") + "...]"
		} else {
			description += "domainSuffix=[" + strings.Join(domainSuffixes, " ") + "]"
		}
	}
	return &DomainItem{
		description,
		domain.NewMatcher(domains, domainSuffixes),
	}
}

func (r *DomainItem) Match(metadata *adapter.InboundContext) bool {
	var domainHost string
	if metadata.Domain != "" {
		domainHost = metadata.Domain
	} else {
		domainHost = metadata.Destination.Fqdn
	}
	if domainHost == "" {
		return false
	}
	return r.matcher.Match(domainHost)
}

func (r *DomainItem) String() string {
	return r.description
}