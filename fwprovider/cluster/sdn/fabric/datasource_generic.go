/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type fabricDataSourceConfig struct {
	typeNameSuffix string
	fabricProtocol string
	modelFunc      func() fabricModel
}

type genericFabricDataSource struct {
	client *fabrics.Client
	config fabricDataSourceConfig
}

func newGenericFabricDataSource(cfg fabricDataSourceConfig) *genericFabricDataSource {
	return &genericFabricDataSource{config: cfg}
}

func (d *genericFabricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.config.typeNameSuffix
}

func (d *genericFabricDataSource) Configure(
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

	d.client = cfg.Client.Cluster().SDNFabrics(d.config.fabricProtocol)
}

func genericDataSourceAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric.",
			Required:    true,
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

func (d *genericFabricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := d.config.modelFunc()
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fabric, err := d.client.GetFabricWithParams(ctx, state.getID(), &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Fabric Not Found",
				fmt.Sprintf("SDN fabric with ID '%s' was not found", state.getID()),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read SDN Fabric",
			err.Error(),
		)

		return
	}

	// Verify the fabric protocol matches what this datasource expects
	if fabric.Protocol != nil && *fabric.Protocol != d.config.fabricProtocol {
		resp.Diagnostics.AddError(
			"SDN Fabric Protocol Mismatch",
			fmt.Sprintf(
				"Expected fabric protocol '%s' but found '%s' for fabric '%s'",
				d.config.fabricProtocol,
				*fabric.Protocol,
				fabric.ID,
			),
		)

		return
	}

	readModel := d.config.modelFunc()

	readModel.fromAPI(fabric.ID, fabric, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
