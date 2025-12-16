/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// Model represents the cloned VM resource.
type Model struct {
	ID                               types.Int64             `tfsdk:"id"`
	NodeName                         types.String            `tfsdk:"node_name"`
	Name                             types.String            `tfsdk:"name"`
	Description                      types.String            `tfsdk:"description"`
	Tags                             stringset.Value         `tfsdk:"tags"`
	Clone                            CloneModel              `tfsdk:"clone"`
	Network                          map[string]NetworkModel `tfsdk:"network"`
	Disk                             map[string]DiskModel    `tfsdk:"disk"`
	Delete                           *DeleteModel            `tfsdk:"delete"`
	CPU                              cpu.Value               `tfsdk:"cpu"`
	RNG                              rng.Value               `tfsdk:"rng"`
	VGA                              vga.Value               `tfsdk:"vga"`
	CDROM                            cdrom.Value             `tfsdk:"cdrom"`
	StopOnDestroy                    types.Bool              `tfsdk:"stop_on_destroy"`
	PurgeOnDestroy                   types.Bool              `tfsdk:"purge_on_destroy"`
	DeleteUnreferencedDisksOnDestroy types.Bool              `tfsdk:"delete_unreferenced_disks_on_destroy"`
	Timeouts                         timeouts.Value          `tfsdk:"timeouts"`
}

// CloneModel captures clone parameters.
type CloneModel struct {
	SourceVMID      types.Int64  `tfsdk:"source_vm_id"`
	SourceNodeName  types.String `tfsdk:"source_node_name"`
	Full            types.Bool   `tfsdk:"full"`
	TargetDatastore types.String `tfsdk:"target_datastore"`
	TargetFormat    types.String `tfsdk:"target_format"`
	SnapshotName    types.String `tfsdk:"snapshot_name"`
	PoolID          types.String `tfsdk:"pool_id"`
	Retries         types.Int64  `tfsdk:"retries"`
	BandwidthLimit  types.Int64  `tfsdk:"bandwidth_limit"`
}

// DeleteModel holds explicit delete lists.
type DeleteModel struct {
	Network []types.String `tfsdk:"network"`
	Disk    []types.String `tfsdk:"disk"`
}

// NetworkModel represents a managed network device slot.
type NetworkModel struct {
	Bridge     types.String  `tfsdk:"bridge"`
	Model      types.String  `tfsdk:"model"`
	Firewall   types.Bool    `tfsdk:"firewall"`
	LinkDown   types.Bool    `tfsdk:"link_down"`
	MACAddress types.String  `tfsdk:"mac_address"`
	MTU        types.Int64   `tfsdk:"mtu"`
	Queues     types.Int64   `tfsdk:"queues"`
	RateLimit  types.Float64 `tfsdk:"rate_limit"`
	Tag        types.Int64   `tfsdk:"tag"`
	Trunks     types.Set     `tfsdk:"trunks"`
}

// DiskModel represents a managed disk slot.
type DiskModel struct {
	File        types.String `tfsdk:"file"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	SizeGB      types.Int64  `tfsdk:"size_gb"`
	Format      types.String `tfsdk:"format"`
	AIO         types.String `tfsdk:"aio"`
	Backup      types.Bool   `tfsdk:"backup"`
	Discard     types.String `tfsdk:"discard"`
	Cache       types.String `tfsdk:"cache"`
	IOThread    types.Bool   `tfsdk:"iothread"`
	Replicate   types.Bool   `tfsdk:"replicate"`
	Serial      types.String `tfsdk:"serial"`
	SSD         types.Bool   `tfsdk:"ssd"`
	ImportFrom  types.String `tfsdk:"import_from"`
	Media       types.String `tfsdk:"media"`
}
