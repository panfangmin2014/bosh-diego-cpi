package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

type CreateVM struct {
	stemcellFinder bdcstem.Finder
	vmCreator      bdcvm.Creator
}

type CreateVMCloudProps struct {
	MemoryMB int `json:"memory_mb"`
	DiskMB   int `json:"disk_mb"`
}

type Environment map[string]interface{}

func NewCreateVM(stemcellFinder bdcstem.Finder, vmCreator bdcvm.Creator) CreateVM {
	return CreateVM{
		stemcellFinder: stemcellFinder,
		vmCreator:      vmCreator,
	}
}

func (a CreateVM) Run(agentID string, stemcellCID StemcellCID, cloudProps CreateVMCloudProps, networks Networks, _ []DiskCID, env Environment) (VMCID, error) {
	stemcell, found, err := a.stemcellFinder.Find(string(stemcellCID))
	if err != nil {
		return "", bosherr.WrapError(err, "Finding stemcell '%s'", stemcellCID)
	}

	if !found {
		return "", bosherr.New("Expected to find stemcell '%s'", stemcellCID)
	}

	props := bdcvm.Props{
		MemoryMB: cloudProps.MemoryMB,
		DiskMB:   cloudProps.DiskMB,
	}

	vmNetworks := networks.AsVMNetworks()

	vmEnv := bdcvm.Environment(env)

	vm, err := a.vmCreator.Create(agentID, stemcell, props, vmNetworks, vmEnv)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating VM with agent ID '%s'", agentID)
	}

	return VMCID(vm.ID()), nil
}
