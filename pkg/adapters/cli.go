package adapters

import (
	"agent-in-go/pkg/session"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLIAdapter struct {
	store *session.SessionStore
}

func NewCLIAdapter(store *session.SessionStore) *CLIAdapter {
	return &CLIAdapter{store: store}
}

func (a *CLIAdapter) Name() string { return "CLI" }

func (a *CLIAdapter) Start() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("[CLI] agent 451 ready. ctrl+c to exit")
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		} else if input == "bye" {
			fmt.Println("good bye")
			break
		}
		fmt.Println("answer:", a.store.Ask("cli", input))
	}
	return nil
}
