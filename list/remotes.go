package list

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eris-cli/util"
	//	"github.com/docker/machine/libmachine"
	// "github.com/docker/machine/libmachine/persist"
)

func ListExistingRemotes() ([]string, error) {

	api, err := util.RemotesConnect("")
	if err != nil {
		return []string{}, err
	}

	//f := persist.NewFilestore(util.GetBaseDir(), "", "")

	// list, err := libmachine.API.List(f)
	machines, err := api.List()
	if err != nil {
		fmt.Printf("Error listing machines: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Machines:\n%v\n", machines)
	return machines, nil
}
