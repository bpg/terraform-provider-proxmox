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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for CPU settings.
type Value = types.Object

type DNSValue = types.Object

// NewValue returns a new Value with the given CPU settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, vmID int, diags *diag.Diagnostics) Value {
	ci := Model{}

	devices := config.CustomStorageDevices.Filter(func(device *vms.CustomStorageDevice) bool {
		return device.IsCloudInitDrive(vmID)
	})

	if len(devices) != 1 {
		types.ObjectNull(attributeTypes())
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
			dnsObj, d := types.ObjectValueFrom(ctx, attributeTypesDNS(), dns)
			diags.Append(d...)

			ci.DNS = dnsObj
		}

		obj, d := types.ObjectValueFrom(ctx, attributeTypes(), ci)
		diags.Append(d...)

		return obj
	}

	return types.ObjectNull(attributeTypes())
}

// FillCreateBody fills the CreateRequestBody with the Cloud-Init settings from the Value.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	var plan Model

	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if d.HasError() {
		return
	}

	ci := vms.CustomCloudInitConfig{}

	// TODO: should we check for !null?
	if !plan.DNS.IsUnknown() {
		var dns ModelDNS

		plan.DNS.As(ctx, &dns, basetypes.ObjectAsOptions{})

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
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	isClone bool,
	diags *diag.Diagnostics,
) {
	var plan, state Model

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	d = stateValue.As(ctx, &state, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	del := func(field ...string) {
		updateBody.Delete = append(updateBody.Delete, field...)
	}

	// TODO: migrate cloud init to another datastore

	if !reflect.DeepEqual(plan.DNS, state.DNS) {
		if attribute.ShouldBeRemoved(plan.DNS, state.DNS, isClone) {
			del("searchdomain", "nameserver")
		} else if attribute.IsDefined(plan.DNS) {
			ci := vms.CustomCloudInitConfig{}

			var planDNS, stateDNS ModelDNS
			d = plan.DNS.As(ctx, &planDNS, basetypes.ObjectAsOptions{})
			diags.Append(d...)
			d = state.DNS.As(ctx, &stateDNS, basetypes.ObjectAsOptions{})
			diags.Append(d...)

			if diags.HasError() {
				return
			}

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
					var servers []string

					planDNS.Servers.ElementsAs(ctx, &servers, false)

					ci.Nameserver = ptr.Ptr(strings.Join(servers, " "))
				}
			}

			updateBody.CloudInitConfig = &ci
		}
	}
}
