package assistant

import "strings"

type RetrievalPolicy struct {
	MinQueryChars int
	Limit         int
}

func (p RetrievalPolicy) ShouldRetrieve(query string) bool {
	query = strings.TrimSpace(query)
	min := p.MinQueryChars
	if min <= 0 {
		min = 8
	}
	return len([]rune(query)) >= min
}

func (p RetrievalPolicy) TopK() int {
	if p.Limit <= 0 {
		return 5
	}
	return p.Limit
}
