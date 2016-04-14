package cli

import (
	"testing"
	"fmt"
	"github.com/codegangsta/cli"
	"flag"
)

func TestGetCommands (t *testing.T) {
	command := getCommands()
	length := len(command)

	if len(command) != 3 {
		t.Errorf("Length of command list is not 3 it is %s", length)
	}
}

func TestAddCommand (t *testing.T) {
	AddCommand("test", []string{""}, "", "", "", "", func(c *cli.Context) { fmt.Println("the test action method") })

	command := getCommands()
	if len(command) != 4 {
		t.Errorf("Add command has not added a command!")
	}

	context := cli.NewContext(cli.NewApp(), flag.NewFlagSet("testflag", 1), nil)
	commandToRun := command[3]

	err := commandToRun.Run(context)

	if err != nil {
		t.Errorf("Failed to run test command!")
	}
}