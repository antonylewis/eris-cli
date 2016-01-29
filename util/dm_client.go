package util

import (
	//"fmt"
	//"os"
	"path/filepath"

	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine"
	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

//var api libmachine.API

func RemotesConnect(remName string) (api libmachine.API, err error) {
	//check if dm is actually instaled!

	api = libmachine.NewClient(GetBaseDir(), GetMachineCertDir())
	defer api.Close()

	//fmt.Print("Get Machines Dir: %s\n", api.GetMachinesDir())

	return api, nil

}

//---from github.com/docker/machine/commands/mcndirs/util.go"
// modified to suit our purposes

/*var (
	BaseDir = os.Getenv("MACHINE_STORAGE_PATH")
)*/

func GetBaseDir() string {
	//	if BaseDir == "" {
	BaseDir := filepath.Join(RemotesPath, "machine")
	//	}
	return BaseDir
}

func GetMachineDir() string {
	return filepath.Join(GetBaseDir(), "machines")
}

func GetMachineCertDir() string {
	return filepath.Join(GetBaseDir(), "certs")
}
