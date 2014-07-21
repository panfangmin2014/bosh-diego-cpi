package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bdcaction "github.com/cppforlife/bosh-diego-cpi/action"
)

type Config struct {
	ETCD ETCDConfig

	Actions bdcaction.ConcreteFactoryOptions
}

type ETCDConfig struct {
	// e.g. [127.0.0.1:4001, 127.0.0.2:4001]
	ConnectAddresses []string
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapError(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) Validate() error {
	err := c.ETCD.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating ETCD configuration")
	}

	err = c.Actions.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Actions configuration")
	}

	return nil
}

func (c ETCDConfig) Validate() error {
	if len(c.ConnectAddresses) == 0 {
		return bosherr.New("Must provide non-empty ConnectAddresses")
	}

	return nil
}
