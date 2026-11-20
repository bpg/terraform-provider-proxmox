/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

var (
	_ resource.Resource              = &bootstrapResource{}
	_ resource.ResourceWithConfigure = &bootstrapResource{}
)

type bootstrapResource struct {
	client proxmox.Client
}

// NewResource creates the proxmox_cloudinit_bootstrap resource.
func NewResource() resource.Resource {
	return &bootstrapResource{}
}

func (r *bootstrapResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_cloudinit_bootstrap"
}

// contentString is a seed file attribute. cloud-init applies the seed only on first boot, so Update rejects changes.
func contentString(desc, mdDesc string, required bool) schema.StringAttribute {
	return schema.StringAttribute{
		Description:         desc,
		MarkdownDescription: mdDesc,
		Required:            required,
		Optional:            !required,
		Sensitive:           true,
	}
}

func (r *bootstrapResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Bootstraps a stopped VM with cloud-init through a temporary NoCloud seed ISO, without SSH.",
		MarkdownDescription: "Bootstraps a stopped VM with cloud-init through a temporary " +
			"[NoCloud](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) seed ISO, without SSH. " +
			"The ISO is removed once cloud-init finishes, so the VM can be migrated freely.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("The VM ID."),
			"vm_id": schema.Int64Attribute{
				Description:   "The VM to bootstrap.",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"datastore_id": schema.StringAttribute{
				Description:         "The datastore to upload the temporary seed ISO to.",
				MarkdownDescription: "The datastore on the VM's node to upload the temporary seed ISO to. It must allow `iso` content.",
				Required:            true,
			},
			"file_name": schema.StringAttribute{
				Description: "The file name of the temporary seed ISO.",
				MarkdownDescription: "The file name of the temporary seed ISO on `datastore_id`. Defaults to " +
					"`cloudinit-bootstrap-<vm_id>.iso`. It must end in `.iso` and use only letters, digits, `.`, `_` and `-` " +
					"(Proxmox would rename other characters). Creation fails if a file with this name already exists on the " +
					"datastore. A custom name must be unique on the datastore across all configurations that bootstrap VMs " +
					"there.",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
					stringvalidator.RegexMatches(isoFileNameRE, isoFileNameMsg),
				},
			},
			"interface": schema.StringAttribute{
				Description: "The empty CD-ROM drive to attach the seed ISO to.",
				MarkdownDescription: "The empty CD-ROM drive to attach the seed ISO to, one of `ide[0-3]`, `sata[0-5]` or " +
					"`scsi[0-30]`. Defaults to `ide0`. Creation fails unless the VM has an empty CD-ROM (`file_id = \"none\"`) " +
					"there.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ide0"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(cdromInterfaceRE, "one of `ide[0-3]`, `sata[0-5]`, `scsi[0-30]`"),
				},
			},
			"user_data": contentString("The user-data file contents.",
				"The `user-data` file contents, e.g. a `#cloud-config` document.", true),
			"meta_data": contentString("The meta-data file contents.",
				"The `meta-data` file contents. Defaults to `instance-id` and `local-hostname` set to the VM name.", false),
			"network_config": contentString("The network-config file contents.",
				"The `network-config` file contents (network config version 1 or 2).", false),
			"vendor_data": contentString("The vendor-data file contents.", "The `vendor-data` file contents.", false),
			"timeout": schema.Int64Attribute{
				Description:         "How long to wait for cloud-init to finish, in seconds.",
				MarkdownDescription: "How long to wait for cloud-init to finish, in seconds. Defaults to `1800`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1800),
				Validators:          []validator.Int64{int64validator.AtLeast(1)},
			},
			"node_name": schema.StringAttribute{
				Description: "The node the VM was last seen on.",
				MarkdownDescription: "The node the VM was last seen on. It follows migrations, and is used to find the VM " +
					"when the cluster's resource list misses it.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"seed_file_id": schema.StringAttribute{
				Description: "The volume ID of the seed ISO this resource uploaded, while it is still on the datastore.",
				MarkdownDescription: "The volume ID of the seed ISO this resource uploaded, while it is still on the " +
					"datastore (after a failed or interrupted bootstrap). Empty otherwise; refresh clears it once the file is " +
					"gone. Destroying the resource deletes only this file.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"ipv4_addresses": schema.MapAttribute{
				Description: "IPv4 addresses of the VM's network devices, keyed by guest interface name.",
				MarkdownDescription: "IPv4 addresses of the VM's network devices reported by the guest agent, keyed by " +
					"guest interface name. Refreshed while the VM is running and every device has an address; the " +
					"last known addresses are kept otherwise.",
				Computed:      true,
				ElementType:   ipListType(),
				PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *bootstrapResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

// isoFileNameRE matches the names the PVE upload API keeps as they are; it replaces other characters with "_".
var isoFileNameRE = regexp.MustCompile(`^[A-Za-z0-9_.-]+\.iso$`)

const isoFileNameMsg = "must end in .iso and use only letters, digits, '.', '_' and '-'"

// cdromInterfaceRE is the slot range the cdrom block accepts (qemu-server: MAX_IDE_DISKS=4, MAX_SATA_DISKS=6,
// MAX_SCSI_DISKS=31).
var cdromInterfaceRE = regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi([0-9]|[12][0-9]|30))$`)

func seedVolumeID(datastore, fileName string) string {
	return datastore + ":" + storage.ContentTypeISO + "/" + fileName
}

func defaultFileName(vmID int64) string {
	return fmt.Sprintf("cloudinit-bootstrap-%d.iso", vmID)
}

// seedCleanupError means cloud-init finished, so the VM is bootstrapped, but removing the seed ISO failed.
type seedCleanupError struct{ error }

// manualCleanup says how to remove a seed ISO by hand.
func manualCleanup(node string, vmID int64, iface, volID string) string {
	if node == "" {
		node = "<the VM's node>"
	}

	return fmt.Sprintf("on node %s eject it (qm set %d --%s none,media=cdrom) and delete it (pvesm free %s)",
		node, vmID, iface, volID)
}

func (r *bootstrapResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan bootstrapModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vmID := plan.VMID.ValueInt64()
	plan.ID = types.StringValue(fmt.Sprintf("%d", vmID))

	if plan.FileName.IsUnknown() {
		plan.FileName = types.StringValue(defaultFileName(vmID))
	}

	ips, err := r.bootstrap(ctx, &plan, &resp.Diagnostics)
	if err != nil {
		summary := fmt.Sprintf("Unable to Bootstrap VM %d", vmID)

		if _, ok := errors.AsType[*seedCleanupError](err); ok {
			summary = fmt.Sprintf("Unable to Remove Seed ISO of VM %d", vmID)
		}

		resp.Diagnostics.AddError(summary, err.Error())

		// seed_file_id is known once an ISO may have been uploaded. Only then is a (tainted) state stored, so destroy or
		// the next replacement removes exactly that file; earlier failures leave nothing behind and store nothing.
		if !plan.SeedFileID.IsUnknown() {
			plan.IPv4Addresses = emptyIPv4Map()
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		}

		return
	}

	plan.IPv4Addresses = ips

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// findVM returns the VM's node and client. The cluster's resource list follows migrations but can lag behind a VM
// created moments ago on another node, so a miss there is checked by reading the VM's config on the node it was last
// seen on, or on every online node when none is known yet. Only when that read also misses is the VM gone
// (cluster.ErrVMDoesNotExist); any other error leaves it undecided.
func (r *bootstrapResource) findVM(ctx context.Context, vmID int64, lastNode string) (string, *vms.Client, error) {
	nodeName, err := r.client.Cluster().GetVMNodeName(ctx, int(vmID))
	if err == nil {
		return *nodeName, r.client.Node(*nodeName).VM(int(vmID)), nil
	}

	if !errors.Is(err, cluster.ErrVMDoesNotExist) {
		return "", nil, fmt.Errorf("unable to find the node of VM %d: %w", vmID, err)
	}

	candidates := []string{lastNode}

	if lastNode == "" {
		list, err := r.client.Node("").ListNodes(ctx)
		if err != nil {
			return "", nil, fmt.Errorf("unable to find the node of VM %d: %w", vmID, err)
		}

		candidates = nil

		for _, n := range list {
			if n.Status != nil && *n.Status == "online" {
				candidates = append(candidates, n.Name)
			}
		}
	}

	var undecided error

	for _, n := range candidates {
		vm := r.client.Node(n).VM(int(vmID))

		_, err := vm.GetVM(ctx)
		if err == nil {
			return n, vm, nil
		}

		// A node that doesn't answer may still hold the VM, so keep looking but don't call it gone.
		if !errors.Is(err, api.ErrResourceDoesNotExist) && undecided == nil {
			undecided = fmt.Errorf("unable to read VM %d on node %s: %w", vmID, n, err)
		}
	}

	if undecided != nil {
		return "", nil, undecided
	}

	return "", nil, fmt.Errorf("VM %d: %w", vmID, cluster.ErrVMDoesNotExist)
}

// seedUploadMu serializes the exists check and the upload, so two bootstraps with the same file_name in one apply
// can't both pass the check and boot each other's seed. It doesn't cover other provider processes.
//
//nolint:gochecknoglobals
var seedUploadMu sync.Mutex

// bootstrap returns the guest's IPv4 addresses. Warnings that don't fail the bootstrap go to diags.
func (r *bootstrapResource) bootstrap(ctx context.Context, m *bootstrapModel, diags *diag.Diagnostics) (types.Map, error) {
	none := types.MapNull(ipListType())
	vmID := m.VMID.ValueInt64()
	iface := m.Interface.ValueString()
	datastore := m.DatastoreID.ValueString()
	fileName := m.FileName.ValueString()

	nodeName, vm, err := r.findVM(ctx, vmID, "")
	if err != nil {
		return none, err
	}

	m.NodeName = types.StringValue(nodeName)
	datastoreClient := r.client.Node(nodeName).Storage(datastore)

	status, err := vm.GetVMStatus(ctx)
	if err != nil {
		return none, err
	}

	if status.Status == "running" {
		return none, fmt.Errorf("VM %d is already running, and cloud-init applies the seed only on a VM's first boot. "+
			"For a new VM, set started = false and lifecycle { ignore_changes = [started] } on the VM resource, then "+
			"replace the VM. If an earlier apply reported \"Unable to Remove Seed ISO\", the VM is already bootstrapped: "+
			"remove this resource from the configuration instead", vmID)
	}

	cfg, err := vm.GetVM(ctx)
	if err != nil {
		return none, err
	}

	if err := checkVMConfig(cfg, vmID, iface); err != nil {
		return none, fmt.Errorf("%w. Nothing was changed", err)
	}

	vmName := fmt.Sprintf("vm-%d", vmID)
	if cfg.Name != nil && *cfg.Name != "" {
		vmName = *cfg.Name
	}

	files := seedFiles(vmName, m.UserData.ValueString(), m.MetaData.ValueString(),
		m.NetworkConfig.ValueString(), m.VendorData.ValueString())

	volID := seedVolumeID(datastore, fileName)

	if err := r.uploadSeed(ctx, m, datastoreClient, nodeName, files); err != nil {
		return none, err
	}

	if err := setCDROM(ctx, vm, iface, volID); err != nil {
		return none, err
	}

	res := vm.StartVM(ctx, int(m.Timeout.ValueInt64()))
	if err := res.Err(); err != nil {
		return none, fmt.Errorf("unable to start VM %d: %w", vmID, err)
	}

	res.AddDiagsAsWarnings(diags, fmt.Sprintf("Warnings While Starting VM %d", vmID))

	if err := waitForCloudInit(ctx, vm, time.Duration(m.Timeout.ValueInt64())*time.Second); err != nil {
		return none, fmt.Errorf("%w; the seed ISO is left attached for debugging. cloud-init applies the seed "+
			"only on a VM's first boot, so replace the VM to bootstrap it again", err)
	}

	if err := r.removeSeedISO(ctx, vm, datastoreClient, iface, volID, diags); err != nil {
		return none, &seedCleanupError{fmt.Errorf("cloud-init finished and VM %d is bootstrapped; only removing "+
			"%s failed: %w. Don't replace the VM and don't apply yet: %s, then run untaint on this resource",
			vmID, volID, err, manualCleanup(nodeName, vmID, iface, volID))}
	}

	m.SeedFileID = types.StringNull()

	// The VM is bootstrapped by now, so missing addresses only warn: failing would taint it and replace a good VM.
	ips, complete, err := waitForIPv4Addresses(ctx, vm)

	switch {
	case err != nil:
		diags.AddWarning(fmt.Sprintf("IPv4 Addresses Unavailable for VM %d", vmID),
			fmt.Sprintf("cloud-init finished and VM %d is bootstrapped, but the guest agent didn't report its network "+
				"interfaces: %s. ipv4_addresses is empty; later refreshes fill it in once the agent answers.", vmID, err))

		return emptyIPv4Map(), nil
	case !complete:
		diags.AddWarning(fmt.Sprintf("Incomplete IPv4 Addresses for VM %d", vmID),
			"cloud-init finished, but not every network device of the VM had an IPv4 address after 2 minutes "+
				"(e.g. DHCP still pending). ipv4_addresses is incomplete; later refreshes fill it in once every "+
				"device has one.")
	}

	return ips, nil
}

// uploadSeed checks that the seed's file name is free, records it in seed_file_id and uploads the ISO.
func (r *bootstrapResource) uploadSeed(
	ctx context.Context,
	m *bootstrapModel,
	datastoreClient *storage.Client,
	nodeName string,
	files map[string]string,
) error {
	datastore := m.DatastoreID.ValueString()
	volID := seedVolumeID(datastore, m.FileName.ValueString())

	seedUploadMu.Lock()
	defer seedUploadMu.Unlock()

	// Checked once, outside the upload's retries: PVE overwrites on upload, so this is the only guard against
	// replacing a file this resource didn't upload.
	exists, err := isoListed(ctx, datastoreClient, volID)
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		return fmt.Errorf("datastore %q not found on node %s", datastore, nodeName)
	}

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("%s already exists in the datastore and this resource didn't upload it, so it "+
			"was left alone. If it's a leftover seed ISO (it may hold secrets), eject it from any VM using it and delete "+
			"it on node %s (pvesm free %s), then apply again; otherwise set file_name to an unused name",
			volID, nodeName, volID)
	}

	// Recorded before the upload: a failed or cancelled upload can still leave the file, and destroy or the next
	// replacement must remove it.
	m.SeedFileID = types.StringValue(volID)

	return uploadISO(ctx, r.client, datastoreClient, files, m.FileName.ValueString())
}

// checkVMConfig fails before anything is uploaded when the bootstrap can't work: the wait needs the guest agent, and
// attaching the seed replaces whatever is on iface.
func checkVMConfig(cfg *vms.GetResponseData, vmID int64, iface string) error {
	if cfg.Agent == nil || cfg.Agent.Enabled == nil || !bool(*cfg.Agent.Enabled) {
		return fmt.Errorf("the QEMU guest agent is disabled on VM %d, and this resource reads cloud-init's result "+
			"through it. Add agent { enabled = true } to the VM and install qemu-guest-agent in the guest (user_data "+
			"can do it), then apply again", vmID)
	}

	for slot, dev := range cfg.StorageDevices {
		if dev != nil && dev.IsCloudInitDrive(int(vmID)) {
			return fmt.Errorf("VM %d has a Proxmox cloud-init drive on %s, so the guest would see two NoCloud seeds and "+
				"could boot either. Remove the drive: create the VM without an initialization block, or, for a clone of a "+
				"template that has the drive, replace it with an empty CD-ROM (cdrom { interface = \"%s\" file_id = "+
				"\"none\" }) and set interface = \"%s\" here", vmID, slot, slot, slot)
		}
	}

	dev := cfg.StorageDevices[iface]

	var found string

	switch {
	case dev == nil:
		found = "not configured"
	case dev.Media == nil || *dev.Media != "cdrom":
		found = fmt.Sprintf("a disk (%s)", dev.FileVolume)
	case dev.FileVolume == "cdrom":
		found = "the host's CD drive"
	case dev.FileVolume != "none":
		found = fmt.Sprintf("a CD-ROM holding %s", dev.FileVolume)
	default:
		return nil
	}

	return fmt.Errorf("%s on VM %d must be an empty CD-ROM, but it is %s, and attaching the seed ISO would replace "+
		"it. Add an empty CD-ROM there (cdrom { interface = \"%s\" file_id = \"none\" }), or set interface to a slot "+
		"that has one", iface, vmID, found, iface)
}

// isoListed reports whether volID is in the datastore's ISO listing. A missing datastore is returned as
// api.ErrResourceDoesNotExist.
func isoListed(ctx context.Context, datastoreClient *storage.Client, volID string) (bool, error) {
	contentType := storage.ContentTypeISO

	files, err := datastoreClient.ListDatastoreFiles(ctx, &contentType)
	if err != nil {
		return false, fmt.Errorf("unable to list ISOs for %s: %w", volID, err)
	}

	for _, f := range files {
		if f.VolumeID == volID {
			return true, nil
		}
	}

	return false, nil
}

// removeSeedISO ejects volID if the drive holds it, then deletes it if it's listed as an ISO, so it never touches
// anything but the seed ISO this resource recorded. vm is nil when the VM is gone, so there is nothing to eject.
func (r *bootstrapResource) removeSeedISO(
	ctx context.Context,
	vm *vms.Client,
	datastoreClient *storage.Client,
	iface, volID string,
	diags *diag.Diagnostics,
) error {
	if vm != nil {
		cfg, err := vm.GetVM(ctx)
		if err != nil {
			return err
		}

		// Eject first: PVE deletes an ISO even while a VM has it attached, and that VM then can't start.
		if dev, ok := cfg.StorageDevices[iface]; ok && dev != nil && dev.FileVolume == volID {
			if err := setCDROM(ctx, vm, iface, "none"); err != nil {
				return err
			}
		}
	}

	// Listing first, because the file may already be gone (a failed upload, or a manual cleanup), and a delete of a
	// missing file fails as a task rather than with a not-found error.
	listed, err := isoListed(ctx, datastoreClient, volID)
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		// The datastore is gone, so there is nothing left to delete and destroy must not block.
		return nil
	}

	if err != nil || !listed {
		return err
	}

	res := datastoreClient.DeleteDatastoreFile(ctx, volID)
	if err := res.Err(); err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		return fmt.Errorf("unable to delete %s: %w", volID, err)
	}

	res.AddDiagsAsWarnings(diags, "Warnings While Deleting Seed ISO "+volID)

	return nil
}

func setCDROM(ctx context.Context, vm *vms.Client, iface, volID string) error {
	media := "cdrom"
	body := &vms.UpdateRequestBody{}
	body.AddCustomStorageDevice(iface, vms.CustomStorageDevice{FileVolume: volID, Media: &media})

	if err := vm.UpdateVM(ctx, body); err != nil {
		return fmt.Errorf("unable to set %s to %s: %w", iface, volID, err)
	}

	return nil
}

func ipListType() types.ListType {
	return types.ListType{ElemType: types.StringType}
}

func emptyIPv4Map() types.Map {
	return types.MapValueMust(ipListType(), map[string]attr.Value{})
}

// ipWaitTimeout bounds how long Create waits for DHCP on every NIC after cloud-init finished; the docs and the
// incomplete-addresses warning say "2 minutes".
const ipWaitTimeout = 2 * time.Minute

var errIncompleteIPv4 = errors.New("not every network device has an IPv4 address yet")

// waitForIPv4Addresses returns the guest's addresses once every NIC has one, or what it has after ipWaitTimeout,
// and whether they are complete. A missing privilege or agent ends the wait at once.
func waitForIPv4Addresses(ctx context.Context, vm *vms.Client) (types.Map, bool, error) {
	waitCtx, cancel := context.WithTimeout(ctx, ipWaitTimeout)
	defer cancel()

	var (
		ips      types.Map
		complete bool
	)

	op := retry.NewPollOperation("guest IPv4 addresses",
		retry.WithBaseDelay(5*time.Second),
		retry.WithRetryIf(func(err error) bool { return !isPermanentAgentError(err) }),
	)

	err := op.DoPoll(waitCtx, func() error {
		m, c, err := guestIPv4Addresses(waitCtx, vm)
		if err != nil {
			return err
		}

		ips, complete = m, c

		if !complete {
			return errIncompleteIPv4
		}

		return nil
	})
	if err != nil && ips.IsNull() {
		return types.MapNull(ipListType()), false, fmt.Errorf("unable to read the guest's IPv4 addresses: %w", err)
	}

	return ips, complete, nil
}

// guestIPv4Addresses returns the IPv4 addresses of the VM's network devices, keyed by guest interface name, and whether
// every device has one. Only interfaces whose MAC belongs to a Proxmox network device are included, so container
// networking inside the guest (CNI bridges, kube-ipvs0, ...) doesn't churn the map.
func guestIPv4Addresses(ctx context.Context, vm *vms.Client) (types.Map, bool, error) {
	none := types.MapNull(ipListType())

	cfg, err := vm.GetVM(ctx)
	if err != nil {
		return none, false, err
	}

	macs := map[string]bool{}

	for _, dev := range cfg.NetworkDevices {
		if dev != nil && dev.MACAddress != nil {
			macs[strings.ToLower(*dev.MACAddress)] = true
		}
	}

	data, err := vm.GetVMNetworkInterfacesFromAgent(ctx)
	if err != nil {
		return none, false, err
	}

	var nics []vms.GetQEMUNetworkInterfacesResponseResult
	if data.Result != nil {
		nics = *data.Result
	}

	out, complete := ipv4ByInterface(macs, nics)

	m, diags := types.MapValue(ipListType(), out)
	if diags.HasError() {
		tflog.Error(ctx, "unable to build ipv4_addresses", map[string]any{"diags": diags})

		return none, false, errors.New("unable to build ipv4_addresses")
	}

	return m, complete, nil
}

// ipv4ByInterface keeps the guest interfaces whose MAC is in macs and reports whether each of macs has an IPv4 address.
// VLAN subinterfaces and bridges in the guest reuse a NIC's MAC, so a MAC counts as covered when any interface
// with it has an address, and an interface without one is only listed if nothing else covers its MAC.
func ipv4ByInterface(macs map[string]bool, nics []vms.GetQEMUNetworkInterfacesResponseResult) (map[string]attr.Value, bool) {
	out := map[string]attr.Value{}
	covered := map[string]bool{}
	empty := map[string][]string{}

	for _, nic := range nics {
		mac := strings.ToLower(nic.MACAddress)
		if !macs[mac] {
			continue
		}

		list := []attr.Value{}

		if nic.IPAddresses != nil {
			for _, a := range *nic.IPAddresses {
				if a.Type == "ipv4" {
					list = append(list, types.StringValue(a.Address))
				}
			}
		}

		out[nic.Name] = types.ListValueMust(types.StringType, list)

		if len(list) > 0 {
			covered[mac] = true
		} else {
			empty[mac] = append(empty[mac], nic.Name)
		}
	}

	for mac, names := range empty {
		if covered[mac] {
			for _, name := range names {
				delete(out, name)
			}
		}
	}

	return out, len(covered) == len(macs)
}

// Read refreshes ipv4_addresses from the guest agent on the VM's current node, so DHCP changes show up, and node_name
// after a migration. The bootstrap itself is one-time, so everything else is kept. So are the last known addresses
// while the VM is stopped, the agent can't answer, or a NIC has no IPv4 address yet (e.g. DHCP during a reboot).
func (r *bootstrapResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state bootstrapModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodeName, vm, err := r.findVM(ctx, state.VMID.ValueInt64(), state.NodeName.ValueString())

	switch {
	case errors.Is(err, cluster.ErrVMDoesNotExist):
		// The VM was deleted outside Terraform. It is planned as a create, which replace_triggered_by doesn't catch,
		// so the bootstrap has to be created again too.
		warnSeedLeft(&resp.Diagnostics, state)
		resp.State.RemoveResource(ctx)

		return
	case err != nil:
		tflog.Debug(ctx, "keeping the last known state", map[string]any{"error": err.Error()})
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

		return
	}

	state.NodeName = types.StringValue(nodeName)

	if ips, err := currentIPv4Addresses(ctx, vm); err != nil {
		tflog.Debug(ctx, "keeping the last known ipv4_addresses", map[string]any{"error": err.Error()})
	} else {
		state.IPv4Addresses = ips
	}

	if volID := state.SeedFileID.ValueString(); volID != "" {
		// Once the file is gone (e.g. deleted by hand after a failed cleanup), there is nothing left to delete, and a
		// file that later appears under the same name isn't this resource's.
		datastore, _, _ := strings.Cut(volID, ":")

		listed, err := isoListed(ctx, r.client.Node(nodeName).Storage(datastore), volID)

		switch {
		case errors.Is(err, api.ErrResourceDoesNotExist), err == nil && !listed:
			state.SeedFileID = types.StringNull()
		case err != nil:
			tflog.Debug(ctx, "keeping seed_file_id", map[string]any{"error": err.Error()})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// warnSeedLeft warns that a recorded seed ISO may outlive its deleted VM: PVE keeps ISOs when it destroys a VM, and
// nothing deletes it later.
func warnSeedLeft(diags *diag.Diagnostics, state bootstrapModel) {
	volID := state.SeedFileID.ValueString()
	if volID == "" {
		return
	}

	node := state.NodeName.ValueString()
	if node == "" {
		node = "a node that has it"
	}

	diags.AddWarning("Seed ISO Left on Datastore",
		fmt.Sprintf("VM %d no longer exists, but the seed ISO %s that a failed bootstrap left may still be on the "+
			"datastore (it may hold secrets). Delete it with pvesm free %s on %s.",
			state.VMID.ValueInt64(), volID, volID, node))
}

func currentIPv4Addresses(ctx context.Context, vm *vms.Client) (types.Map, error) {
	none := types.MapNull(ipListType())

	status, err := vm.GetVMStatus(ctx)
	if err != nil {
		return none, err
	}

	if status.Status != "running" {
		return none, fmt.Errorf("VM %d is %s", vm.VMID, status.Status)
	}

	ips, complete, err := guestIPv4Addresses(ctx, vm)
	if err != nil {
		return none, err
	}

	if !complete {
		return none, errIncompleteIPv4
	}

	return ips, nil
}

// changedSeedAttrs lists the seed content attributes that differ. Empty and unset are the same seed.
func changedSeedAttrs(plan, state bootstrapModel) []string {
	attrs := []struct {
		name        string
		plan, state types.String
	}{
		{"user_data", plan.UserData, state.UserData},
		{"meta_data", plan.MetaData, state.MetaData},
		{"network_config", plan.NetworkConfig, state.NetworkConfig},
		{"vendor_data", plan.VendorData, state.VendorData},
	}

	var changed []string

	for _, a := range attrs {
		if a.plan.ValueString() != a.state.ValueString() {
			changed = append(changed, a.name)
		}
	}

	return changed
}

// Update stores changes that don't affect the VM. A seed content change can't reach it, so it fails rather than being
// stored silently. The error is raised here and not at plan time, because a plan error would also block replacing the
// VM, which is how such a change is applied.
func (r *bootstrapResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state bootstrapModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Returning before State.Set keeps the old state, so the next plan shows the change again.
	if changed := changedSeedAttrs(plan, state); len(changed) > 0 {
		vmID := state.VMID.ValueInt64()
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Update Cloud-init Bootstrap of VM %d", vmID),
			fmt.Sprintf("%s changed, but cloud-init applies the seed only on a VM's first boot, so the change can't reach "+
				"VM %d. No changes were saved. To bootstrap with the new content, replace the VM and this resource (apply "+
				"with -replace=<VM address> -replace=<this resource's address>). To keep the existing VM, add the changed "+
				"attributes to this resource's lifecycle.ignore_changes.", strings.Join(changed, ", "), vmID))

		return
	}

	plan.NodeName = state.NodeName

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the seed ISO a failed bootstrap left (seed_file_id); after a successful one there is nothing to undo.
// It runs before the VM is destroyed, because this resource depends on the VM. seed_file_id is state, not private
// state, because Terraform doesn't pass a tainted object's private state to its Delete.
func (r *bootstrapResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state bootstrapModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	volID := state.SeedFileID.ValueString()
	if volID == "" {
		return
	}

	vmID := state.VMID.ValueInt64()
	iface := state.Interface.ValueString()
	datastore, _, _ := strings.Cut(volID, ":")

	nodeName, vm, err := r.findVM(ctx, vmID, state.NodeName.ValueString())
	if errors.Is(err, cluster.ErrVMDoesNotExist) {
		// The ISO outlives the VM, so it can still be deleted from the node the VM was last seen on.
		nodeName = state.NodeName.ValueString()
		if nodeName == "" {
			warnSeedLeft(&resp.Diagnostics, state)

			return
		}

		vm, err = nil, nil
	}

	if err == nil {
		err = r.removeSeedISO(ctx, vm, r.client.Node(nodeName).Storage(datastore), iface, volID, &resp.Diagnostics)
	}

	if err != nil {
		if nodeName == "" {
			nodeName = state.NodeName.ValueString()
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Clean Up Seed ISO of VM %d", vmID),
			fmt.Sprintf("%s. The state is kept, so applying again retries. To clean up by hand instead, %s, then apply "+
				"again.", err.Error(), manualCleanup(nodeName, vmID, iface, volID)))
	}
}
