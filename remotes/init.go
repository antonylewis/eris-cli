package remotes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//"strconv"
	"encoding/json"
	"sync"
	"time"

	//	def "github.com/eris-ltd/eris-cli/definitions"
	ini "github.com/eris-ltd/eris-cli/initialize"
	"github.com/eris-ltd/eris-cli/loaders"
	"github.com/eris-ltd/eris-cli/util"

	"github.com/docker/machine/libmachine/drivers"
	//"github.com/docker/machine/libmachine/drivers/rpc"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	//"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/drivers/digitalocean"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine"
	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

var api libmachine.API

var wg sync.WaitGroup
var maxTimeout = 15 * time.Minute
var script = "docker.sh"

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

	machines := make([]string, remDef.Nodes)
	for i, _ := range machines {
		mach := fmt.Sprintf("eris-%s-%v", remDef.Name, i)
		log.WithField("=>", mach).Warn("Gonna build machine:")
		machines[i] = mach
	}

	failOut := make(chan bool, len(machines))
	go timeOutTicker(machines)

	wg.Add(len(machines))
	for _, m := range machines {
		go createMachines(m, remDef.Driver, failOut)
	}

	go func(failOut chan bool) {
		if _, ok := <-failOut; ok {
			for _, m := range machines {
				fmt.Println(m)
			}
			os.Exit(1)
		}
	}(failOut)
	wg.Wait()

	for _, m := range machines {
		log.WithField("=>", m).Warn("Machine created:")
	}
	/*

		log.Warn("Initializing docker images:")
		for _, m := range machines {
			log.WithField("=>", m).Warn("Initializing docker images for machine:")
			if err := initMachine(m, remDef.Services); err != nil {
				log.Warn(fmt.Sprintf("some dumb error: %v\n", err))
			}
			//go initMachines(m, failOut)
		}
		/*go func(failOut chan bool) {
			if _, ok := <-failOut; ok {
				for _, m := range machines {
					fmt.Println(m)
				}
				os.Exit(1)
			}
		}(failOut)
		wg.Wait()*/
}

func dmRun(typ string, machines []string) {
}

func createMachines(name, driver string, failOut chan<- bool) {
	defer wg.Done()
	if err := createMachine(name, driver); err != nil {
		failOut <- true
		//return
	}
	/*if err := setUpMachine(name); err != nil {
		failOut <- true
	}*/
}

/*func initMachines(name string, failOut chan<- bool) error {
	defer wg.Done()
	if err := initMachine(name); err != nil {
		failOut <- true
	}
	return nil
}*/

// straight outta digitalocean.Drivers
type Driver struct {
	*drivers.BaseDriver
	AccessToken       string
	DropletID         int
	DropletName       string
	Image             string
	Region            string
	SSHKeyID          int
	Size              string
	IPv6              bool
	Backups           bool
	PrivateNetworking bool
	UserDataFile      string
}

const (
	defaultSSHPort = 22
	defaultSSHUSer = "root"
	defaultImage   = "ubuntu-15-10-x64"
	defaultRegion  = "nyc3"
	defaultSize    = "512mb"
)

func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		Image:  defaultImage,
		Size:   defaultSize,
		Region: defaultRegion,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

