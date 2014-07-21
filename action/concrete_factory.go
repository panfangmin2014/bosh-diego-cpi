package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshcmd "github.com/cloudfoundry/bosh-agent/platform/commands"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	bdcbbs "github.com/cppforlife/bosh-diego-cpi/bbs"
	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

type concreteFactory struct {
	availableActions map[string]Action
}

func NewConcreteFactory(
	bbs bdcbbs.BBS,
	fs boshsys.FileSystem,
	uuidGen boshuuid.Generator,
	compressor boshcmd.Compressor,
	options ConcreteFactoryOptions,
	logger boshlog.Logger,
) concreteFactory {
	stemcellImporter := bdcstem.NewFSImporter(options.StemcellsDir, fs, uuidGen, compressor, logger)
	stemcellFinder := bdcstem.NewFSFinder(options.StemcellsDir, fs, logger)

	vmCreator := bdcvm.NewLRPCreator(uuidGen, bbs, options.LRPDomain, options.Agent, logger)
	vmFinder := bdcvm.NewLRPFinder(bbs, logger)

	return concreteFactory{
		availableActions: map[string]Action{
			// Stemcell management
			"create_stemcell": NewCreateStemcell(stemcellImporter),
			"delete_stemcell": NewDeleteStemcell(stemcellFinder),

			// VM management
			"create_vm":          NewCreateVM(stemcellFinder, vmCreator),
			"delete_vm":          NewDeleteVM(vmFinder),
			"has_vm":             NewHasVM(vmFinder),
			"reboot_vm":          NewRebootVM(),
			"set_vm_metadata":    NewSetVMMetadata(),
			"configure_networks": NewConfigureNetworks(),

			// Not implemented:
			//   create_disk
			//   delete_disk
			//   attach_disk
			//   detach_disk
			//   current_vm_id
			//   snapshot_disk
			//   delete_snapshot
			//   get_disks
			//   ping
		},
	}
}

func (f concreteFactory) Create(method string) (Action, error) {
	action, found := f.availableActions[method]
	if !found {
		return nil, bosherr.New("Could not create action with method '%s'", method)
	}

	return action, nil
}
