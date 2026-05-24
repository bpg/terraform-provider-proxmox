/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

type nodeConfigModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Description types.String `tfsdk:"description"`
}

func (m *nodeConfigModel) toAPI() *nodes.ConfigUpdateRequestBody {
	return &nodes.ConfigUpdateRequestBody{
		Description: attribute.StringPtrFromValue(m.Description),
	}
}

func (m *nodeConfigModel) fromAPI(data *nodes.ConfigGetResponseData) {
	if data.Description != nil && *data.Description != "" {
		// PVE stores description as a comment in the config file and appends one trailing newline.
		trimmed := strings.TrimSuffix(*data.Description, "\n")
		m.Description = types.StringValue(trimmed)
	} else {
		m.Description = types.StringNull()
	}
}
