/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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
