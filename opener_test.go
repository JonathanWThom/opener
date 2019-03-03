package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type TestSystem struct{}

var commands []string
var oldPath string
var oldSys System
var oldArgs []string

// Run mocks exec.Command().Run()
func (r TestSystem) Run(command string, args ...string) error {
	joined := strings.Join(args, " ")
	cmd := fmt.Sprintf("%s %s", command, joined)
	commands = append(commands, cmd)

	return nil
}

// Output mocks exec.Command().Output()
// It is used when checking if apps are open before closing them, so here
// returns "running" so the pkill command will be run
func (r TestSystem) Output(command string, args ...string) ([]byte, error) {
	return []byte("running"), nil
}

// TestOpenerDefault ensures that commands to open apps in default list are run
// Tests opener
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

// TestCloseDefault ensures that commands to close apps in default list are run
// Tests opener -c
func TestCloseDefault(t *testing.T) {
	setup()
	defer cleanup()

	// need to stub that they're running too
	os.Args = []string{"cmd", "-c"}
	Run()

	cases := []string{
		"pkill -SIGINT -a Google Chrome Canary",
		"pkill -SIGINT -a Slack",
		"pkill -SIGINT -a Calendar",
		"pkill -SIGINT -a Notion",
		"pkill -SIGINT -a Postico",
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
	oldArgs = os.Args
	sys = TestSystem{}
	path = "./applications.json"
}

func cleanup() {
	path = oldPath
	sys = oldSys
	os.Args = oldArgs
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
