package zones

import (
	"context"
)

type API interface {
	GetZones(ctx context.Context) ([]ZoneData, error)
	GetZone(ctx context.Context, id string) (*ZoneData, error)
	CreateZone(ctx context.Context, req *ZoneRequestData) error
	UpdateZone(ctx context.Context, req *ZoneRequestData) error
	DeleteZone(ctx context.Context, id string) error
}
