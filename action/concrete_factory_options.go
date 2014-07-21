package action

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

type ConcreteFactoryOptions struct {
	StemcellsDir string

	// e.g diego-cpi
	LRPDomain string

	Agent bdcvm.AgentOptions
}

func (o ConcreteFactoryOptions) Validate() error {
	if o.StemcellsDir == "" {
		return bosherr.New("Must provide non-empty StemcellsDir")
	}

	err := o.Agent.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Agent configuration")
	}

	return nil
}
