/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type fabricNodeDataSourceConfig struct {
	typeNameSuffix string
	fabricProtocol string
	modelFunc      func() fabricNodeModel
}

type genericFabricNodeDataSource struct {
	client *cluster.Client
	config fabricNodeDataSourceConfig
}

func newGenericFabricNodeDataSource(cfg fabricNodeDataSourceConfig) *genericFabricNodeDataSource {
	return &genericFabricNodeDataSource{config: cfg}
}

func (d *genericFabricNodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.config.typeNameSuffix
}

func (d *genericFabricNodeDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf(
				"Expected config.DataSource, got: %T",
				req.ProviderData,
			),
		)

		return
	}

	d.client = cfg.Client.Cluster()
}

func genericDataSourceAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"fabric_id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric.",
			Required:    true,
		},
		"node_id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric node.",
			Required:    true,
		},
		"interface_names": schema.SetAttribute{
			Description: "Set of interface names associated with the fabric node.",
			Computed:    true,
			ElementType: types.StringType,
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

func (d *genericFabricNodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := d.config.modelFunc()
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	id := state.getID()
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected SDN Fabric Node ID Format",
			fmt.Sprintf("Expected SDN Fabric Node ID to be in the format <fabric_id>/<node_id>, got: %s", id),
		)

		return
	}
	fabricID := parts[0]
	nodeID := parts[1]

	client := d.client.SDNFabricNodes(fabricID, d.config.fabricProtocol)
	fabricNode, err := client.GetFabricNodeWithParams(ctx, nodeID, &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Fabric Node Not Found",
				fmt.Sprintf("SDN fabric node with ID '%s' was not found", state.getID()),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read SDN Fabric Node",
			err.Error(),
		)

		return
	}

	// Verify the fabric protocol matches what this datasource expects
	if fabricNode.Protocol != nil && *fabricNode.Protocol != d.config.fabricProtocol {
		resp.Diagnostics.AddError(
			"SDN Fabric Protocol Mismatch",
			fmt.Sprintf(
				"Expected fabric protocol '%s' but found '%s' for fabric '%s' and node '%s'",
				d.config.fabricProtocol,
				*fabricNode.Protocol,
				fabricNode.FabricID,
				fabricNode.NodeID,
			),
		)

		return
	}

	readModel := d.config.modelFunc()

	readModel.fromAPI(id, fabricNode, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
