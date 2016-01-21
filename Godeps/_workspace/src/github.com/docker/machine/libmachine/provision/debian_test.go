package provision

import (
	"testing"

	"github.com/docker/machine/drivers/fakedriver"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/auth"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/engine"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/provision/provisiontest"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/swarm"
)

func TestDebianDefaultStorageDriver(t *testing.T) {
	p := NewDebianProvisioner(&fakedriver.Driver{}).(*DebianProvisioner)
	p.SSHCommander = provisiontest.NewFakeSSHCommander(provisiontest.FakeSSHCommanderOptions{})
	p.Provision(swarm.Options{}, auth.Options{}, engine.Options{})
	if p.EngineOptions.StorageDriver != "aufs" {
		t.Fatal("Default storage driver should be aufs")
	}
}
