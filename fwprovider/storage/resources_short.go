/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

// Short-name alias wrappers for storage resources (ADR-007 Phase 2).
// Each wrapper embeds the concrete resource, overrides Metadata/Schema/MoveState
// to register under the short proxmox_storage_* name without deprecation and
// with state migration support from the old proxmox_virtual_environment_storage_* name.

// --- CIFS ---

var (
	_ resource.Resource              = &cifsStorageShort{}
	_ resource.ResourceWithMoveState = &cifsStorageShort{}
)

type cifsStorageShort struct{ cifsStorageResource }

// NewCIFSStorageShortResource creates the short-named proxmox_storage_cifs resource.
func NewCIFSStorageShortResource() resource.Resource {
	return &cifsStorageShort{
		cifsStorageResource: cifsStorageResource{
			storageResource: &storageResource[*CIFSStorageModel, CIFSStorageModel]{
				storageType:  "cifs",
				resourceName: "proxmox_storage_cifs",
			},
		},
	}
}

func (r *cifsStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_cifs"
}

func (r *cifsStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.cifsStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *cifsStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.cifsStorageResource.Schema, "proxmox_virtual_environment_storage_cifs")
}

// --- Directory ---

var (
	_ resource.Resource              = &directoryStorageShort{}
	_ resource.ResourceWithMoveState = &directoryStorageShort{}
)

type directoryStorageShort struct{ directoryStorageResource }

// NewDirectoryStorageShortResource creates the short-named proxmox_storage_directory resource.
func NewDirectoryStorageShortResource() resource.Resource {
	return &directoryStorageShort{
		directoryStorageResource: directoryStorageResource{
			storageResource: &storageResource[*DirectoryStorageModel, DirectoryStorageModel]{
				storageType:  "dir",
				resourceName: "proxmox_storage_directory",
			},
		},
	}
}

func (r *directoryStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_directory"
}

func (r *directoryStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.directoryStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *directoryStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.directoryStorageResource.Schema, "proxmox_virtual_environment_storage_directory")
}

// --- LVM ---

var (
	_ resource.Resource              = &lvmPoolStorageShort{}
	_ resource.ResourceWithMoveState = &lvmPoolStorageShort{}
)

type lvmPoolStorageShort struct{ lvmPoolStorageResource }

// NewLVMPoolStorageShortResource creates the short-named proxmox_storage_lvm resource.
func NewLVMPoolStorageShortResource() resource.Resource {
	return &lvmPoolStorageShort{
		lvmPoolStorageResource: lvmPoolStorageResource{
			storageResource: &storageResource[*LVMStorageModel, LVMStorageModel]{
				storageType:  "lvm",
				resourceName: "proxmox_storage_lvm",
			},
		},
	}
}

func (r *lvmPoolStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_lvm"
}

func (r *lvmPoolStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.lvmPoolStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *lvmPoolStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.lvmPoolStorageResource.Schema, "proxmox_virtual_environment_storage_lvm")
}

// --- LVM Thin ---

var (
	_ resource.Resource              = &lvmThinPoolStorageShort{}
	_ resource.ResourceWithMoveState = &lvmThinPoolStorageShort{}
)

type lvmThinPoolStorageShort struct{ lvmThinPoolStorageResource }

// NewLVMThinPoolStorageShortResource creates the short-named proxmox_storage_lvmthin resource.
func NewLVMThinPoolStorageShortResource() resource.Resource {
	return &lvmThinPoolStorageShort{
		lvmThinPoolStorageResource: lvmThinPoolStorageResource{
			storageResource: &storageResource[*LVMThinStorageModel, LVMThinStorageModel]{
				storageType:  "lvmthin",
				resourceName: "proxmox_storage_lvmthin",
			},
		},
	}
}

func (r *lvmThinPoolStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_lvmthin"
}

func (r *lvmThinPoolStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.lvmThinPoolStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *lvmThinPoolStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.lvmThinPoolStorageResource.Schema, "proxmox_virtual_environment_storage_lvmthin")
}

// --- NFS ---

var (
	_ resource.Resource              = &nfsStorageShort{}
	_ resource.ResourceWithMoveState = &nfsStorageShort{}
)

type nfsStorageShort struct{ nfsStorageResource }

// NewNFSStorageShortResource creates the short-named proxmox_storage_nfs resource.
func NewNFSStorageShortResource() resource.Resource {
	return &nfsStorageShort{
		nfsStorageResource: nfsStorageResource{
			storageResource: &storageResource[*NFSStorageModel, NFSStorageModel]{
				storageType:  "nfs",
				resourceName: "proxmox_storage_nfs",
			},
		},
	}
}

func (r *nfsStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_nfs"
}

func (r *nfsStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.nfsStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *nfsStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.nfsStorageResource.Schema, "proxmox_virtual_environment_storage_nfs")
}

// --- PBS ---

var (
	_ resource.Resource              = &pbsStorageShort{}
	_ resource.ResourceWithMoveState = &pbsStorageShort{}
)

type pbsStorageShort struct{ pbsStorageResource }

// NewProxmoxBackupServerStorageShortResource creates the short-named proxmox_storage_pbs resource.
func NewProxmoxBackupServerStorageShortResource() resource.Resource {
	return &pbsStorageShort{
		pbsStorageResource: pbsStorageResource{
			storageResource: &storageResource[*PBSStorageModel, PBSStorageModel]{
				storageType:  "pbs",
				resourceName: "proxmox_storage_pbs",
			},
		},
	}
}

func (r *pbsStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_pbs"
}

func (r *pbsStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.pbsStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *pbsStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.pbsStorageResource.Schema, "proxmox_virtual_environment_storage_pbs")
}

// --- ZFS Pool ---

var (
	_ resource.Resource              = &zfsPoolStorageShort{}
	_ resource.ResourceWithMoveState = &zfsPoolStorageShort{}
)

type zfsPoolStorageShort struct{ zfsPoolStorageResource }

// NewZFSPoolStorageShortResource creates the short-named proxmox_storage_zfspool resource.
func NewZFSPoolStorageShortResource() resource.Resource {
	return &zfsPoolStorageShort{
		zfsPoolStorageResource: zfsPoolStorageResource{
			storageResource: &storageResource[*ZFSStorageModel, ZFSStorageModel]{
				storageType:  "zfspool",
				resourceName: "proxmox_storage_zfspool",
			},
		},
	}
}

func (r *zfsPoolStorageShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_storage_zfspool"
}

func (r *zfsPoolStorageShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.zfsPoolStorageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *zfsPoolStorageShort) MoveState(ctx context.Context) []resource.StateMover {
	return storageShortMoveState(ctx, r.zfsPoolStorageResource.Schema, "proxmox_virtual_environment_storage_zfspool")
}

// --- Shared helper ---

type schemaFunc func(context.Context, resource.SchemaRequest, *resource.SchemaResponse)

func storageShortMoveState(ctx context.Context, schemaFn schemaFunc, oldTypeName string) []resource.StateMover {
	var schemaResp resource.SchemaResponse

	schemaFn(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState(oldTypeName, &schemaResp.Schema),
	}
}
