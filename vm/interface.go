package vm

import (
	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
)

type Creator interface {
	// Create takes an agent id and creates a VM with provided configuration
	Create(string, bdcstem.Stemcell, Props, Networks, Environment) (VM, error)
}

type Finder interface {
	Find(string) (VM, bool, error)
}

type Props struct {
	MemoryMB int
	DiskMB   int
}

type VM interface {
	ID() string

	Delete() error
}

type Environment map[string]interface{}
