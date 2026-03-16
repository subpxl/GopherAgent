package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
}

func calculate(input string) string {
	var a, b float64
	var op string

	_, err := fmt.Sscanf(input, "%f %s %f", &a, &op, &b)
	if err != nil {
		return "error: expected format 'a op b', got: " + input
	}

	switch op {
	case "+":
		return fmt.Sprintf("%g", a+b)
	case "-":
		return fmt.Sprintf("%g", a-b)
	case "*":
		return fmt.Sprintf("%g", a*b)
	case "/":
		if b == 0 {
			return "error: division by zero"
		}
		return fmt.Sprintf("%g", a/b)
	default:
		return "error: unknown operator: " + op
	}
}

func CalculatorTool() Tool {
	return Tool{
		Name:        "calculator",
		Description: "performs arithmetic on two numbers. input format: 'a op b' where op is +, -, *, /. example: '12.5 * 3'",
		Fn:          calculate,
	}
}

// Skill
func MathSkill() Skill {
	return Skill{
		Name: "math-solver",
		SystemTemplate: strings.TrimSpace(`
You are a math assistant. Given a problem, break it into ONE arithmetic step at a time.
For each step call the calculator tool with format: 'a op b'
Problem: {input}`),
	}
}

func extractAction(response string) string {
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TOOL:") || strings.HasPrefix(line, "FINAL:") {
			return line
		}
	}
	return response // no action found, return raw for default case
}
func main() {
	print("agent  451")
	agent := NewAgent("qwen3:4b-instruct")

	agent.Tools.Register(CalculatorTool())

	agent.Skills.Register(MathSkill())

	agent.Tools.Register(Tool{
		Name:        "search",
		Description: "searches the web for a query, input: a search query string",
		Fn: func(input string) string {
			// stub — wire in a real search API
			return fmt.Sprintf("(search stub: top result for %q)", input)
		},
	})

	// test the skill

	prompt, _ := agent.Skills.Render("math-solver", "what is 128 divided by 4, then multiplied by 4555555?")
	answer := agent.Run(prompt, 6)
	fmt.Println(answer)

}

func callOllama(model, prompt string) string {
	reqBody := Request{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "FINAL: LLM unavailable — " + err.Error()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result Response
	json.Unmarshal(body, &result)

	return strings.TrimSpace(result.Response)

}

type ToolFunc func(input string) string

type Tool struct {
	Name        string
	Description string
	Fn          ToolFunc
}

type Tools map[string]Tool

func (t Tools) Register(tool Tool) {
	t[tool.Name] = tool
}

func (t Tools) Call(name, input string) string {
	tool, ok := t[name]
	if !ok {
		return fmt.Sprintf("error:unknown tool %q", name)
	}
	return tool.Fn(input)
}

func (t Tools) Description() string {
	var sb strings.Builder
	for _, tool := range t {
		fmt.Fprintf(&sb, "-%s: %s\n", tool.Name, tool.Description)
	}
	return sb.String()
}

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

// skills
type Skill struct {
	Name           string
	SystemTemplate string
}

type Skills map[string]Skill

func (s Skills) Register(skill Skill) {
	s[skill.Name] = skill
}

func (s Skills) Render(name, input string) (string, error) {
	skill, ok := s[name]
	if !ok {
		return "", fmt.Errorf("unknown skill %q", name)
	}
	return strings.ReplaceAll(skill.SystemTemplate, "{input}", input), nil
}

func (s Skills) String() string {
	var sb strings.Builder
	for _, skill := range s {
		fmt.Fprintf(&sb, "-%s: %s\n", skill.Name, skill.SystemTemplate)
	}
	return sb.String()

}

// context

type Context struct {
	History []string
}

func (c *Context) Append(role, msg string) {
	c.History = append(c.History, fmt.Sprintf("[%s] %s", role, msg))

}

func (c *Context) String() string {
	return strings.Join(c.History, "\n")
}

// Agent

type Agent struct {
	Model   string
	Plan    Plan
	Context Context
	Tools   Tools
	Skills  Skills
	Memory  *Memory
}

func NewAgent(model string) *Agent {
	return &Agent{
		Model:  model,
		Tools:  make(Tools),
		Skills: make(Skills),
		Memory: NewMemory(10),
	}
}

func (a *Agent) systemPrompt() string {
	return fmt.Sprintf(`you are a autonomus agent. reason step by step.
	
	Available Tools:
	%s
	
	To use a tool, respons Exactly with:
	TOOL:<tool_name>|<input>

	Available Skills:
	%s

	CURRENT memory:
	%s

	Current Plan:
	%s

	When you have the final answer, respond with:
	FINAL:<your answer>
	`, a.Tools.Description(), a.Skills.String(), a.Memory.Recall(), a.Plan.String())

}

func (a *Agent) Run(goal string, maxSteps int) string {
	a.Plan.Goal = goal
	a.Context.Append("user", goal)

	for step := 0; step < maxSteps; step++ {
		fmt.Printf("\n-- step %d --\n", step+1)

		prompt := a.systemPrompt() + "\n\nConversation so far:\n" + a.Context.String()
		thought := callOllama(a.Model, prompt)
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
			fmt.Printf("tool %q -> %s\n", toolName, result) // ← fixed format

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
