package vm

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bdcbbs "github.com/cppforlife/bosh-diego-cpi/bbs"
)

type LRPVM struct {
	id     string
	bbs    bdcbbs.BBS
	logger boshlog.Logger
}

func NewLRPVM(
	id string,
	bbs bdcbbs.BBS,
	logger boshlog.Logger,
) LRPVM {
	return LRPVM{
		id:     id,
		bbs:    bbs,
		logger: logger,
	}
}

func (vm LRPVM) ID() string { return vm.id }

func (vm LRPVM) Delete() error {
	// todo if we kill desired lrp before waiting for all instances to die
	// #has_vm will return false when instances are potentially running
	err := vm.bbs.RemoveDesiredLRPByProcessGuid(vm.id)
	if err != nil {
		return bosherr.WrapError(err, "Removing LRP")
	}

	// fyi converger runs every 30 secs
	for i := 0; i < 91; i++ {
		lrps, err := vm.bbs.GetActualLRPsByProcessGuid(vm.id)
		if err != nil {
			if err == bdcbbs.LRPNotFound {
				return nil
			}

			return bosherr.WrapError(err, "Getting actual LRP after deleting desired LRP")
		}

		if len(lrps) == 0 {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return bosherr.WrapError(err, "Timed out waiting for actual LRP to be deleted")
}
