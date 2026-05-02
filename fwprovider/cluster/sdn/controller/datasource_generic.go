/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controller

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

type controllerDataSourceConfig struct {
	typeNameSuffix string
	controllerType string
	modelFunc      func() controllerModel
}

type genericControllerDataSource struct {
	client *cluster.Client
	config controllerDataSourceConfig
}

func newGenericControllerDataSource(cfg controllerDataSourceConfig) *genericControllerDataSource {
	return &genericControllerDataSource{config: cfg}
}

func (d *genericControllerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox" + d.config.typeNameSuffix
}

func (d *genericControllerDataSource) Configure(
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
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN controller",
			Required:    true,
		},
		"digest": schema.StringAttribute{
			Description: "Digest of the controller section.",
			Computed:    true,
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

func (d *genericControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := d.config.modelFunc()
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.getID()

	client := d.client.SDNControllers()

	controller, err := client.GetController(ctx, id)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Controller Not Found",
				fmt.Sprintf("SDN controller with ID '%s' was not found", state.getID()),
			)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read SDN Controller %q", id),
			err.Error(),
		)

		return
	}

	// Verify the controller type matches what this datasource expects
	if controller.Type != nil && *controller.Type != d.config.controllerType {
		resp.Diagnostics.AddError(
			"SDN Controller Type Mismatch",
			fmt.Sprintf(
				"Expected controller type '%s' but found '%s' for id '%s'",
				d.config.controllerType,
				*controller.Type,
				id,
			),
		)

		return
	}

	readModel := d.config.modelFunc()

	readModel.fromAPIForDatasource(id, controller, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
