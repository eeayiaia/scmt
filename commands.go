//command line interface
package main

import (
	"github.com/codegangsta/cli"
	"github.com/eeayiaia/scmt/daemon"
	"github.com/eeayiaia/scmt/invoker"

	log "github.com/Sirupsen/logrus"

	"bytes"
	"fmt"
)

/* add plugin, remove plugin, enable/disable plugin, list nodes (status), list plugins (status), status of
specific node, ..?
*/

type ActionFunction func(*cli.Context)

var commands []cli.Command = []cli.Command{
	{
		Name:        "install-plugin",
		Aliases:     []string{""},
		Usage:       "scmt install-plugin <plugin name>",
		UsageText:   "Installs the plugin listed in plugins.d",
		Description: "",
		ArgsUsage:   "<plugin name> as first argument followed by nodes to install on, if no nodes are listed the plugin is installed on every node",
		Category:    "Plugin Control",
		Action:      func(c *cli.Context) { installPlugin(c) },
	},
	{
		Name:        "uninstall-plugin",
		Aliases:     []string{""},
		Usage:       "scmt uninstall-plugin <plugin name>",
		UsageText:   "Uninstalls the plugin listed in plugins.d",
		Description: "",
		ArgsUsage:   "<plugin name> as first argument followed by nodes to uninstall from, if no nodes are listed the plugin is uninstalled on every node",
		Category:    "Plugin Control",
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

	{
		Name:        "stop-daemon",
		Aliases:     []string{""},
		Usage:       "scmt stop-daemon",
		UsageText:   "Stops the daemon if it is running in the background",
		Description: "",
		ArgsUsage:   "",
		Category:    "Daemon Control",
		Action:      stopDaemon,
	},
	{
		Name:        "start-daemon",
		Aliases:     []string{""},
		Usage:       "scmt start-daemon",
		UsageText:   "Starts the daemon if it is not running in the background",
		Description: "",
		ArgsUsage:   "",
		Category:    "Daemon Control",
		Action:      startDaemon,
	},
}

func getCommands() []cli.Command {
	return commands
}

func AddCommandShort(name string, af ActionFunction) {
	AddCommand(name, []string{""}, "No usage set", "No usage text set", "No args usage description", "", af)
}

func AddCommand(name string, aliases []string, usage string, usageText string, argsUsage string, cat string, af ActionFunction) {
	commands = append(commands, cli.Command{
		Name:      name,
		Aliases:   aliases,
		Usage:     usage,
		UsageText: usageText,
		ArgsUsage: argsUsage,
		Category:  cat,
		Action:    af,
	})
}

func installPlugin(c *cli.Context) {
	//TODO: handle installation of plugins on specific nodes only, examine
	buffer := bytes.NewBufferString(c.Args().First())
	buffer.WriteString(" ")
	invoker.SendPacket(invoker.TYPE_INSTALL_PLUGIN, *buffer)
}

func uninstallPlugin(c *cli.Context) {
	fmt.Println("uninstalling plugin " + c.Args().First() + " :Not implemented")
}

func printNodeInfo(c *cli.Context) {
	fmt.Println("Not implemented: print node info")
}

func stopDaemon(c *cli.Context) {
	if !daemon.IsDaemonized() {
		log.Error("No daemon is running!")
		return
	}

	invoker.SendPacket(invoker.TYPE_STOP_DAEMON, *bytes.NewBufferString(""))
	log.Info("Stopping the daemon ..")
}

func startDaemon(c *cli.Context) {
	if daemon.IsDaemonized() {
		log.Error("Daemon is already running!")
		return
	}

	daemon.Daemonize(background, termHandler)
}
