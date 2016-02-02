package commands

import (
	"strconv"

	rem "github.com/eris-ltd/eris-cli/remotes"

	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
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
	Use:   "new NAME NODES", //ax NODES in favour of flag
	Short: "Command-line tool to deploy new remotes",
	Long: `Tool will prompt for deploy options.
	Requires docker-machine installed and a Digital Ocean API Token.`,
	Run: NewRemote,
}

var remotesInit = &cobra.Command{
	Use:   "init NAME",
	Short: "Initialize remotes from remotes definition file.",
	Long: `Initialize remotes from remotes definition file.
	Creates N machines & pulls specified service images.`,
	Run: InitRemote,
}

var remotesList = &cobra.Command{
	Use:   "ls",
	Short: "List all registered remotes.",
	Long:  `List all registered remotes`,
	Run:   ListRemotes,
}

var remotesDo = &cobra.Command{
	Use:   "do NAME ACTION",
	Short: "Perform an action on a remote.",
	Long:  `Perform an action on a remote according to the action definition file.`,
	Run:   DoRemote,
}

var remotesEdit = &cobra.Command{
	Use:   "edit NAME",
	Short: "Edit a remote definition file.",
	Long:  `Edit a remote definition file`,
	Run:   EditRemote,
}

var remotesProv = &cobra.Command{
	Use:   "prov NAME",
	Short: "",
	Long:  ``,
	Run:   ProvRemote,
}

var remotesCat = &cobra.Command{
	Use:   "cat NAME",
	Short: "Cat a remote definition file.",
	Long:  `Cat a remote definition file`,
	Run:   CatRemote,
}

var remotesRename = &cobra.Command{
	Use:   "rename OLD NEW",
	Short: "Rename a remote and all its hosts.",
	Long:  `Rename a remote and all its hosts.`,
	Run:   RenameRemote,
}

var remotesRemove = &cobra.Command{
	Use:   "rm NAME",
	Short: "Remove a remote and all its hosts.",
	Long:  `Remove a remote and all its hosts.`,
	Run:   RemoveRemote,
}

func buildRemotesFlags() {

}

// all remotes take a single arg (NAME)
// save for new & ls
func NewRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(2, "eq", cmd, args))

	var err error
	do.Name = args[0]
	do.Nodes, err = strconv.Atoi(args[1])
	IfExit(err)
	IfExit(rem.NewRemote(do))
}

func InitRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.Init(do))
}

func ListRemotes(cmd *cobra.Command, args []string) {
	//flags only (coming)
	IfExit(ArgCheck(0, "eq", cmd, args))
	IfExit(rem.ListRemotes(do))
}

func DoRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.Do(do))
}

func EditRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.EditRemote(do))
}

func ProvRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.ProvRemote(do))
}

func CatRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.CatRemote(do))
}

func RenameRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.Rename(do))
}

func RemoveRemote(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))
	do.Name = args[0]
	IfExit(rem.RemoveRemote(do))
}
