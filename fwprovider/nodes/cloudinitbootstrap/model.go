/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bootstrapModel struct {
	ID            types.String `tfsdk:"id"`
	VMID          types.Int64  `tfsdk:"vm_id"`
	DatastoreID   types.String `tfsdk:"datastore_id"`
	FileName      types.String `tfsdk:"file_name"`
	Interface     types.String `tfsdk:"interface"`
	UserData      types.String `tfsdk:"user_data"`
	MetaData      types.String `tfsdk:"meta_data"`
	NetworkConfig types.String `tfsdk:"network_config"`
	VendorData    types.String `tfsdk:"vendor_data"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	NodeName      types.String `tfsdk:"node_name"`
	SeedFileID    types.String `tfsdk:"seed_file_id"`
	IPv4Addresses types.Map    `tfsdk:"ipv4_addresses"`
}
