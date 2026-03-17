package memory

import (
	"strings"
)

// memory

type Memory struct {
	ShortTerm []string
	LongTerm  []string
	MaxItems  int
}

func NewMemory(maxItems int) *Memory {
	return &Memory{MaxItems: maxItems}
}

func (m *Memory) Remember(fact string) {
	m.ShortTerm = append(m.ShortTerm, fact)
	if len(m.ShortTerm) > m.MaxItems {
		m.ShortTerm = m.ShortTerm[1:]
	}
}

func (m *Memory) Commit(fact string) {
	m.LongTerm = append(m.LongTerm, fact)
}

func (m *Memory) Recall() string {
	all := append(m.LongTerm, m.ShortTerm...)
	return strings.Join(all, "\n")
}
