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

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value represents the type for CPU settings.
type Value = types.Object

type DNSValue = types.Object

// NewValue returns a new Value with the given CPU settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, vmID int, diags *diag.Diagnostics) Value {
	cloudinit := Model{}

	devices := config.CustomStorageDevices.Filter(func(device *vms.CustomStorageDevice) bool {
		return device.IsCloudInitDrive(vmID)
	})

	if len(devices) != 1 {
		return types.ObjectNull(attributeTypes())
	}

	for iface, device := range devices {
		cloudinit.Interface = types.StringValue(iface)
		cloudinit.DatastoreId = types.StringValue(device.GetDatastoreID())

		dns := ModelDNS{}
		dns.Domain = types.StringPointerValue(config.CloudInitDNSDomain)

		if config.CloudInitDNSServer != nil {
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

			cloudinit.DNS = dnsObj
		}

		obj, d := types.ObjectValueFrom(ctx, attributeTypes(), cloudinit)
		diags.Append(d...)

		return obj
	}

	return types.ObjectNull(attributeTypes())
}

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
		FileVolume: fmt.Sprintf("%s:cloudinit", plan.DatastoreId.ValueString()),
		Media:      ptr.Ptr("cdrom"),
	}

	body.AddCustomStorageDevice(plan.Interface.ValueString(), device)
}
