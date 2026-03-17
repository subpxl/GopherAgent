package planning

import (
	"fmt"
	"strings"
)

// plan
type Plan struct {
	Goal  string
	Steps []string
}

func (p *Plan) AddStep(step string) {
	p.Steps = append(p.Steps, fmt.Sprintf("step %d: %s", len(p.Steps)+1, step))
}

func (p *Plan) String() string {
	return fmt.Sprintf("Goal: %s\n%s", p.Goal, strings.Join(p.Steps, "\n"))
}
