package tools

import (
	"fmt"
	"strings"
)

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
