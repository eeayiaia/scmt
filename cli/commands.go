//command line interface
package cli

import (
	"fmt"
	"github.com/codegangsta/cli"
	//"github.com/eeayiaia/scmt/invoker"
	//"bytes"
)

/* add plugin, remove plugin, enable/disable plugin, list nodes (status), list plugins (status), status of
specific node, ..?
*/

type ActionFunction func(*cli.Context)

var commands []cli.Command = []cli.Command {
	{
		Name:        "install-plugin",
		Aliases:     []string{""},
		Usage:       "scmt install-plugin <plugin name>",
		UsageText:   "Installs the plugin listed in plugins.d",
		Description: "",
		ArgsUsage:   "<plugin name> as first argument followed by nodes to install on, if no nodes are listed the plugin is installed on every node",
		Category:    "Plugin",
		Action:      func(c *cli.Context) { installPlugin(c) },
	},
	{
		Name:        "uninstall-plugin",
		Aliases:     []string{""},
		Usage:       "scmt uninstall-plugin <plugin name>",
		UsageText:   "Uninstalls the plugin listed in plugins.d",
		Description: "",
		ArgsUsage:   "<plugin name> as first argument followed by nodes to uninstall from, if no nodes are listed the plugin is uninstalled on every node",
		Category:    "Plugin",
		Action:      func(c *cli.Context) { uninstallPlugin(c) },
	},
	{
		Name:        "node-status",
		Aliases:     []string{""},
		Usage:       "scmt node-status <node name | node ip>",
		UsageText:   "Presents status information on node",
		Description: "",
		ArgsUsage:   "A list of <node name | node ip> , if left blank, status on all nodes are presented",
		Category:    "Cluster information",
		Action:      func(c *cli.Context) { printNodeInfo(c) },
	},
}

func getCommands() []cli.Command {
	return commands
}

func AddCommandShort(name string, af ActionFunction) {
	AddCommand(name, []string{""}, "No usage set", "No usage text set", "No args usage description", "", af)
}

func AddCommand(name string, aliases []string, usage string, usageText string, argsUsage string, cat string, af ActionFunction){
	commands = append(commands, cli.Command {
			Name:        name,
			Aliases:     aliases,
			Usage:       usage,
			UsageText:   usageText,
			ArgsUsage:   argsUsage,
			Category:    cat,
			Action:      af,
		})
}

func installPlugin(c *cli.Context) {
	//TODO: handle installation of plugins on specific nodes only, examine
	/*buffer := bytes.NewBufferString(c.Args().First())
	invoker.SendPacket(invoker.TYPE_INSTALL_PLUGIN, *buffer)*/
	fmt.Println("installing plugin " + c.Args().First() + " :Not implemented")
}

func uninstallPlugin(c *cli.Context) {
	fmt.Println("uninstalling plugin " + c.Args().First() + " :Not implemented")
}

func printNodeInfo(c *cli.Context) {
	fmt.Println("Not implemented: print node info")
}
