/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents the cloud-init model.
type Model struct {
	DatastoreID types.String `tfsdk:"datastore_id"`
	Interface   types.String `tfsdk:"interface"`
	DNS         *ModelDNS    `tfsdk:"dns"`
}

type ModelDNS struct {
	Domain  types.String `tfsdk:"domain"`
	Servers types.List   `tfsdk:"servers"`
}
