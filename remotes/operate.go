package remotes

import (
	"fmt"
	"strings"

	"github.com/eris-ltd/eris-cli/loaders"
	"github.com/eris-ltd/eris-cli/perform"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func Init(args []string) {
	remDef, err := loaders.LoadRemoteDefinition(args[0])
	if err != nil {
		log.Warn(fmt.Sprintf("error loading remote defintion: %v", err))
	}
	log.WithFields(log.Fields{
		"name":     remDef.Name,
		"nodes":    remDef.Nodes,
		"driver":   remDef.Driver,
		"regions":  strings.Join(remDef.Regions, ","),
		"services": strings.Join(remDef.Services, ","),
	}).Warn("Initializing remote:")

	log.Debug("With machines:")
	for _, mach := range remDef.Machines {
		log.WithField("=>", mach).Debug()
	}

	if err := perform.CreateRemote(remDef); err != nil {
		log.Error(err)
	}

	/*
		if err := perform.PullImagesToRemote(remDef); err != nil {
			log.Error(err)
		}*/

}

func ProvRemote(args []string) error {
	remDef, err := loaders.LoadRemoteDefinition(args[0])
	if err != nil {
		log.Warn(fmt.Sprintf("error loading remote defintion: %v", err))
	}
	log.Warn("reprovisioning remote")
	if err := perform.ReprovisionRemote(remDef); err != nil {
		log.Error(err)
	}
	return nil
}

func Start() {
}

func Stop() {
}
