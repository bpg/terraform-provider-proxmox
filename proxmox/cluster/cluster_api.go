package cluster

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type API struct {
	types.Client
}

func (a *API) Firewall() *firewall.API {
	return &firewall.API{Client: a}
}
