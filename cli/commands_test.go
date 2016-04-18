package cli

import (
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"testing"
)

func TestGetCommands(t *testing.T) {
	command := getCommands()
	length := len(command)

	if len(command) != len(commands) {
		t.Errorf("Length of command list is not 3 it is %s", length)
	}
}

func TestAddCommand(t *testing.T) {
	formerLen := len(getCommands())

	AddCommand("test", []string{""}, "", "", "", "", func(c *cli.Context) { fmt.Println("the test action method") })

	command := getCommands()
	if len(command) != formerLen+1 {
		t.Errorf("Add command has not added a command!")
	}

	context := cli.NewContext(cli.NewApp(), flag.NewFlagSet("testflag", 1), nil)
	commandToRun := command[len(commands)-1]

	err := commandToRun.Run(context)

	if err != nil {
		t.Errorf("Failed to run test command!")
	}
}
