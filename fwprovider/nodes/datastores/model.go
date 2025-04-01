/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

type Model struct {
	NodeName types.String `tfsdk:"node_name"`
	Filters  *struct {
		ContentTypes stringset.Value `tfsdk:"content_types"`
		ID           types.String    `tfsdk:"id"`
		Target       types.String    `tfsdk:"target"`
	} `tfsdk:"filters"`
	Datastores []Datastore `tfsdk:"datastores"`
}

type Datastore struct {
	Active            types.Bool      `tfsdk:"active"`
	ContentTypes      stringset.Value `tfsdk:"content_types"`
	Enabled           types.Bool      `tfsdk:"enabled"`
	ID                types.String    `tfsdk:"id"`
	NodeName          types.String    `tfsdk:"node_name"`
	Shared            types.Bool      `tfsdk:"shared"`
	SpaceAvailable    types.Int64     `tfsdk:"space_available"`
	SpaceTotal        types.Int64     `tfsdk:"space_total"`
	SpaceUsed         types.Int64     `tfsdk:"space_used"`
	SpaceUsedFraction types.Float64   `tfsdk:"space_used_fraction"`
	Type              types.String    `tfsdk:"type"`
}
