package vm

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bdcbbs "github.com/cppforlife/bosh-diego-cpi/bbs"
)

const lrpFinderLogTag = "LRPFinder"

type LRPFinder struct {
	bbs    bdcbbs.BBS
	logger boshlog.Logger
}

func NewLRPFinder(bbs bdcbbs.BBS, logger boshlog.Logger) LRPFinder {
	return LRPFinder{bbs: bbs, logger: logger}
}

func (f LRPFinder) Find(id string) (VM, bool, error) {
	f.logger.Debug(lrpFinderLogTag, "Finding LRP with ID '%s'", id)

	lrp, err := f.bbs.GetDesiredLRPByProcessGuid(id)
	if err != nil {
		if err == bdcbbs.LRPNotFound { //storetcd.ErrorKeyNotFound
			f.logger.Debug(lrpFinderLogTag, "Did not find LRP with ID '%s'", id)
			return nil, false, nil
		}

		return nil, false, bosherr.WrapError(err, "Finding LRP by ID")
	}

	f.logger.Debug(lrpFinderLogTag, "Found LRP with ID '%s'", id)

	return NewLRPVM(lrp.ProcessGuid, f.bbs, f.logger), true, nil
}
