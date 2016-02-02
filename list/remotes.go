package list

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eris-cli/util"
)

func ListExistingRemotes() ([]string, error) {

	api, err := util.RemotesConnect("")
	if err != nil {
		return []string{}, err
	}

	machines, err := api.List()
	if err != nil {
		fmt.Printf("Error listing machines: %v\n", err)
		os.Exit(1)
	}
	return machines, nil
}
