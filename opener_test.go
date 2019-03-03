package main

import (
	"fmt"
	"strings"
	"testing"
)

type TestSystem struct{}

var commands []string
var oldPath string
var oldSys System

// Run mocks exec.Command().Run()
func (r TestSystem) Run(command string, args ...string) error {
	joined := strings.Join(args, " ")
	cmd := fmt.Sprintf("%s %s", command, joined)
	commands = append(commands, cmd)

	return nil
}

// TestOpenerDefault ensures that commands to open apps in default list are run
func TestOpenerDefault(t *testing.T) {
	setup()
	defer cleanup()

	Run()

	cases := []string{
		"open -a Google Chrome Canary",
		"open -a Slack",
		"open -a Calendar",
		"open -a Notion",
		"open -a Postico",
	}

	for _, c := range cases {
		if !contains(commands, c) {
			t.Errorf("Command was not executed: %s", c)
		}
	}
}

// Helpers
func setup() {
	oldPath = path
	oldSys = sys
	sys = TestSystem{}
	path = "./applications.json"
}

func cleanup() {
	path = oldPath
	sys = oldSys
	commands = []string{}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
