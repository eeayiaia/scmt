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
  		Name:			"install-plugin",
  		Aliases:		[]string{""},
  		Usage:			"scmt install-plugin <plugin name>",
  		UsageText:		"Installs the plugin listed in plugins.d",
  		Description:	"",
  		ArgsUsage:		"<plugin name> as first argument followed by nodes to install on, if no nodes are listed
  						the plugin is installed on every node",
  		Category:		"Plugin",
  		Action:			func (c *cli.Context) { installPlugin(c) },

  	},
  	{
  		Name:			"uninstall-plugin",
  		Aliases:		[]string{""},
  		Usage:			"scmt uninstall-plugin <plugin name>",
  		UsageText:		"Uninstalls the plugin listed in plugins.d",
  		Description:	"",
  		ArgsUsage:		"<plugin name> as first argument followed by nodes to uninstall from, if no nodes are listed
  						the plugin is uninstalled on every node",
  		Category:		"Plugin",
  		Action:			func (c *cli.Context) { uninstallPlugin(c) },
	},

}

func getCommands() []cli.Command {
	return Commands
}

func installPlugin(c *cli.Context) {
	fmt.Println("installing plugin " + c.Args().First() + " :Not implemented")
}

func uninstallPlugin(c *cli.Context) {
	fmt.Println("uninstalling plugin " + c.Args().First() + " :Not implemented")
}


