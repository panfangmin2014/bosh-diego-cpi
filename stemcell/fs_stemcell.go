package stemcell

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

const fsStemcellLogTag = "FSStemcell"

type FSStemcell struct {
	id      string
	dirPath string
	props   Props

	fs     boshsys.FileSystem
	logger boshlog.Logger
}

func NewFSStemcell(
	id string,
	dirPath string,
	props Props,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) FSStemcell {
	return FSStemcell{id: id, dirPath: dirPath, props: props, fs: fs, logger: logger}
}

func (s FSStemcell) ID() string { return s.id }

func (s FSStemcell) RootFSPath() string { return s.props.RootFSPath }
func (s FSStemcell) Stack() string      { return s.props.Stack }

func (s FSStemcell) Delete() error {
	s.logger.Debug(fsStemcellLogTag, "Deleting stemcell '%s'", s.id)

	err := s.fs.RemoveAll(s.dirPath)
	if err != nil {
		return bosherr.WrapError(err, "Deleting stemcell directory '%s'", s.dirPath)
	}

	return nil
}
