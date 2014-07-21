package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
)

type DeleteStemcell struct {
	stemcellFinder bdcstem.Finder
}

func NewDeleteStemcell(stemcellFinder bdcstem.Finder) DeleteStemcell {
	return DeleteStemcell{stemcellFinder: stemcellFinder}
}

func (a DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	stemcell, found, err := a.stemcellFinder.Find(string(stemcellCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding stemcell '%s'", stemcellCID)
	}

	if found {
		err := stemcell.Delete()
		if err != nil {
			return nil, bosherr.WrapError(err, "Deleting stemcell '%s'", stemcellCID)
		}
	}

	return nil, nil
}
