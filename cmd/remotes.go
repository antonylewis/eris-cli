package commands

import (
	"strconv"

	//"github.com/eris-ltd/eris-cli/list"
	rem "github.com/eris-ltd/eris-cli/remotes"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

//TODO add ArgCheck -> call rem pkg indirectly
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
	Remotes.AddCommand(remotesList)
	//Remotes.AddCommand(remotesDo)
	Remotes.AddCommand(remotesEdit)
	Remotes.AddCommand(remotesProv)
	Remotes.AddCommand(remotesCat)
	//Remotes.AddCommand(remotesRename)
	Remotes.AddCommand(remotesRemove)
}

var remotesNew = &cobra.Command{
	Use:   "new NAME NODES", //eventuall DRIVER
	Short: "Command-line tool to deploy new remotes",
	Long: `Tool will prompt for deploy options.
	Requires docker-machine installed and a Digital Ocean API Token.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		do.Name = args[0]
		do.Nodes, err = strconv.Atoi(args[1])
		if err != nil {
			log.Warn("strconv err:")
			log.Error(err)
		}
		rem.NewRemote(do)
	},
}

var remotesInit = &cobra.Command{
	Use:   "init NAME",
	Short: "Initialize remotes from remotes definition file.",
	Long: `Initialize remotes from remotes definition file.
	Creates N machines & pulls specified service images.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Init(args)
	},
}

var remotesList = &cobra.Command{
	Use:   "ls",
	Short: "List all registered remotes.",
	Long:  `List all registered remotes`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.ListRemotes()
	},
}

var remotesDo = &cobra.Command{
	Use:   "do NAME ACTION",
	Short: "Perform an action on a remote.",
	Long:  `Perform an action on a remote according to the action definition file.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Do(args)
	},
}

var remotesEdit = &cobra.Command{
	Use:   "edit NAME",
	Short: "Edit a remote definition file.",
	Long:  `Edit a remote definition file`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = rem.EditRemote(args)
	},
}

var remotesProv = &cobra.Command{
	Use:   "prov NAME",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = rem.ProvRemote(args)
	},
}

var remotesCat = &cobra.Command{
	Use:   "cat NAME",
	Short: "Cat a remote definition file.",
	Long:  `Cat a remote definition file`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = rem.CatRemote(args)
	},
}

var remotesRename = &cobra.Command{
	Use:   "rename OLD NEW",
	Short: "Rename a remote and all its hosts.",
	Long:  `Rename a remote and all its hosts.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.Rename(args)
	},
}

var remotesRemove = &cobra.Command{
	Use:   "rm NAME",
	Short: "Remove a remote and all its hosts.",
	Long:  `Remove a remote and all its hosts.`,
	Run: func(cmd *cobra.Command, args []string) {
		rem.RemoveRemote(args)
	},
}

func buildRemotesFlags() {

}

func ListRemotes() {
	if err := rem.ListRemotes(); err != nil {
		return
	}
}
