package main

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/machine/drivers/digitalocean"
	"github.com/docker/machine/libmachine"

	log "github.com/Sirupsen/logrus"
	. "github.com/eris-ltd/common/go/common"
)

var wg sync.WaitGroup
var maxTimeout = 15 * time.Minute

func CreateRemote(machines []string) error {

	failOut := make(chan bool, len(machines))
	go timeOutTicker(machines)

	wg.Add(len(machines))
	for _, m := range machines {
		go createMachines(m, "digitalocean", failOut)
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

func createMachines(name, driver string, failOut chan<- bool) {
	defer wg.Done()
	if err := createMachine(name, driver); err != nil {
		failOut <- true
		//return
	}
}

func createMachine(name, driverTyp string) error {
	api, err := RemotesConnect(name)
	if err != nil {
		return err
	}

	driver := digitalocean.NewDriver(name, GetBaseDir())

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

	if err := api.Save(h); err != nil {
		log.Warn("save client error:")
		log.Error(err)
		return err
	}

	return nil
}

func RemotesConnect(remName string) (api libmachine.API, err error) {

	api = libmachine.NewClient(GetBaseDir(), GetMachineCertDir())
	defer api.Close()

	return api, nil
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

// modified from
// github.com/docker/machine/commands/mcndirs/util.go

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
