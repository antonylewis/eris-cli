package commands

import (
	"strconv"

	rem "github.com/eris-ltd/eris-cli/remotes"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

var Remotes = &cobra.Command{
	Use:   "remotes",
	Short: "Manage and Perform Remote Machines and Services.",
	Long: `Display and Manage remote machines which are operating
various services reachable by the Eris platform.

Actions, if configured as such, can utilize remote machines.
To register and manage remote machines for sending of actions
to those machines, use this command.`,
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

// build the services subcommand
func buildRemotesCommand() {
	Remotes.AddCommand(remotesNew)
	Remotes.AddCommand(remotesInit)
	//Remotes.AddCommand(remotesLoad)
	//Remotes.AddCommand(remotesAdd)
	//Remotes.AddCommand(remotesList)
	//Remotes.AddCommand(remotesDo)
	//Remotes.AddCommand(remotesEdit)
	//Remotes.AddCommand(remotesRename)
	//Remotes.AddCommand(remotesRemove)
}

var remotesNew = &cobra.Command{
	Use:   "new NAME NODES",
	Short: "Command-line tool to deploy new remotes",
	Long: `Tool will prompt for deploy options.
	Requires docker-machine installed and a Digital Ocean API Token.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		do.Name = args[0]
		do.Nodes, err = strconv.Atoi(args[1])
		if err != nil {
			log.Warn("fuck")
		}
		rem.NewRemote(do)
	},
}

var remotesInit = &cobra.Command{
	Use:   "init NAME", //Not yet like that
	Short: "Initialize remotes from remotes definition file.",
	Long: `Initialize remotes from remotes definition file.
	Creates N machines & pulls specified service images.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Init(args)
	},
}

var remotesLoad = &cobra.Command{
	Use:   "load NAME",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		rem.Load(args[0])
	},
}

var remotesAdd = &cobra.Command{
	Use:   "add [name] [remote-definition-file]",
	Short: "Adds a remote to Eris.",
	Long:  `Adds a remote to Eris in JSON, TOML, or YAML format.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Add(args)
	},
}

// flags to add: --global --project
var remotesList = &cobra.Command{
	Use:   "ls",
	Short: "List all registered remotes.",
	Long:  `List all registered remotes`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.List()
	},
}

var remotesDo = &cobra.Command{
	Use:   "do [name]",
	Short: "Perform an action on a remote.",
	Long:  `Perform an action on a remote according to the action definition file.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Do(args)
	},
}

var remotesEdit = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a remote definition file.",
	Long:  `Edit a remote definition file`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Edit(args)
	},
}

var remotesRename = &cobra.Command{
	Use:   "rename [old] [new]",
	Short: "Rename a remote.",
	Long:  `Rename a remote`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Rename(args)
	},
}

var remotesRemove = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a remote definition file.",
	Long:  `Remove a remote definition file`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Remove(args)
	},
}
