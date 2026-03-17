package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func runCommand(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return "error: empty command"
	}

	// timeout guard — don't let the agent hang forever
	cmd := exec.Command("sh", "-c", input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		out := strings.TrimSpace(stdout.String())
		errOut := strings.TrimSpace(stderr.String())

		if err != nil {
			if errOut != "" {
				return fmt.Sprintf("error: %s", errOut)
			}
			return fmt.Sprintf("error: %s", err)
		}
		if out == "" {
			return "(no output)"
		}
		return out

	case <-time.After(10 * time.Second):
		cmd.Process.Kill()
		return "error: command timed out after 10s"
	}
}

func ShellTool() Tool {
	return Tool{
		Name:        "shell",
		Description: "executes a shell command and returns stdout. input: any valid shell command. example: 'ls -la' or 'echo hello' or 'cat file.txt'",
		Fn:          runCommand,
	}
}
