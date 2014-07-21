package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

type HasVM struct {
	vmFinder bdcvm.Finder
}

func NewHasVM(vmFinder bdcvm.Finder) HasVM {
	return HasVM{vmFinder: vmFinder}
}

func (a HasVM) Run(vmCID VMCID) (bool, error) {
	_, found, err := a.vmFinder.Find(string(vmCID))
	if err != nil {
		return false, bosherr.WrapError(err, "Finding VM '%s'", vmCID)
	}

	return found, nil
}