//TODO add region option
func createMachine(name, driverTyp string) error {
	//mach := filepath.Join(RemotesPath, "machines")
	//certs := "/Users/zicter/.docker/machine/certs"
	certs := filepath.Join(RemotesPath, "certs")
	//if err := os.MkdirAll(certs, 0777); err != nil {
	//	return err
	//}

	dafuq := filepath.Join(RemotesPath, "machines", name)
	if err := os.MkdirAll(dafuq, 0777); err != nil {
		return err
	}

	//client := libmachine.NewClient("/tmp/automatic", "/tmp/automatic/certs")
	client := libmachine.NewClient("/Users/zicter/.docker/machine/", certs)
	defer client.Close()

	driver := NewDriver(name, "/Users/zicter/.docker/machine/")

	driver.AccessToken = getToken()
	driver.DropletName = name
	driver.SSHUser = "root"
	driver.SSHPort = 22

	rawDriver, err := json.Marshal(driver)
	if err != nil {
		log.Warn("unmarshal error:")
		log.Error(err)
		return err
	}

	//fmt.Printf("RAW: %s\n", string(rawDriver))
	//driver should be "digitalocean", set in remDef file
	h, err := client.NewHost("digitalocean", rawDriver)
	if err != nil {
		log.Warn("newHost error:")
		log.Error(err)
		return err
	}

	/*if err := driver.PreCreateCheck(); err != nil {
		log.Warn("pre-create error:")
		log.Error(err)
		return err
	}*/

	//	h.HostOptions.EngineOptions.StorageDriver = "overlay"

	if err := client.Create(h); err != nil {
		log.Warn("client create error:")
		log.Error(err)
		return err
	}

	if err := client.Save(h); err != nil {
		log.Warn("save client error:")
		log.Error(err)
		return err
	}
	/*out, err := h.RunSSHCommand("df -h")
	if err != nil {
		log.Error(err)
		return err
	}

	fmt.Printf("Results of your disk space query:\n%s\n", out)*/

	fmt.Println("DONE")
	return nil
}

/*func createMachine(name, driver string) error {
	var cmd *exec.Cmd
	doToken := getToken()
	cmd = exec.Command("docker-machine", "create", "--driver", driver, "--digitalocean-access-token", doToken, name)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Cannot make the machine (%s): (%s)\n\n%s", name, err, out.String())
	}
	log.Warn(out.String())
	return nil
}*/

func initMachine(name string, services []string) error {
	util.DockerConnect(true, name) //verbose should take falg ?
	//TODO specify imgs -> marshal from file that has deps in it
	var images []string

	if len(services) == 1 && services[0] == "all" {
		images = []string{
			"quay.io/eris/tor-relay", //super lightweight
			//"quay.io/eris/base",
			//"quay.io/eris/keys",
			//"quay.io/eris/data",
			//"quay.io/eris/ipfs",
			//"quay.io/eris/erisdb",
			//"quay.io/eris/epm",
		}
	} else {
		images = services
	}

	//marshal service names into srv.Image
	if err := ini.GetTheImages(true, images); err != nil {
		return err
	}
	return nil
}

func regenerateCert(name string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("docker-machine", "regenerate-certs", "--force", name)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Cannot regenerate cert on machine (%s): (%s)\n\n%s", name, err, out.String())
	}
	log.Warn(out.String())
	return nil

}

func getToken() string {
	doToken := os.Getenv("DO_TOKEN")
	if doToken == "" {
		log.Warn("Digital Ocean Access Token not found.")
		log.Warn("Please set it with `$ export DO_TOKEN=secret_token`")
		os.Exit(1)
	}
	return doToken
}

func timeOutTicker(machines []string) {
	time.Sleep(maxTimeout)
	for _, m := range machines {
		fmt.Println(m)
	}
	os.Exit(1)
}

func setUpMachine(machine string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	file := filepath.Join(dir, script)
	cmd := exec.Command("docker-machine", "scp", file, fmt.Sprintf("%s:", machine))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Cannot scp into the machine (%s): (%s)\n\n%s", machine, err, out.String())
	}

	file = filepath.Base(file)
	cmd = exec.Command("docker-machine", "ssh", machine, fmt.Sprintf("sudo $HOME/%s", file))
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Cannot execute the command to change docker daemon on machine (%s): (%s)\n\n%s", machine, err, out.String())
	}
	return nil
}

func Load(name string) {
	//remDef, err := loaders.LoadRemoteDefinition(name)
	//if err != nil {
	//	log.WithField("=>", fmt.Sprintf("%v", err)).Warn("Loaders err:")
	//}
}
