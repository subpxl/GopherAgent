package agentcore

import (
	"log"
	"os"
	"strings"
)

type Personality struct {
	Name        string
	Description string
	Traits      []string
	Rules       []string
}

func LoadPersonality(filepath string) (*Personality, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	content := string(file)
	p := &Personality{}
	lines := strings.Split(content, "\n")
	var currentSection string
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "name:") {
			p.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "##") {
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "##"))
		} else if strings.HasPrefix(line, "-") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			switch currentSection {
			case "Traits":
				p.Traits = append(p.Traits, value)
			case "Rules":
				p.Rules = append(p.Rules, value)
			}
		} else if currentSection == "Description" && line != "" {
			p.Description += line + " "
		}
	}

	return p, nil
}

func (p *Personality) String() string {
	return p.Name + "\n" +
		p.Description + "\n" +
		"Traits: " + strings.Join(p.Traits, ", ") + "\n" +
		"Rules: " + strings.Join(p.Rules, ", ")
}
