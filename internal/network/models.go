/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	pvetypes "github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
)

type linuxBridgeResourceModel struct {
	// Base attributes
	ID        types.String         `tfsdk:"id"`
	NodeName  types.String         `tfsdk:"node_name"`
	Iface     types.String         `tfsdk:"iface"`
	Address   pvetypes.IPCIDRValue `tfsdk:"address"`
	Gateway   pvetypes.IPAddrValue `tfsdk:"gateway"`
	Address6  pvetypes.IPCIDRValue `tfsdk:"address6"`
	Gateway6  pvetypes.IPAddrValue `tfsdk:"gateway6"`
	Autostart types.Bool           `tfsdk:"autostart"`
	MTU       types.Int64          `tfsdk:"mtu"`
	Comment   types.String         `tfsdk:"comment"`
	// Linux bridge attributes
	BridgePorts     []types.String `tfsdk:"bridge_ports"`
	BridgeVLANAware types.Bool     `tfsdk:"bridge_vlan_aware"`
}

//nolint:lll
func (m *linuxBridgeResourceModel) exportToNetworkInterfaceCreateUpdateBody() *nodes.NetworkInterfaceCreateUpdateRequestBody {
	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     m.Iface.ValueString(),
		Type:      "bridge",
		Autostart: pvetypes.CustomBool(m.Autostart.ValueBool()).Pointer(),
	}

	body.CIDR = m.Address.ValueStringPointer()
	body.Gateway = m.Gateway.ValueStringPointer()
	body.CIDR6 = m.Address6.ValueStringPointer()
	body.Gateway6 = m.Gateway6.ValueStringPointer()

	if !m.MTU.IsUnknown() {
		body.MTU = m.MTU.ValueInt64Pointer()
	}

	body.Comments = m.Comment.ValueStringPointer()

	var sanitizedPorts []string

	for i := 0; i < len(m.BridgePorts); i++ {
		port := strings.TrimSpace(m.BridgePorts[i].ValueString())
		if len(port) > 0 {
			sanitizedPorts = append(sanitizedPorts, port)
		}
	}
	sort.Strings(sanitizedPorts)
	bridgePorts := strings.Join(sanitizedPorts, " ")

	if len(bridgePorts) > 0 {
		body.BridgePorts = &bridgePorts
	}

	body.BridgeVLANAware = pvetypes.CustomBool(m.BridgeVLANAware.ValueBool()).Pointer()

	return body
}

func (m *linuxBridgeResourceModel) importFromNetworkInterfaceList(
	ctx context.Context,
	iface *nodes.NetworkInterfaceListResponseData,
) error {
	m.Address = pvetypes.NewIPCIDRPointerValue(iface.CIDR)
	m.Gateway = pvetypes.NewIPAddrPointerValue(iface.Gateway)
	m.Address6 = pvetypes.NewIPCIDRPointerValue(iface.CIDR6)
	m.Gateway6 = pvetypes.NewIPAddrPointerValue(iface.Gateway6)
	m.Autostart = types.BoolPointerValue(iface.Autostart.PointerBool())

	if iface.MTU != nil {
		if v, err := strconv.Atoi(*iface.MTU); err == nil {
			m.MTU = types.Int64Value(int64(v))
		}
	} else {
		m.MTU = types.Int64Null()
	}

	if iface.Comments != nil {
		m.Comment = types.StringValue(strings.TrimSpace(*iface.Comments))
	} else {
		m.Comment = types.StringNull()
	}

	if iface.BridgeVLANAware != nil {
		m.BridgeVLANAware = types.BoolPointerValue(iface.BridgeVLANAware.PointerBool())
	} else {
		m.BridgeVLANAware = types.BoolValue(false)
	}

	if iface.BridgePorts != nil && len(*iface.BridgePorts) > 0 {
		ports, diags := types.ListValueFrom(ctx, types.StringType, strings.Split(*iface.BridgePorts, " "))
		if diags.HasError() {
			return fmt.Errorf("failed to parse bridge ports: %s", *iface.BridgePorts)
		}

		diags = ports.ElementsAs(ctx, &m.BridgePorts, false)
		if diags.HasError() {
			return fmt.Errorf("failed to build bridge ports list: %s", *iface.BridgePorts)
		}
	}

	return nil
}
