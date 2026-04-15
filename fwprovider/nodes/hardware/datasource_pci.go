/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardware

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/hardware"
)

var (
	_ datasource.DataSource              = &pciDataSource{}
	_ datasource.DataSourceWithConfigure = &pciDataSource{}
)

// pciDataSource is the implementation of the proxmox_hardware_pci data source.
type pciDataSource struct {
	client proxmox.Client
}

// NewPCIDataSource creates a new PCI hardware data source.
func NewPCIDataSource() datasource.DataSource {
	return &pciDataSource{}
}

// Metadata defines the data source type name.
func (d *pciDataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_hardware_pci"
}

// Schema defines the schema for the data source.
func (d *pciDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of PCI devices present on a specific Proxmox VE node.",
		MarkdownDescription: "Retrieves the list of PCI devices present on a specific Proxmox VE node. " +
			"This is useful for discovering PCI devices available for passthrough or " +
			"for configuring [hardware mappings](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/hardware_mapping_pci).",
		Attributes: map[string]schema.Attribute{
			"node_name": schema.StringAttribute{
				Description: "The name of the node to list PCI devices from.",
				Required:    true,
			},
			"pci_class_blacklist": schema.ListAttribute{
				Description: "A list of PCI class IDs (hex prefixes) to exclude from the results. " +
					"If not set, the Proxmox default blacklist is used which filters out " +
					"memory controllers (05), bridges (06), and processors (0b). " +
					"Set to an empty list to return all devices.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"filters": schema.SingleNestedAttribute{
				Description: "Client-side filters for narrowing down results. All filters use prefix matching.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Filter by PCI address prefix (e.g. `0000:01` to match all devices on bus 01).",
						Optional:    true,
					},
					"class": schema.StringAttribute{
						Description: "Filter by PCI class code prefix (e.g. `03` to match all display controllers). " +
							"The `0x` prefix in class codes is stripped before matching.",
						Optional: true,
					},
					"vendor_id": schema.StringAttribute{
						Description: "Filter by vendor ID prefix (e.g. `8086` for Intel devices). " +
							"The `0x` prefix in vendor IDs is stripped before matching.",
						Optional: true,
					},
					"device_id": schema.StringAttribute{
						Description: "Filter by device ID prefix. " +
							"The `0x` prefix in device IDs is stripped before matching.",
						Optional: true,
					},
				},
			},
			"devices": schema.ListNestedAttribute{
				Description: "The list of PCI devices.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The PCI address in `domain:bus:device.function` format (e.g. `0000:00:02.0`).",
							Computed:    true,
						},
						"class": schema.StringAttribute{
							Description: "The PCI class code (hex, e.g. `0x030000`).",
							Computed:    true,
						},
						"device": schema.StringAttribute{
							Description: "The PCI device ID (hex, e.g. `0x5916`).",
							Computed:    true,
						},
						"device_name": schema.StringAttribute{
							Description: "The human-readable device name.",
							Computed:    true,
						},
						"iommu_group": schema.Int64Attribute{
							Description: "The IOMMU group number. `-1` indicates that the device is not in an IOMMU group.",
							Computed:    true,
						},
						"mdev": schema.BoolAttribute{
							Description: "Whether the device supports mediated devices (vGPU).",
							Computed:    true,
						},
						"subsystem_device": schema.StringAttribute{
							Description: "The PCI subsystem device ID (hex).",
							Computed:    true,
						},
						"subsystem_device_name": schema.StringAttribute{
							Description: "The human-readable subsystem device name.",
							Computed:    true,
						},
						"subsystem_vendor": schema.StringAttribute{
							Description: "The PCI subsystem vendor ID (hex).",
							Computed:    true,
						},
						"subsystem_vendor_name": schema.StringAttribute{
							Description: "The human-readable subsystem vendor name.",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "The PCI vendor ID (hex, e.g. `0x8086`).",
							Computed:    true,
						},
						"vendor_name": schema.StringAttribute{
							Description: "The human-readable vendor name (e.g. `Intel Corporation`).",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure sets the client for the data source.
func (d *pciDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client
}

// Read fetches PCI device data from the Proxmox API.
func (d *pciDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model pciDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hwClient := d.client.Node(model.NodeName.ValueString()).Hardware()

	reqBody := &hardware.ListPCIDevicesRequestBody{}

	if !model.PCIClassBlacklist.IsNull() {
		var classes []string

		resp.Diagnostics.Append(model.PCIClassBlacklist.ElementsAs(ctx, &classes, false)...)

		if resp.Diagnostics.HasError() {
			return
		}

		blacklist := strings.Join(classes, ";")
		reqBody.ClassBlacklist = &blacklist
	}

	devices, err := hwClient.ListPCIDevices(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read PCI Devices", err.Error())
		return
	}

	model.Devices = make([]pciDevice, 0, len(devices))

	for _, dev := range devices {
		if model.Filters != nil && !matchesPCIFilters(dev, model.Filters) {
			continue
		}

		model.Devices = append(model.Devices, pciDeviceFromAPI(dev))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// matchesPCIFilters checks whether a PCI device matches all configured filters.
// All filters use prefix matching. The "0x" prefix is stripped from hex values before comparison.
func matchesPCIFilters(dev *hardware.PCIDeviceData, f *pciFilters) bool {
	if !f.ID.IsNull() && !strings.HasPrefix(dev.ID, f.ID.ValueString()) {
		return false
	}

	if !f.Class.IsNull() && !hexPrefixMatch(dev.Class, f.Class.ValueString()) {
		return false
	}

	if !f.VendorID.IsNull() && !hexPrefixMatch(dev.Vendor, f.VendorID.ValueString()) {
		return false
	}

	if !f.DeviceID.IsNull() && !hexPrefixMatch(dev.Device, f.DeviceID.ValueString()) {
		return false
	}

	return true
}

// hexPrefixMatch compares two hex values by prefix after stripping "0x" and lowercasing.
func hexPrefixMatch(value, prefix string) bool {
	v := strings.TrimPrefix(strings.ToLower(value), "0x")
	p := strings.TrimPrefix(strings.ToLower(prefix), "0x")

	return strings.HasPrefix(v, p)
}
