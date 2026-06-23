/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clone

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Value is the type alias for the clone block object.
type Value = types.Object

// Model represents the clone block's Terraform model.
type Model struct {
	VMID        types.Int64  `tfsdk:"vm_id"`
	NodeName    types.String `tfsdk:"node_name"`
	Full        types.Bool   `tfsdk:"full"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	Retries     types.Int64  `tfsdk:"retries"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vm_id":        types.Int64Type,
		"node_name":    types.StringType,
		"full":         types.BoolType,
		"datastore_id": types.StringType,
		"retries":      types.Int64Type,
	}
}

// NullValue returns a null clone block value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}
