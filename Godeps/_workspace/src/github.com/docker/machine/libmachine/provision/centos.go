package provision

import (
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/drivers"
)

func init() {
	Register("Centos", &RegisteredProvisioner{
		New: NewCentosProvisioner,
	})
}

func NewCentosProvisioner(d drivers.Driver) Provisioner {
	return &CentosProvisioner{
		NewRedHatProvisioner("centos", d),
	}
}

type CentosProvisioner struct {
	*RedHatProvisioner
}

func (provisioner *CentosProvisioner) String() string {
	return "centos"
}
