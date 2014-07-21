package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
)

type CreateStemcell struct {
	stemcellImporter bdcstem.Importer
}

type CreateStemcellCloudProps struct {
	RootFSPath string `json:"rootfs_path"`
	Stack      string `json:"stack"`
}

func NewCreateStemcell(stemcellImporter bdcstem.Importer) CreateStemcell {
	return CreateStemcell{stemcellImporter: stemcellImporter}
}

func (a CreateStemcell) Run(imagePath string, cloudProps CreateStemcellCloudProps) (StemcellCID, error) {
	props := bdcstem.Props{
		RootFSPath: cloudProps.RootFSPath,
		Stack:      cloudProps.Stack,
	}

	stemcell, err := a.stemcellImporter.ImportFromPath(imagePath, props)
	if err != nil {
		return "", bosherr.WrapError(err, "Importing stemcell from '%s'", imagePath)
	}

	return StemcellCID(stemcell.ID()), nil
}
