package subnets

import (
	"context"
)

type API interface {
	GetSubnets(ctx context.Context, vnetID string) ([]SubnetData, error)
	GetSubnet(ctx context.Context, vnetID string, id string) (*SubnetData, error)
	CreateSubnet(ctx context.Context, vnetID string, data *SubnetRequestData) error
	UpdateSubnet(ctx context.Context, vnetID string, data *SubnetRequestData) error
	DeleteSubnet(ctx context.Context, vnetID string, id string) error
}
