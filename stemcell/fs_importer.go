package stemcell

import (
	"os"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshcmd "github.com/cloudfoundry/bosh-agent/platform/commands"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"
)

const fsImporterLogTag = "FSImporter"

type FSImporter struct {
	dirPath string

	fs         boshsys.FileSystem
	uuidGen    boshuuid.Generator
	compressor boshcmd.Compressor

	logger boshlog.Logger
}

func NewFSImporter(
	dirPath string,
	fs boshsys.FileSystem,
	uuidGen boshuuid.Generator,
	compressor boshcmd.Compressor,
	logger boshlog.Logger,
) FSImporter {
	return FSImporter{
		dirPath: dirPath,

		fs:         fs,
		uuidGen:    uuidGen,
		compressor: compressor,

		logger: logger,
	}
}

func (i FSImporter) ImportFromPath(imagePath string, props Props) (Stemcell, error) {
	i.logger.Debug(fsImporterLogTag, "Importing stemcell from path '%s'", imagePath)

	id, err := i.uuidGen.Generate()
	if err != nil {
		return nil, bosherr.WrapError(err, "Generating stemcell id")
	}

	stemcellPath := filepath.Join(i.dirPath, id)

	stemcellImagePath := filepath.Join(stemcellPath, "image")

	err = i.fs.MkdirAll(stemcellImagePath, os.FileMode(0755))
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating stemcell image directory '%s'", stemcellImagePath)
	}

	err = i.compressor.DecompressFileToDir(imagePath, stemcellImagePath)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unpacking stemcell '%s' to '%s'", imagePath, stemcellImagePath)
	}

	stemcellPropsPath := filepath.Join(stemcellPath, "props")

	err = WritePropsToPath(props, stemcellPropsPath, i.fs)
	if err != nil {
		removeAllErr := i.fs.RemoveAll(stemcellPath)
		if removeAllErr != nil {
			i.logger.Error(fsImporterLogTag, "Failed to clean up stemcell from path '%s'", stemcellPath)
		}

		return nil, bosherr.WrapError(err, "Writing props '%v' to '%s'", props, stemcellPropsPath)
	}

	i.logger.Debug(fsImporterLogTag, "Imported stemcell from path '%s'", imagePath)

	return NewFSStemcell(id, stemcellImagePath, props, i.fs, i.logger), nil
}
