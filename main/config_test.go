package main_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"

	bdcaction "github.com/cppforlife/bosh-diego-cpi/action"
	. "github.com/cppforlife/bosh-diego-cpi/main"
	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

var validConfig = Config{
	ETCD:    validETCDConfig,
	Actions: validActionsOptions,
}

var validETCDConfig = ETCDConfig{
	ConnectAddresses: []string{"fake-address"},
}

var validActionsOptions = bdcaction.ConcreteFactoryOptions{
	StemcellsDir: "/tmp/stemcells",

	Agent: bdcvm.AgentOptions{
		Mbus: "fake-mbus",
		NTP:  []string{},

		Blobstore: bdcvm.BlobstoreOptions{
			Type: "fake-blobstore-type",
		},
	},
}

var _ = Describe("NewConfigFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("returns error if config is not valid", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Validating config"))
	})

	It("returns error if file contains invalid json", func() {
		err := fs.WriteFileString("/config.json", "-")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Unmarshalling config"))
	})

	It("returns error if file cannot be read", func() {
		err := fs.WriteFileString("/config.json", "{}")
		Expect(err).ToNot(HaveOccurred())

		fs.ReadFileError = errors.New("fake-read-err")

		_, err = NewConfigFromPath("/config.json", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-read-err"))
	})
})

var _ = Describe("Config", func() {
	var (
		config Config
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all warden and agent sections are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if warden section is not valid", func() {
			config.ETCD.ConnectAddresses = []string{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating ETCD configuration"))
		})

		It("returns error if actions section is not valid", func() {
			config.Actions.StemcellsDir = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Actions configuration"))
		})
	})
})

var _ = Describe("ETCDConfig", func() {
	var (
		config ETCDConfig
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validETCDConfig
		})

		It("does not return error if all fields are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if ConnectAddress is empty", func() {
			config.ConnectAddresses = []string{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide non-empty ConnectAddresses"))
		})
	})
})
