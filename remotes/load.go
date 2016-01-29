package remotes

import (
	"github.com/eris-ltd/eris-cli/util"
)

func FindRemoteDefinitionFile(name string) string {
	return util.GetFileByNameAndType("remotes", name)
}
