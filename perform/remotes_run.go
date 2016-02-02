package perform

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	//"path/filepath"
	"sync"
	"time"

	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/util"

	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/commands/mcndirs"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/host"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/persist"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/drivers/digitalocean"
	//. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

var wg sync.WaitGroup
var maxTimeout = 15 * time.Minute

//will take a ops.Remotes or w/e
func CreateRemote(rem *definitions.RemoteDefinition) error {
	machines := rem.Machines

	failOut := make(chan bool, len(machines))
	go timeOutTicker(machines)

	wg.Add(len(machines))
	for _, m := range machines {
		go createMachines(m, rem.Driver, failOut)
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

	return nil
}

func ReprovisionRemote(rem *definitions.RemoteDefinition) error {
	api, err := util.RemotesConnect("")
	if err != nil {
		return err
	}

	hosts, hostsInError := persist.LoadHosts(api, rem.Machines)

	if len(hostsInError) > 0 {
		for _, err := range hostsInError {
			log.Error(err)
		}

		return fmt.Errorf("uh oh")
	}

	if len(hosts) == 0 {
		return fmt.Errorf("no hosts")
	}

	hNames := make([]string, len(hosts))
	for i, h := range hosts {
		log.WithField("=>", h.Name).Warn("host to prov")
		hNames[i] = h.Name
	}

	failOut := make(chan bool, len(hosts))
	go timeOutTicker(hNames)

	wg.Add(len(hosts))
	for _, h := range hosts {
		go provisionHosts(h, failOut)
	}

	go func(failOut chan bool) {
		if _, ok := <-failOut; ok {
			for _, h := range hosts {
				fmt.Println(h.Name)
			}
			os.Exit(1)
		}
	}(failOut)
	wg.Wait()

	for _, h := range hosts {
		log.WithField("=>", h.Name).Warn("Saving host to store")
		if err := api.Save(h); err != nil {
			return fmt.Errorf("Error saving host to store: %s", err)
		}
	}
	return nil
}

func RemoveRemote(rem *definitions.RemoteDefinition) error {
	machines := rem.Machines

	failOut := make(chan bool, len(machines))
	go timeOutTicker(machines)

	wg.Add(len(machines))
	for _, m := range machines {
		go removeAllMachines(m, failOut)
		//go createMachines(m, rem.Driver, failOut)
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
		log.WithField("=>", m).Warn("Machine eliminated:")
	}
	return nil
}

func createMachines(name, driver string, failOut chan<- bool) {
	defer wg.Done()
	if err := createMachine(name, driver); err != nil {
		failOut <- true
		//return
	}
}
func provisionHosts(h *host.Host, failOut chan<- bool) {
	defer wg.Done()
	log.Warn("provision host")
	if err := h.Provision(); err != nil {
		failOut <- true
		//return
	}
}

/*func provisionHost(h *host.Host) error {
	log.Warn("provision host")
	if err := h.Provision(); err != nil {
		return err
	}
	return nil
}*/

//All = local & cloud
func removeAllMachines(mach string, failOut chan<- bool) {
	defer wg.Done()
	if err := removeMachine(mach); err != nil {
		failOut <- true
		//return
	}
}

func removeMachine(mach string) error {
	api, err := util.RemotesConnect("")
	if err != nil {
		return err
	}

	if err := removeRemoteMachine(mach, api); err != nil {
		return err
	}

	if err := removeLocalMachine(mach, api); err != nil {
		return err
	}
	return nil
}

//TODO add region option
func createMachine(name, driverTyp string) error {

	api, err := util.RemotesConnect(name)
	if err != nil {
		return err
	}

	driver := digitalocean.NewDriver(name, mcndirs.GetBaseDir())

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

	//driver should be "digitalocean", set in remDef file
	h, err := api.NewHost("digitalocean", rawDriver)
	if err != nil {
		log.Warn("newHost error:")
		log.Error(err)
		return err
	}

	if err := api.Create(h); err != nil {
		log.Warn("client create error:")
		log.Error(err)
		return err
	}

	//maybe not needed?
	if err := api.Save(h); err != nil {
		log.Warn("save client error:")
		log.Error(err)
		return err
	}

	fmt.Println("DONE")
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

func removeRemoteMachine(hostName string, api libmachine.API) error {
	currentHost, loaderr := api.Load(hostName)
	if loaderr != nil {
		return loaderr
	}

	return currentHost.Driver.Remove()
}

func removeLocalMachine(hostName string, api libmachine.API) error {
	exist, _ := api.Exists(hostName)
	if !exist {
		return errors.New(hostName + " does not exist.")
	}
	return api.Remove(hostName)
}

////////////// TODO much later

func PullImagesToRemote(rem *definitions.RemoteDefinition) error {
	machines := rem.Machines

	log.Warn("Initializing docker images:")
	for _, m := range machines {
		log.WithField("=>", m).Warn("Initializing docker images for machine:")
		if err := initMachine(m, rem.Services); err != nil {
			log.Warn(fmt.Sprintf("some dumb error: %v\n", err))
		}
	}

	//can be better
	/*failOut := make(chan bool, len(machines))
	go timeOutTicker(machines)

	wg.Add(len(machines))
	for _, m := range machines {
		go initMachines(m, failOut)
	}

	go func(failOut chan bool) {
		if _, ok := <-failOut; ok {
			for _, m := range machines {
				fmt.Println(m)
			}
			os.Exit(1)
		}
	}(failOut)
	wg.Wait()*/

	return nil
}

func initMachine(name string, services []string) error {
	util.DockerConnect(true, name) //verbose should take falg ?
	//TODO specify imgs -> marshal from file that has deps in it
	var images []string

	if len(services) == 1 && services[0] == "all" {
		images = []string{
			"quay.io/eris/tor-relay", //super lightweight
		}
	} else {
		images = services
	}

	srv := definitions.BlankService()
	srv.Image = images[0]
	srv.Name = "tor"
	opts := definitions.BlankOperation()
	//marshal service names into srv.Image
	if err := DockerPull(srv, opts); err != nil {
		return err
	}
	return nil
}

/*func initMachines(name string, failOut chan<- bool) error {
	defer wg.Done()
	if err := initMachine(name); err != nil {
		failOut <- true
	}
	return nil
}*/
