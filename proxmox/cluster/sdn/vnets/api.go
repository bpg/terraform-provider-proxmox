package vnets

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

type API interface {
	GetVnets(ctx context.Context) ([]VnetData, error)
	GetVnet(ctx context.Context, id string) (*VnetData, error)
	CreateVnet(ctx context.Context, req *VnetRequestData) error
	UpdateVnet(ctx context.Context, req *VnetRequestData) error
	DeleteVnet(ctx context.Context, id string) error
	GetParentZone(ctx context.Context, zoneId string) (*zones.ZoneData, error)
}
