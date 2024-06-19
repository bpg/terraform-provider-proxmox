/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents the CPU model.
type Model struct {
	DatastoreId types.String `tfsdk:"datastore_id"`
	Interface   types.String `tfsdk:"interface"`
	DNS         DNSValue     `tfsdk:"dns"`
}

type ModelDNS struct {
	Domain  types.String `tfsdk:"domain"`
	Servers types.List   `tfsdk:"servers"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"datastore_id": types.StringType,
		"interface":    types.StringType,
		"dns":          types.ObjectType{}.WithAttributeTypes(attributeTypesDNS()),
	}
}

func attributeTypesDNS() map[string]attr.Type {
	return map[string]attr.Type{
		"domain":  types.StringType,
		"servers": types.ListType{ElemType: types.StringType},
	}
}
