package main

import (
	"agent-in-go/pkg/adapters"
	"agent-in-go/pkg/agentcore"
	"agent-in-go/pkg/session"
	"agent-in-go/pkg/skills"
	"agent-in-go/pkg/tools"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Model               string
	MaxSteps            int
	PersonalityFilePath string
}

func setupAgent(config *Config) *agentcore.Agent {
	personality, err := agentcore.LoadPersonality(config.PersonalityFilePath)
	if err != nil {
		log.Fatal(err)
	}

	a := agentcore.NewAgent(config.Model, config.MaxSteps, personality)

	a.Tools.Register(tools.CalculatorTool())
	a.Tools.Register(tools.ShellTool())

	loadedSkills, err := skills.LoadFromDir("skills")
	if err != nil {
		log.Printf("warn: could not load skills: %v", err)
	}
	a.Skills = loadedSkills
	return a
}

func main() {
	config := &Config{
		Model:               "qwen3:4b-instruct",
		MaxSteps:            6,
		PersonalityFilePath: "personality.md",
	}

	store := session.NewSessionStore(func() *agentcore.Agent {
		return setupAgent(config)
	})

	adapterList := []adapters.Adapter{
		adapters.NewRESTAdapter("8080", store),
		adapters.NewCLIAdapter(store),
		adapters.NewWSAdapter("8081", store),
	}

	errCh := make(chan error, len(adapterList))
	for _, a := range adapterList {
		a := a
		go func() {
			fmt.Printf("starting %s adapter\n", a.Name())
			errCh <- a.Start()
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatal("adapter error:", err)
	case <-quit:
		fmt.Println("\nshutting down")
		os.Exit(0)
	}
}
