//command line interface
package cli

import (
	"github.com/codegangsta/cli"
	"fmt"
)

/* add plugin, remove plugin, enable/disable plugin, list nodes (status), list plugins (status), status of 
specific node, ..? 
*/

var Commands []cli.Command = []cli.Command {
  	{
  		Name:			"welcome",
  		Aliases:		[]string{"w"},
  		Usage:			"",
  		UsageText:		"",
  		Description:	"",
  		ArgsUsage:		"",
  		Category:		"",
  		Action:			func (c *cli.Context) { fmt.Println("Welcome to scmt")},

  	},
  	{
  		Name:			"add",
  		Aliases:		[]string{"Install, INSTALL"},
  		Usage:			"",
  		UsageText:		"",
  		Description:	"",
  		ArgsUsage:		"",
  		Category:		"",
  		Action:			func (c *cli.Context) {
							installPlugin(c)
						},
	},

}

func getCommands() []cli.Command {
	return Commands
}

func installPlugin(c *cli.Context) {
	fmt.Println("installing plugin " + c.Args().First() + " :Not implemented")
	return
}


