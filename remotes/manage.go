package remotes

import (
	"path/filepath"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

func NewRemote(do *definitions.Do) error {
	rem := definitions.BlankRemoteDefinition()
	rem.Name = do.Name
	rem.Nodes = do.Nodes

	var err error
	//get maintainer info
	rem.Maintainer.Name, rem.Maintainer.Email, err = config.GitConfigUser()
	if err != nil {
		log.Debug(err.Error())
	}

	log.WithFields(log.Fields{
		"name":  rem.Name,
		"nodes": rem.Nodes,
	}).Debug("Creating a new remote definition file")
	err = WriteRemoteDefinitionFile(rem, filepath.Join(RemotesPath, rem.Name+".toml"))
	if err != nil {
		return err
	}
	do.Result = "success"
	return nil
}

func Add(args []string) {

}

func List() {

}

func Edit(args []string) {

}

func Rename(args []string) {

}

func Remove(args []string) {

}
