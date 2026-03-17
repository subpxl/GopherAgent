package agentcore

import (
	"agent-in-go/pkg/llm"
	"agent-in-go/pkg/memory"
	"agent-in-go/pkg/planning"
	"agent-in-go/pkg/skills"
	"agent-in-go/pkg/tools"
	"fmt"
	"strings"
)

// Agent
type Agent struct {
	Model       string
	MaxSteps    int
	Personality *Personality
	Plan        planning.Plan
	Context     AgentContext
	Tools       tools.Tools
	Skills      skills.Skills
	Memory      *memory.Memory
}

func NewAgent(model string, maxSteps int, personality *Personality) *Agent {

	return &Agent{
		Personality: personality,
		Model:       model,
		MaxSteps:    maxSteps,
		Tools:       make(tools.Tools),
		Skills:      make(skills.Skills),
		Memory:      memory.NewMemory(10),
	}

}

func (a *Agent) systemPrompt() string {
	return fmt.Sprintf(`You are an autonomous agent. Think step by step, then act.

RESPONSE FORMAT — follow exactly, one line per response, no exceptions:

  To call a tool:   TOOL:<tool_name>|<input>
  To give an answer: FINAL:<answer>

RULES:
- Output ONE line only — either a TOOL: call or a FINAL: answer
- Never output "Final answer:" or any other prefix — only TOOL: or FINAL:
- After a tool returns a result, use the ACTUAL value — never say "result" or "the output"
- If no tool is needed, respond immediately with FINAL:<answer>
- FINAL answers may contain spaces, punctuation, full sentences — everything after FINAL: is the answer


YOUR PERSONALITY:
%s


AVAILABLE TOOLS:
%s

AVAILABLE SKILLS:
%s

MEMORY:
%s

CURRENT PLAN:
%s

EXAMPLES:
  User: what is 10 + 5?
  Agent: TOOL:calculator|10 + 5
  Tool returned: 15
  Agent: FINAL:15

  User: tell me a joke
  Agent: FINAL:Why did the scarecrow win an award? Because he was outstanding in his field.
`, a.Personality.String(),

		a.Tools.Description(),
		a.Skills.String(),
		a.Memory.Recall(),
		a.Plan.String(),
	)
}

func (a *Agent) Run(goal string) string {
	a.Plan.Goal = goal
	a.Context.Append("user", goal)

	for step := 0; step < a.MaxSteps; step++ {
		fmt.Printf("\n-- step %d --\n", step+1)

		prompt := a.systemPrompt() + "\n\nConversation so far:\n" + a.Context.String()
		thought := llm.CallOllama(a.Model, prompt)
		fmt.Printf("thought: %s\n", thought)

		a.Context.Append("assistant", thought)
		a.Plan.AddStep(thought)

		// ← scan lines, not the whole blob
		action := extractAction(thought)

		switch {
		case strings.HasPrefix(action, "TOOL:"):
			rest := strings.TrimPrefix(action, "TOOL:")
			parts := strings.SplitN(rest, "|", 2)
			if len(parts) != 2 {
				a.Memory.Remember("malformed TOOL: " + action)
				continue
			}
			toolName := strings.TrimSpace(parts[0])
			toolInput := strings.TrimSpace(parts[1])
			result := a.Tools.Call(toolName, toolInput)
			fmt.Printf("tool %q -> %s\n", toolName, result)

			obs := fmt.Sprintf("tool %s returned: %s", toolName, result)
			a.Memory.Remember(obs)
			a.Context.Append("tool", obs)

		case strings.HasPrefix(action, "FINAL:"):
			answer := strings.TrimSpace(strings.TrimPrefix(action, "FINAL:"))
			a.Memory.Commit("completed: " + answer)
			return answer

		default:
			a.Memory.Remember(thought)
		}
	}

	return "max steps reached without a final answer"
}

func extractAction(response string) string {
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// canonical: TOOL:calculator|44 + 45
		if strings.HasPrefix(line, "TOOL:") {
			return line
		}

		// canonical: FINAL:89
		if strings.HasPrefix(line, "FINAL:") {
			return line
		}

		// drift case: "Final answer: FINAL:89" or "Final answer: 89"
		if lower := strings.ToLower(line); strings.HasPrefix(lower, "final") {
			// try to find an embedded FINAL: token first
			if idx := strings.Index(line, "FINAL:"); idx != -1 {
				return strings.TrimSpace(line[idx:])
			}
			// fallback: strip any "Final answer:" / "Final:" prefix and wrap
			colonIdx := strings.Index(line, ":")
			if colonIdx != -1 {
				rest := strings.TrimSpace(line[colonIdx+1:])
				if rest != "" {
					return "FINAL:" + rest
				}
			}
		}
	}

	// nothing matched — treat the whole response as a final answer
	// better to surface it than loop forever
	trimmed := strings.TrimSpace(response)
	if trimmed != "" {
		return "FINAL:" + trimmed
	}
	return ""
}
