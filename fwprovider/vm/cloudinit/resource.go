/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// NewValue returns a new Value with the given CPU settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, vmID int, diags *diag.Diagnostics) *Model {
	ci := Model{}

	devices := config.CustomStorageDevices.Filter(func(device *vms.CustomStorageDevice) bool {
		return device.IsCloudInitDrive(vmID)
	})

	if len(devices) != 1 {
		return nil
	}

	for iface, device := range devices {
		ci.Interface = types.StringValue(iface)
		ci.DatastoreID = types.StringValue(device.GetDatastoreID())

		dns := ModelDNS{}
		dns.Domain = types.StringPointerValue(config.CloudInitDNSDomain)

		if config.CloudInitDNSServer != nil && strings.Trim(*config.CloudInitDNSServer, " ") != "" {
			dnsServers := strings.Split(*config.CloudInitDNSServer, " ")
			servers, d := types.ListValueFrom(ctx, customtypes.IPAddrType{}, dnsServers)
			diags.Append(d...)

			dns.Servers = servers
		} else {
			dns.Servers = types.ListNull(customtypes.IPAddrType{})
		}

		if !reflect.DeepEqual(dns, ModelDNS{}) {
			ci.DNS = &dns
		}

		return &ci
	}

	return nil
}

// FillCreateBody fills the CreateRequestBody with the Cloud-Init settings from the Value.
func FillCreateBody(ctx context.Context, plan *Model, body *vms.CreateRequestBody) {
	if plan == nil {
		return
	}

	ci := vms.CustomCloudInitConfig{}

	if plan.DNS != nil {
		dns := *plan.DNS

		if !dns.Domain.IsUnknown() {
			ci.SearchDomain = dns.Domain.ValueStringPointer()
		}

		if !dns.Servers.IsUnknown() {
			var servers []string

			dns.Servers.ElementsAs(ctx, &servers, false)

			ci.Nameserver = ptr.Ptr(strings.Join(servers, " "))
		}
	}

	body.CloudInitConfig = &ci

	device := vms.CustomStorageDevice{
		Enabled:    true,
		FileVolume: fmt.Sprintf("%s:cloudinit", plan.DatastoreID.ValueString()),
		Media:      ptr.Ptr("cdrom"),
	}

	body.AddCustomStorageDevice(plan.Interface.ValueString(), device)
}

// FillUpdateBody fills the UpdateRequestBody with the Cloud-Init settings from the Value.
func FillUpdateBody(
	ctx context.Context,
	plan, state *Model,
	updateBody *vms.UpdateRequestBody,
	isClone bool,
	diags *diag.Diagnostics,
) {
	if plan == nil || reflect.DeepEqual(plan, state) {
		return
	}

	del := func(field ...string) {
		updateBody.Delete = append(updateBody.Delete, field...)
	}

	// TODO: migrate cloud init to another datastore

	if !reflect.DeepEqual(plan.DNS, state.DNS) {
		if plan.DNS == nil && state.DNS != nil && !isClone {
			del("searchdomain", "nameserver")
		} else if plan.DNS != nil {
			ci := vms.CustomCloudInitConfig{}

			planDNS := plan.DNS
			stateDNS := state.DNS

			if !planDNS.Domain.Equal(stateDNS.Domain) {
				if attribute.ShouldBeRemoved(planDNS.Domain, stateDNS.Domain, isClone) {
					del("searchdomain")
				} else if attribute.IsDefined(planDNS.Domain) {
					ci.SearchDomain = planDNS.Domain.ValueStringPointer()
				}
			}

			if !planDNS.Servers.Equal(stateDNS.Servers) {
				if attribute.ShouldBeRemoved(planDNS.Servers, stateDNS.Servers, isClone) {
					del("nameserver")
				} else if attribute.IsDefined(planDNS.Servers) {
					// TODO: duplicates code from FillCreateBody
					var servers []string

					planDNS.Servers.ElementsAs(ctx, &servers, false)

					//// special case for the servers list, if we want to remove them during update
					//if len(servers) == 0 {
					//	del("nameserver")
					//} else {
					ci.Nameserver = ptr.Ptr(strings.Join(servers, " "))
					//}
				}
			}

			updateBody.CloudInitConfig = &ci
		}
	}
}
