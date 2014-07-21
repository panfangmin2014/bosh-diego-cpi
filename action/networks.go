package action

import (
	bdcvm "github.com/cppforlife/bosh-diego-cpi/vm"
)

type Networks map[string]Network

type Network struct {
	Type string `json:"type"`

	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`

	DNS     []string `json:"dns"`
	Default []string `json:"default"`

	MAC string `json:"mac"`

	CloudProperties map[string]interface{} `json:"cloud_properties"`
}

func (ns Networks) AsVMNetworks() bdcvm.Networks {
	networks := bdcvm.Networks{}

	for netName, network := range ns {
		networks[netName] = bdcvm.Network{
			Type: network.Type,

			IP:      network.IP,
			Netmask: network.Netmask,
			Gateway: network.Gateway,

			DNS:     network.DNS,
			Default: network.Default,

			CloudProperties: network.CloudProperties,
		}
	}

	return networks
}
