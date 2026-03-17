package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadFromDir(dir string) (Skills, error) {
	result := make(Skills)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return result, fmt.Errorf("cannot read skills dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		mdPath := filepath.Join(dir, entry.Name(), "skill.md")
		data, err := os.ReadFile(mdPath)
		if err != nil {
			fmt.Printf("skipping %s: %v\n", entry.Name(), err)
			continue
		}
		skill, err := parseSkillMD(string(data))
		if err != nil {
			fmt.Printf("skipping %s: %v\n", entry.Name(), err)
			continue
		}
		result.Register(skill)
		fmt.Printf("loaded skill: %s — %s\n", skill.Name, skill.Brief)
	}

	return result, nil
}

func parseSkillMD(content string) (Skill, error) {
	parts := strings.SplitN(content, "---", 2)
	if len(parts) != 2 {
		return Skill{}, fmt.Errorf("missing '---' separator")
	}

	s := Skill{Template: strings.TrimSpace(parts[1])}

	for _, line := range strings.Split(strings.TrimSpace(parts[0]), "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch strings.TrimSpace(k) {
		case "name":
			s.Name = strings.TrimSpace(v)
		case "brief":
			s.Brief = strings.TrimSpace(v)
		}
	}

	if s.Name == "" {
		return Skill{}, fmt.Errorf("missing 'name' field")
	}
	return s, nil
}
