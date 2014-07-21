package vm

import (
	"encoding/json"
	"fmt"

	runmodels "github.com/cloudfoundry-incubator/runtime-schema/models"
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	bdcbbs "github.com/cppforlife/bosh-diego-cpi/bbs"
	bdcstem "github.com/cppforlife/bosh-diego-cpi/stemcell"
)

const lrpCreatorLogTag = "LRPCreator"

type LRPCreator struct {
	uuidGen      boshuuid.Generator
	bbs          bdcbbs.BBS
	lrpDomain    string
	agentOptions AgentOptions
	logger       boshlog.Logger
}

func NewLRPCreator(
	uuidGen boshuuid.Generator,
	bbs bdcbbs.BBS,
	lrpDomain string,
	agentOptions AgentOptions,
	logger boshlog.Logger,
) LRPCreator {
	return LRPCreator{
		uuidGen:      uuidGen,
		bbs:          bbs,
		lrpDomain:    lrpDomain,
		agentOptions: agentOptions,
		logger:       logger,
	}
}

func (c LRPCreator) Create(agentID string, stemcell bdcstem.Stemcell, props Props, networks Networks, env Environment) (VM, error) {
	id, err := c.uuidGen.Generate()
	if err != nil {
		return LRPVM{}, bosherr.WrapError(err, "Generating VM id")
	}

	err = c.ensureSingleDynamicNetwork(networks)
	if err != nil {
		return LRPVM{}, err
	}

	agentEnv := NewAgentEnvForVM(agentID, id, networks, env, c.agentOptions)

	agentEnvJSONBytes, err := json.Marshal(agentEnv)
	if err != nil {
		return LRPVM{}, bosherr.WrapError(err, "Marshalling agent env")
	}

	lrp := c.buildLRP(id, stemcell, props, agentEnvJSONBytes)

	c.logger.Debug(lrpCreatorLogTag, "Creating LRP with spec %#v", lrp)

	err = c.bbs.DesireLRP(lrp)
	if err != nil {
		return LRPVM{}, bosherr.WrapError(err, "Creating container")
	}

	vm := NewLRPVM(id, c.bbs, c.logger)

	return vm, nil
}

func (c LRPCreator) ensureSingleDynamicNetwork(networks Networks) error {
	var network Network

	switch len(networks) {
	case 0:
		return bosherr.New("Expected exactly one network; received zero")
	case 1:
		network = networks.First()
	default:
		return bosherr.New("Expected exactly one network; received multiple")
	}

	if !network.IsDynamic() {
		return bosherr.New("Expected network's type to be 'dynamic'")
	}

	return nil
}

func (c LRPCreator) buildLRP(id string, stemcell bdcstem.Stemcell, props Props, agentEnvJSONBytes []byte) runmodels.DesiredLRP {
	// todo generalize agent env configuration?
	writeEnvAction := runmodels.RunAction{
		Path: "bash",
		Args: []string{"-c", "echo $BOSH_SETTINGS_JSON > /var/vcap/bosh/warden-cpi-agent-env.json"},
		Env: []runmodels.EnvironmentVariable{
			runmodels.EnvironmentVariable{
				Name:  "BOSH_SETTINGS_JSON",
				Value: string(agentEnvJSONBytes),
			},
		},
	}

	runRunitAction := runmodels.RunAction{
		Path: "/usr/sbin/runsvdir-start",
	}

	monitorAction := runmodels.MonitorAction{
		// do nothing for monitoring
		Action: runmodels.ExecutorAction{Action: runmodels.ParallelAction{}},

		// todo do not specify a hook at all
		HealthyHook: runmodels.HealthRequest{
			Method: "PUT",
			URL:    fmt.Sprintf("http://127.0.0.1:20515/lrp_running/%s/PLACEHOLDER_INSTANCE_INDEX/PLACEHOLDER_INSTANCE_GUID", id),
		},
	}

	lrp := runmodels.DesiredLRP{
		ProcessGuid: id,
		Instances:   1,
		Domain:      c.lrpDomain,

		RootFSPath: stemcell.RootFSPath(),
		Stack:      stemcell.Stack(),

		MemoryMB: props.MemoryMB,
		DiskMB:   props.DiskMB,

		Actions: []runmodels.ExecutorAction{
			runmodels.ExecutorAction{Action: writeEnvAction},
			runmodels.Parallel(
				runmodels.ExecutorAction{Action: runRunitAction},
				runmodels.ExecutorAction{Action: monitorAction},
			),
		},

		// Routes: []runmodels.Route{}
		// Ports: []runmodels.PortMapping{},
		// Log: runmodels.LogConfig{},
	}

	return lrp
}
