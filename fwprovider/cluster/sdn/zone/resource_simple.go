/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var (
	_ resource.ResourceWithConfigure   = &SimpleResource{}
	_ resource.ResourceWithImportState = &SimpleResource{}
)

type simpleModel struct {
	genericModel
}

type SimpleResource struct {
	*genericZoneResource
}

func NewSimpleResource() resource.Resource {
	return &SimpleResource{
		genericZoneResource: newGenericZoneResource(zoneResourceConfig{
			typeNameSuffix: "_sdn_zone_simple",
			zoneType:       zones.TypeSimple,
			modelFunc:      func() zoneModel { return &simpleModel{} },
		}).(*genericZoneResource),
	}
}

func (r *SimpleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simple Zone in Proxmox SDN.",
		MarkdownDescription: "Simple Zone in Proxmox SDN. It will create an isolated VNet bridge. " +
			"This bridge is not linked to a physical interface, and VM traffic is only local on each the node. " +
			"It can be used in NAT or routed setups.",
		Attributes: genericAttributesWith(nil),
	}
}
