package mcndirs

import (
	"os"
	"path/filepath"

	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/mcnutils"
)

var (
	BaseDir = os.Getenv("MACHINE_STORAGE_PATH")
)

func GetBaseDir() string {
	if BaseDir == "" {
		BaseDir = filepath.Join(mcnutils.GetHomeDir(), ".docker", "machine")
	}
	return BaseDir
}

func GetMachineDir() string {
	return filepath.Join(GetBaseDir(), "machines")
}

func GetMachineCertDir() string {
	return filepath.Join(GetBaseDir(), "certs")
}
