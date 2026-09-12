/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

// rollbackCreatedInterface deletes an interface whose create succeeded but whose read-back failed, so the
// staged change does not outlive the failed apply. Terraform has no state for it, and a retry would fail
// with "interface already exists".
func rollbackCreatedInterface(ctx context.Context, client *nodes.Client, iface string, diags *diag.Diagnostics) {
	err := client.DeleteNetworkInterface(ctx, iface)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		diags.AddError(fmt.Sprintf("Unable to Roll Back Network Interface %q", iface), err.Error())
	}
}
