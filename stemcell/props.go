package stemcell

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type Props struct {
	RootFSPath string
	Stack      string
}

func NewPropsFromPath(path string, fs boshsys.FileSystem) (Props, error) {
	var props Props

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return props, bosherr.WrapError(err, "Reading props '%s'", path)
	}

	err = json.Unmarshal(bytes, &props)
	if err != nil {
		return props, bosherr.WrapError(err, "Unmarshalling props")
	}

	return props, nil
}

func WritePropsToPath(props Props, path string, fs boshsys.FileSystem) error {
	bytes, err := json.Marshal(props)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling props")
	}

	err = fs.WriteFile(path, bytes)
	if err != nil {
		return bosherr.WrapError(err, "Writing props '%s'", path)
	}

	return nil
}
