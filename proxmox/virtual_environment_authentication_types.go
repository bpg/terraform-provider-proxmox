/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// VirtualEnvironmentAuthenticationResponseBody contains the body from an authentication response.
type VirtualEnvironmentAuthenticationResponseBody struct {
	Data *VirtualEnvironmentAuthenticationResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAuthenticationResponseCapabilities contains the supported capabilities for a session.
type VirtualEnvironmentAuthenticationResponseCapabilities struct {
	Access     *types.CustomPrivileges `json:"access,omitempty"`
	Datacenter *types.CustomPrivileges `json:"dc,omitempty"`
	Nodes      *types.CustomPrivileges `json:"nodes,omitempty"`
	Storage    *types.CustomPrivileges `json:"storage,omitempty"`
	VMs        *types.CustomPrivileges `json:"vms,omitempty"`
}

// VirtualEnvironmentAuthenticationResponseData contains the data from an authentication response.
type VirtualEnvironmentAuthenticationResponseData struct {
	ClusterName         *string                                               `json:"clustername,omitempty"`
	CSRFPreventionToken *string                                               `json:"CSRFPreventionToken,omitempty"`
	Capabilities        *VirtualEnvironmentAuthenticationResponseCapabilities `json:"cap,omitempty"`
	Ticket              *string                                               `json:"ticket,omitempty"`
	Username            string                                                `json:"username"`
}
