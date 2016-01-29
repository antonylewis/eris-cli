package loaders

import (
	"fmt"
	//"path/filepath"
	//"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	//"github.com/eris-ltd/eris-cli/util"
	//"github.com/eris-ltd/eris-cli/version"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"

	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/viper"
)

func LoadRemoteDefinition(remName string) (*definitions.RemoteDefinition, error) {

	rem := definitions.BlankRemoteDefinition()
	remConf, err := loadRemoteDefinition(remName)
	if err != nil {
		return nil, err
	}

	if err = MarshalRemoteDefinition(remConf, rem); err != nil {
		return nil, err
	}

	if rem.Name == "" {
		return nil, fmt.Errorf("No remote name given.")
	}

	if rem.Nodes == 0 || rem.Nodes > 50 {
		return nil, fmt.Errorf("Number of nodes is either 0 or greater than 50.")
	}

	if rem.Driver != "digitalocean" {
		return nil, fmt.Errorf("Driver specified is not Digital Ocean.")
	}

	if len(rem.Machines) != rem.Nodes {
		erro := fmt.Sprintf("Number of machines (%v) not equal nodes (%v)", len(rem.Machines), rem.Nodes)
		return nil, fmt.Errorf(erro)
	}

	return rem, nil

}

func MarshalRemoteDefinition(remConf *viper.Viper, rem *definitions.RemoteDefinition) error {
	err := remConf.Unmarshal(rem)
	if err != nil {
		log.WithField("=>", fmt.Sprintf("%v", err)).Warn("Unmarshal error:")

		return fmt.Errorf("Sorry, the marmots could not figure that remote definition out.\nPlease check for known remote with [eris remotes ls --known] and retry.\n")
	}
	return nil
}

func loadRemoteDefinition(remName string) (*viper.Viper, error) {
	return config.LoadViperConfig(RemotesPath, remName, "remote")
}
