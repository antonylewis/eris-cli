package remotes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	def "github.com/eris-ltd/eris-cli/definitions"

	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/BurntSushi/toml"
	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

func WriteRemoteDefinitionFile(remoteDef *def.RemoteDefinition, fileName string) error {
	// writer := os.Stdout

	if filepath.Ext(fileName) == "" {
		fileName = remoteDef.Name + ".toml"
		fileName = filepath.Join(RemotesPath, fileName)
	}

	writer, err := os.Create(fileName)
	defer writer.Close()
	if err != nil {
		return err
	}

	switch filepath.Ext(fileName) {
	case ".json":
		mar, err := json.MarshalIndent(remoteDef, "", "  ")
		if err != nil {
			return err
		}
		mar = append(mar, '\n')
		writer.Write(mar)
	case ".yaml":
		mar, err := yaml.Marshal(remoteDef)
		if err != nil {
			return err
		}
		mar = append(mar, '\n')
		writer.Write(mar)
	default:
		writer.Write([]byte("# This is a TOML config file.\n# For more information, see https://github.com/toml-lang/toml\n\n"))
		enc := toml.NewEncoder(writer)
		enc.Indent = ""
		writer.Write([]byte("name = \"" + remoteDef.Name + "\"\n"))
		writer.Write([]byte("nodes = \"" + strconv.Itoa(remoteDef.Nodes) + "\"\n"))
		writer.Write([]byte("driver = \"" + "digitalocean" + "\"\n"))
		writer.Write([]byte("regions = " + "[ \"all\" ]" + "\n"))
		writer.Write([]byte("services = " + "[ \"all\" ]" + "\n")) //imgs to pull in
		writer.Write([]byte("machines = " + "[\n" + fmtMachs(remoteDef.Machines) + "\n]\n"))
		writer.Write([]byte("\n[maintainer]\n"))
		enc.Encode(remoteDef.Maintainer)

	}
	return nil
}

func fmtMachs(machines []string) string {
	k := len(machines) - 1
	for i, mach := range machines {
		machines[i] = fmt.Sprintf("\t\"%s\",", mach)
		//last string in arr => no comma
		if k == i {
			machines[i] = fmt.Sprintf("\t\"%s\"", mach)
		}
	}
	return strings.Join(machines, "\n")
}
