package host

import "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/drivers"

type V2 struct {
	ConfigVersion int
	Driver        drivers.Driver
	DriverName    string
	HostOptions   *Options
	Name          string
}
