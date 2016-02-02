package remotes

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/list"
	"github.com/eris-ltd/eris-cli/loaders"
	"github.com/eris-ltd/eris-cli/perform"
	"github.com/eris-ltd/eris-cli/util"

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
		"name":            rem.Name,
		"number of nodes": rem.Nodes,
	}).Debug("Creating a new remote definition file")

	rem.Machines = make([]string, rem.Nodes)
	for i, _ := range rem.Machines {
		rem.Machines[i] = fmt.Sprintf("eris-remote-%s-%v", rem.Name, i)
	}

	err = WriteRemoteDefinitionFile(rem, filepath.Join(RemotesPath, rem.Name+".toml"))
	if err != nil {
		return err
	}
	do.Result = "success"
	return nil
}

func ListRemotes(do *definitions.Do) error {
	rems, err := list.ListKnown("remotes")
	if err != nil {
		return err
	}
	if rems == "" {
		log.Warn("No known remotes")
		return nil
	}

	knowns := strings.Split(rems, "\n")
	log.WithField("=>", knowns[0]).Warn("The known remotes on your host kind marmot:")
	knowns = append(knowns[:0], knowns[1:]...)
	for _, known := range knowns {
		log.WithField("=>", known).Warn()
	}

	hosts, err := list.ListExistingRemotes()
	if err != nil {
		return err
	}
	if len(hosts) == 0 {
		log.Warn("No hosts found for that remote.")
		return nil
	}

	log.WithField("=>", hosts[0]).Warn("The existing remote hosts kind marmot:")
	hosts = append(hosts[:0], hosts[1:]...)
	for _, host := range hosts {
		log.WithField("=>", host).Warn()
	}
	return nil
}

func EditRemote(do *definitions.Do) error {
	remDefFile := FindRemoteDefinitionFile(do.Name)
	log.WithField("=>", remDefFile).Info("Editing remote")
	//do.Result = "success"
	return Editor(remDefFile)

}

func CatRemote(do *definitions.Do) error {
	configs := util.GetGlobalLevelConfigFilesByType("remotes", true)
	for _, c := range configs {
		cName := strings.Split(filepath.Base(c), ".")[0]
		if cName == do.Name {
			cat, err := ioutil.ReadFile(c)
			if err != nil {
				return err
			}
			//do.Result = string(cat)
			log.Warn(string(cat))
			return nil
		}
	}
	return fmt.Errorf("Unknown remote %s or invalid file extension", do.Name)
}

func RemoveRemote(do *definitions.Do) error {
	remDef, err := loaders.LoadRemoteDefinition(do.Name)
	if err != nil {
		return err
	}

	if err := perform.RemoveRemote(remDef); err != nil {
		return err
	}
	return nil
}

func Rename(do *definitions.Do) error {
	return nil
}
