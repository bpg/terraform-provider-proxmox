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
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

type bootstrapModel struct {
	ID            types.String `tfsdk:"id"`
	VMID          types.Int64  `tfsdk:"vm_id"`
	DatastoreID   types.String `tfsdk:"datastore_id"`
	FileName      types.String `tfsdk:"file_name"`
	Interface     types.String `tfsdk:"interface"`
	UserData      types.String `tfsdk:"user_data"`
	MetaData      types.String `tfsdk:"meta_data"`
	NetworkConfig types.String `tfsdk:"network_config"`
	VendorData    types.String `tfsdk:"vendor_data"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	IPv4Addresses types.Map    `tfsdk:"ipv4_addresses"`
}

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

// contentString is a seed file attribute. Changing it is stored in place: cloud-init only runs on first boot.
func contentString(desc string, required bool) schema.StringAttribute {
	return schema.StringAttribute{
		Description: desc,
		Required:    required,
		Optional:    !required,
		Sensitive:   true,
	}
}

func (r *bootstrapResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Bootstraps a VM with cloud-init through a temporary NoCloud seed ISO, without SSH.",
		MarkdownDescription: "Bootstraps a VM with cloud-init through a temporary " +
			"[NoCloud](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) seed ISO, without SSH.\n\n" +
			"Create the VM with `started = false`, `lifecycle { ignore_changes = [started] }` (or the next apply stops it " +
			"again) and an empty CD-ROM (`file_id = \"none\"`). On create this resource " +
			"builds the seed ISO, uploads it as `file_name` to the datastore on the node the VM is on, attaches it, starts the VM, waits for " +
			"cloud-init to finish, then ejects and deletes the ISO. Nothing is left on the node afterwards, so the VM can be " +
			"migrated freely.\n\n" +
			"Completion is read from `" + cloudInitResultPath + "` through the QEMU guest agent, which must be installed in " +
			"the guest (installing it from the same cloud-init config works). The wait survives guest reboots such as " +
			"`package_reboot_if_required`. If cloud-init reports errors or the wait times out, creation fails, the ISO " +
			"stays attached for debugging and the resource is tainted. cloud-init only runs on a VM's first boot, so " +
			"replace the VM to bootstrap it again. Applying again while the VM is running, e.g. after a failed or interrupted " +
			"apply, only removes the ISO and fails. Don't stop the VM and apply again: booting it a second time does not " +
			"run cloud-init's setup again, so the bootstrap could report success. Destroying the resource also removes " +
			"the ISO.\n\n" +
			"cloud-init only reports the errors of its last boot. An error before a reboot it triggers itself (e.g. " +
			"`package_reboot_if_required` after a failed package install) is not reported; check " +
			"`/var/log/cloud-init.log` in the guest if a bootstrapped VM looks incomplete.\n\n" +
			"After cloud-init finishes, creation waits up to 2 minutes for every network device " +
			"to get an IPv4 address, then warns if one still has none.\n\n" +
			"Use `lifecycle.replace_triggered_by` on the VM so a replaced VM is bootstrapped again. If the VM is deleted " +
			"outside Terraform, the bootstrap is removed from state on refresh, so it is created again with the VM. " +
			"Changing any attribute but `vm_id` later only updates state and does not affect the running VM.\n\n" +
			"Seed contents are stored in Terraform state, inside the guest under `/var/lib/cloud`, and in provider debug " +
			"logs (`TF_LOG=DEBUG` logs the uploaded ISO). While the ISO is on the datastore (during the bootstrap, or " +
			"after a failed one), anyone who can read ISOs there can read it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The VM ID.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vm_id": schema.Int64Attribute{
				Description:   "The VM to bootstrap. Its node is looked up at apply time.",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"datastore_id": schema.StringAttribute{
				Description: "The datastore on the VM's node to upload the temporary ISO to. It must allow `iso` content.",
				Required:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "The file name of the temporary seed ISO on `datastore_id`. It must end in `.iso` and use " +
					"only letters, digits, `.`, `_` and `-` (Proxmox would rename other characters). Use a name no other " +
					"ISO on the datastore has (across the cluster if it is shared): an existing file with this name is " +
					"replaced when the bootstrap runs and deleted afterwards, even if another VM has it attached.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
					stringvalidator.RegexMatches(isoFileNameRE, isoFileNameMsg),
				},
			},
			"interface": schema.StringAttribute{
				Description: "The CD-ROM interface to attach the ISO to. Defaults to `ide0` " +
					"(`ide2` is the `initialization` default).",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ide0"),
			},
			"user_data": contentString("The `user-data` file contents, e.g. a `#cloud-config` document.", true),
			"meta_data": contentString("The `meta-data` file contents. Defaults to `instance-id` and `local-hostname` set "+
				"to the VM name.", false),
			"network_config": contentString("The `network-config` file contents (network config version 1 or 2).", false),
			"vendor_data":    contentString("The `vendor-data` file contents.", false),
			"timeout": schema.Int64Attribute{
				Description: "How long to wait for cloud-init to finish, in seconds. Defaults to 1800.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1800),
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"ipv4_addresses": schema.MapAttribute{
				Description: "IPv4 addresses of the VM's network devices reported by the guest agent, keyed by " +
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

func seedVolumeID(datastore, fileName string) string {
	return datastore + ":" + storage.ContentTypeISO + "/" + fileName
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

	ips, complete, err := r.bootstrap(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Bootstrap VM %d", vmID), err.Error())

		// A (tainted) state makes destroy or the next replacement remove any seed ISO left attached.
		plan.IPv4Addresses = types.MapValueMust(ipListType(), map[string]attr.Value{})
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

		return
	}

	if !complete {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Incomplete IPv4 Addresses for VM %d", vmID),
			"cloud-init finished, but not every network device of the VM had an IPv4 address after 2 minutes "+
				"(e.g. DHCP still pending). ipv4_addresses is incomplete; later refreshes fill it in once every "+
				"device has one.")
	}

	plan.IPv4Addresses = ips

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// vmOnCurrentNode returns the VM's client on the node it runs on now, which changes when it migrates.
func (r *bootstrapResource) vmOnCurrentNode(ctx context.Context, vmID int64) (string, *vms.Client, error) {
	nodeName, err := r.client.Cluster().GetVMNodeName(ctx, int(vmID))
	if err != nil {
		return "", nil, fmt.Errorf("unable to find the node of VM %d: %w", vmID, err)
	}

	return *nodeName, r.client.Node(*nodeName).VM(int(vmID)), nil
}

// bootstrap returns the guest's IPv4 addresses and whether every network device has one.
func (r *bootstrapResource) bootstrap(ctx context.Context, m *bootstrapModel) (types.Map, bool, error) {
	none := types.MapNull(ipListType())
	vmID := m.VMID.ValueInt64()
	iface := m.Interface.ValueString()
	datastore := m.DatastoreID.ValueString()
	fileName := m.FileName.ValueString()

	nodeName, vm, err := r.vmOnCurrentNode(ctx, vmID)
	if err != nil {
		return none, false, err
	}

	datastoreClient := r.client.Node(nodeName).Storage(datastore)

	status, err := vm.GetVMStatus(ctx)
	if err != nil {
		return none, false, err
	}

	// Remove what an earlier, failed or interrupted attempt left behind.
	if err := r.removeSeedISO(ctx, vm, datastoreClient, iface, datastore, fileName); err != nil {
		return none, false, err
	}

	if status.Status == "running" {
		return none, false, fmt.Errorf("VM %d is already running and cloud-init only runs on first boot: "+
			"replace the VM to bootstrap it again", vmID)
	}

	cfg, err := vm.GetVM(ctx)
	if err != nil {
		return none, false, err
	}

	vmName := fmt.Sprintf("vm-%d", vmID)
	if cfg.Name != nil && *cfg.Name != "" {
		vmName = *cfg.Name
	}

	files := seedFiles(vmName, m.UserData.ValueString(), m.MetaData.ValueString(),
		m.NetworkConfig.ValueString(), m.VendorData.ValueString())
	volID := seedVolumeID(datastore, fileName)

	if err := uploadISO(ctx, r.client, datastoreClient, files, fileName); err != nil {
		return none, false, err
	}

	if err := setCDROM(ctx, vm, iface, volID); err != nil {
		return none, false, err
	}

	if err := vm.StartVM(ctx, int(m.Timeout.ValueInt64())).Err(); err != nil {
		return none, false, fmt.Errorf("unable to start VM %d: %w", vmID, err)
	}

	if err := waitForCloudInit(ctx, vm, time.Duration(m.Timeout.ValueInt64())*time.Second); err != nil {
		return none, false, fmt.Errorf("%w; the seed ISO is left attached for debugging. cloud-init only runs on "+
			"first boot, so replace the VM to bootstrap it again", err)
	}

	if err := r.removeSeedISO(ctx, vm, datastoreClient, iface, datastore, fileName); err != nil {
		return none, false, err
	}

	return waitForIPv4Addresses(ctx, vm)
}

// removeSeedISO ejects the seed ISO if the drive holds it and deletes it from the datastore. It only needs the
// configured name, not what an earlier attempt recorded, so it also cleans up after interrupted or failed runs.
func (r *bootstrapResource) removeSeedISO(
	ctx context.Context,
	vm *vms.Client,
	datastoreClient *storage.Client,
	iface, datastore, fileName string,
) error {
	volID := seedVolumeID(datastore, fileName)

	cfg, err := vm.GetVM(ctx)
	if err != nil {
		return err
	}

	if dev, ok := cfg.StorageDevices[iface]; ok && dev != nil && dev.FileVolume == volID {
		if err := setCDROM(ctx, vm, iface, "none"); err != nil {
			return err
		}
	}

	contentType := storage.ContentTypeISO

	// Listing first, because the file is usually already gone (before an upload, or destroying after a successful
	// bootstrap), and a delete of a missing file fails as a task rather than with a not-found error.
	files, err := datastoreClient.ListDatastoreFiles(ctx, &contentType)
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		// No such datastore (e.g. a typo in a tainted state), so there is nothing to clean up and destroy must not block.
		return nil
	}

	if err != nil {
		return fmt.Errorf("unable to list ISOs on %s: %w", datastore, err)
	}

	for _, f := range files {
		if f.VolumeID != volID {
			continue
		}

		err := datastoreClient.DeleteDatastoreFile(ctx, f.VolumeID).Err()
		if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
			return fmt.Errorf("unable to delete %s: %w", f.VolumeID, err)
		}
	}

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

// ipWaitTimeout bounds how long Create waits for DHCP on every NIC after cloud-init finished; the docs and the
// incomplete-addresses warning say "2 minutes".
const ipWaitTimeout = 2 * time.Minute

var errIncompleteIPv4 = errors.New("not every network device has an IPv4 address yet")

// waitForIPv4Addresses returns the guest's addresses once every NIC has one, or what it has after ipWaitTimeout,
// and whether they are complete.
func waitForIPv4Addresses(ctx context.Context, vm *vms.Client) (types.Map, bool, error) {
	waitCtx, cancel := context.WithTimeout(ctx, ipWaitTimeout)
	defer cancel()

	var (
		ips      types.Map
		complete bool
	)

	err := retry.NewPollOperation("guest IPv4 addresses", retry.WithBaseDelay(5*time.Second)).DoPoll(waitCtx, func() error {
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

// Read refreshes ipv4_addresses from the guest agent on the VM's current node, so DHCP changes show up.
// The bootstrap itself is one-time, so everything else is kept. So are the last known addresses while the VM is
// stopped, the agent can't answer, or a NIC has no IPv4 address yet (e.g. DHCP during a reboot).
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

	ips, err := r.currentIPv4Addresses(ctx, state.VMID.ValueInt64())

	switch {
	case errors.Is(err, cluster.ErrVMDoesNotExist):
		// The VM was deleted outside Terraform. It is planned as a create, which replace_triggered_by doesn't catch,
		// so the bootstrap has to be created again too.
		resp.State.RemoveResource(ctx)

		return
	case err != nil:
		tflog.Debug(ctx, "keeping the last known ipv4_addresses", map[string]any{"error": err.Error()})
	default:
		state.IPv4Addresses = ips
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bootstrapResource) currentIPv4Addresses(ctx context.Context, vmID int64) (types.Map, error) {
	none := types.MapNull(ipListType())

	_, vm, err := r.vmOnCurrentNode(ctx, vmID)
	if err != nil {
		return none, err
	}

	status, err := vm.GetVMStatus(ctx)
	if err != nil {
		return none, err
	}

	if status.Status != "running" {
		return none, fmt.Errorf("VM %d is %s", vmID, status.Status)
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

// Update only stores the new values: the bootstrap already ran, and cloud-init does not run again on the same VM.
func (r *bootstrapResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan bootstrapModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the seed ISO a failed bootstrap left; after a successful one there is nothing to undo.
// It runs before the VM is destroyed, because this resource depends on the VM.
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

	vmID := state.VMID.ValueInt64()
	datastore := state.DatastoreID.ValueString()

	nodeName, vm, err := r.vmOnCurrentNode(ctx, vmID)
	if errors.Is(err, cluster.ErrVMDoesNotExist) {
		return
	}

	if err == nil {
		err = r.removeSeedISO(ctx, vm, r.client.Node(nodeName).Storage(datastore), state.Interface.ValueString(),
			datastore, state.FileName.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Clean Up Seed ISO of VM %d", vmID), err.Error())
	}
}
