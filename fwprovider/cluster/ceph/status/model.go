/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package status

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	clusterceph "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ceph"
)

// model is the Terraform-side representation of the Ceph status data source.
// `NodeName` is the only input attribute - the rest are populated from the
// API response in fromAPI.
type model struct {
	ID           types.String `tfsdk:"id"`
	NodeName     types.String `tfsdk:"node_name"`
	FSID         types.String `tfsdk:"fsid"`
	HealthStatus types.String `tfsdk:"health_status"`
	QuorumNames  types.List   `tfsdk:"quorum_names"`
}

// fromAPI populates the Computed fields of the model from the typed
// API response payload. `NodeName` is left untouched - it is an input
// attribute set by the caller. `id` mirrors `fsid`, which is the stable
// Ceph cluster identifier returned by both the cluster and node endpoints.
func (m *model) fromAPI(ctx context.Context, data *clusterceph.StatusResponseData) diag.Diagnostics {
	m.ID = types.StringValue(data.FSID)
	m.FSID = types.StringValue(data.FSID)
	m.HealthStatus = types.StringValue(data.Health.Status)

	// A null list is invalid for a Computed attribute after Read; coerce a nil
	// API response to an empty list so downstream configs see a known value.
	if data.QuorumNames == nil {
		data.QuorumNames = []string{}
	}

	quorum, diags := types.ListValueFrom(ctx, types.StringType, data.QuorumNames)
	if diags.HasError() {
		return diags
	}

	m.QuorumNames = quorum

	return diags
}
