package stemcell

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type FSFinder struct {
	dirPath string

	fs     boshsys.FileSystem
	logger boshlog.Logger
}

func NewFSFinder(dirPath string, fs boshsys.FileSystem, logger boshlog.Logger) FSFinder {
	return FSFinder{dirPath: dirPath, fs: fs, logger: logger}
}

func (f FSFinder) Find(id string) (Stemcell, bool, error) {
	stemcellPath := filepath.Join(f.dirPath, id)

	if f.fs.FileExists(stemcellPath) {
		stemcellImagePath := filepath.Join(stemcellPath, "image")
		stemcellPropsPath := filepath.Join(stemcellPath, "props")

		props, err := NewPropsFromPath(stemcellPropsPath, f.fs)
		if err != nil {
			return nil, true, bosherr.WrapError(err, "Loading props")
		}

		return NewFSStemcell(id, stemcellImagePath, props, f.fs, f.logger), true, nil
	}

	return nil, false, nil
}
