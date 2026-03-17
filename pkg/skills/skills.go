package skills

import (
	"fmt"
	"strings"
)

// skills
type Skill struct {
	Name     string
	Brief    string
	Template string
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
	return strings.ReplaceAll(skill.Template, "{input}", input), nil
}

func (s Skills) String() string {
	var sb strings.Builder
	for _, skill := range s {
		fmt.Fprintf(&sb, "- %s: %s\n", skill.Name, skill.Brief)
	}
	return sb.String()
}
