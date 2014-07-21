package vm_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStemcell(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vm Suite")
}

type NonJSONMarshable struct{}

func (m NonJSONMarshable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("fake-marshal-err")
}
