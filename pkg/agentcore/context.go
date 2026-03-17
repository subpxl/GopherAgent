package agentcore

import (
	"fmt"
	"strings"
)

// context

type AgentContext struct {
	History []string
}

func (c *AgentContext) Append(role, msg string) {
	c.History = append(c.History, fmt.Sprintf("[%s] %s", role, msg))

}

func (c *AgentContext) String() string {
	return strings.Join(c.History, "\n")
}
