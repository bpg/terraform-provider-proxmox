//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

const vmAddr = "proxmox_virtual_environment_vm.test"

const bootstrapAddr = "proxmox_cloudinit_bootstrap.test"

// otherSuffix names an unrelated ISO the bootstrap must not touch.
const otherSuffix = "-other.iso"

// The guest needs network access to install qemu-guest-agent; the bridge defaults to vmbr0.
func networkDevices(n int) string {
	bridge := os.Getenv("PROXMOX_VE_ACC_CLOUDINIT_BRIDGE")
	if bridge == "" {
		bridge = "vmbr0"
	}

	vlan := os.Getenv("PROXMOX_VE_ACC_CLOUDINIT_VLAN_ID")
	if vlan == "" {
		vlan = "0"
	}

	return strings.Repeat(fmt.Sprintf("network_device {\n  bridge = %q\n  vlan_id = %s\n}\n", bridge, vlan), n)
}

// vmOpts varies the test VM; the zero value is a VM ready to be bootstrapped on ide0.
type vmOpts struct {
	noAgent bool
	cdrom   string // CD-ROM slot, ide0 if empty
	nics    int    // 1 if 0
	// replacedBy is added to the VM's replace_triggered_by.
	replacedBy string
	// initialization gives the VM a Proxmox cloud-init drive on ide2.
	initialization bool
}

// vmConfig renders a stopped VM with an empty ide0 CD-ROM, ready to be bootstrapped.
func vmConfig(te *test.Environment, nodeName string) string {
	return vmConfigWith(te, nodeName, vmOpts{})
}

func vmConfigWith(te *test.Environment, nodeName string, o vmOpts) string {
	if o.cdrom == "" {
		o.cdrom = "ide0"
	}

	if o.nics == 0 {
		o.nics = 1
	}

	agent := "agent {\n    enabled = true\n  }"
	if o.noAgent {
		agent = ""
	}

	if o.initialization {
		agent += "\n\n  initialization {\n    datastore_id = \"local\"\n    interface    = \"ide2\"\n  }"
	}

	lifecycle := "ignore_changes = [started]"
	if o.replacedBy != "" {
		lifecycle += "\n    replace_triggered_by = [" + o.replacedBy + "]"
	}

	return te.RenderConfig(`
resource "proxmox_virtual_environment_download_file" "image" {
  content_type        = "import"
  datastore_id        = "local"
  node_name           = "{{.NodeName}}"
  url                 = "{{.CloudImagesServer}}/noble/current/noble-server-cloudimg-amd64.img"
  file_name           = "{{.TestName}}.qcow2" # Ubuntu .img files are qcow2; .raw would not boot
  overwrite_unmanaged = true
}

resource "proxmox_virtual_environment_vm" "test" {
  node_name       = "` + nodeName + `"
  name            = "{{.TestName}}"
  machine         = "q35"
  started         = false
  migrate         = true
  stop_on_destroy = true

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = "local"
    import_from  = proxmox_virtual_environment_download_file.image.id
    interface    = "virtio0"
    size         = 8
  }

  cdrom {
    interface = "` + o.cdrom + `"
    file_id   = "none"
  }

  ` + networkDevices(o.nics) + `

  ` + agent + `

  lifecycle {
    ` + lifecycle + `
  }
}
`)
}

const userData = `<<-EOT
    #cloud-config
    packages: [qemu-guest-agent]
    runcmd:
      - [systemctl, enable, --now, qemu-guest-agent]
      %s
  EOT`

// networkConfig names the NIC "lan0" by MAC, so the happy path can check the name in ipv4_addresses.
const networkConfig = `network_config = yamlencode({
    version = 1
    config = [{
      type        = "physical"
      name        = "lan0"
      mac_address = lower(proxmox_virtual_environment_vm.test.network_device[0].mac_address)
      subnets     = [{ type = "dhcp" }]
    }]
  })
  meta_data = "instance-id: bootstrap-test\nlocal-hostname: bootstrap-test\n"`

// bootstrapConfig renders the bootstrap; extra is added verbatim, timeout 0 keeps the default.
func bootstrapConfig(extraRuncmd string, timeout int, extra string) string {
	t := extra
	if timeout != 0 {
		t += fmt.Sprintf("\n  timeout = %d", timeout)
	}

	return fmt.Sprintf(`
resource "proxmox_cloudinit_bootstrap" "test" {
  vm_id        = proxmox_virtual_environment_vm.test.vm_id
  datastore_id = "local"
  user_data    = %s
  %s
}
`, fmt.Sprintf(userData, extraRuncmd), t)
}

// seedName is the default seed ISO name for a VM.
func seedName(vmID int) string {
	return fmt.Sprintf("cloudinit-bootstrap-%d.iso", vmID)
}

func vmID(s *terraform.State) (int, error) {
	rs, ok := s.RootModule().Resources[vmAddr]
	if !ok {
		return 0, fmt.Errorf("%s not found in state", vmAddr)
	}

	id, err := strconv.Atoi(rs.Primary.Attributes["vm_id"])
	if err != nil {
		return 0, fmt.Errorf("invalid vm_id: %w", err)
	}

	return id, nil
}

// seedISOs lists the ISOs named fileName on a node's local datastore.
func seedISOs(ctx context.Context, te *test.Environment, nodeName, fileName string) ([]string, error) {
	nodes := te.NodeClient()
	nodes.NodeName = nodeName
	ds := nodes.Storage("local")

	contentType := storage.ContentTypeISO

	files, err := ds.ListDatastoreFiles(ctx, &contentType)
	if err != nil {
		return nil, err
	}

	var out []string

	for _, f := range files {
		if f.VolumeID == "local:iso/"+fileName {
			out = append(out, f.VolumeID)
		}
	}

	return out, nil
}

// checkCleanedUp verifies through the API that the VM is running with an empty drive and its seed ISO is gone.
func checkCleanedUp(te *test.Environment, nodeName string) resource.TestCheckFunc {
	return checkCleanedUpOn(te, nodeName, "ide0", "")
}

// checkCleanedUpOn is checkCleanedUp for a custom CD-ROM slot and seed file name (the default name if empty).
func checkCleanedUpOn(te *test.Environment, nodeName, iface, fileName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		id, err := vmID(s)
		if err != nil {
			return err
		}

		nodes := te.NodeClient()
		nodes.NodeName = nodeName
		vm := nodes.VM(id)

		status, err := vm.GetVMStatus(ctx)
		if err != nil {
			return err
		}

		if status.Status != "running" {
			return fmt.Errorf("VM %d is %s, want running", id, status.Status)
		}

		cfg, err := vm.GetVM(ctx)
		if err != nil {
			return err
		}

		if dev := cfg.StorageDevices[iface]; dev == nil || dev.FileVolume != "none" {
			return fmt.Errorf("%s of VM %d = %+v, want none", iface, id, dev)
		}

		name := fileName
		if name == "" {
			name = seedName(id)
		}

		isos, err := seedISOs(ctx, te, nodeName, name)
		if err != nil {
			return err
		}

		if len(isos) > 0 {
			return fmt.Errorf("seed ISOs left for VM %d: %v", id, isos)
		}

		return nil
	}
}

// testName is the subtest's VM name.
func testName(te *test.Environment) string {
	// RenderConfig prepends the provider block; the name is the last line.
	rendered := strings.Split(strings.TrimSpace(te.RenderConfig("{{.TestName}}")), "\n")

	return strings.TrimSpace(rendered[len(rendered)-1])
}

// testVM looks up the test VM by name, for steps that need it before it is in state or after a failed apply.
func testVM(t *testing.T, te *test.Environment) (string, *vms.Client) {
	t.Helper()

	name := testName(te)

	// A just-created VM shows up in the cluster resources without its name until pvestatd's next update.
	for range 10 {
		list, err := te.ClusterClient().GetClusterResourcesVM(context.Background())
		require.NoError(t, err)

		for _, v := range list {
			if v.Name != name {
				continue
			}

			nodes := te.NodeClient()
			nodes.NodeName = v.NodeName
			vm := nodes.VM(v.VMID)

			// Cluster resources can list a reused VM ID under the name of an earlier VM for a while, so check the config.
			cfg, err := vm.GetVM(context.Background())
			if err == nil && cfg.Name != nil && *cfg.Name == name {
				return v.NodeName, vm
			}
		}

		time.Sleep(3 * time.Second)
	}

	require.FailNow(t, "test VM not found: "+name)

	return "", nil
}

// uploadStaleISO leaves an ISO named fileName that no state knows about, like a killed run would.
func uploadStaleISO(t *testing.T, te *test.Environment, fileName string) {
	t.Helper()

	nodeName, _ := testVM(t, te)

	f, err := os.CreateTemp(t.TempDir(), "stale-*.iso")
	require.NoError(t, err)

	_, err = f.WriteString("stale")
	require.NoError(t, err)
	_, err = f.Seek(0, 0)
	require.NoError(t, err)

	nodes := te.NodeClient()
	nodes.NodeName = nodeName

	_, err = nodes.Storage("local").APIUpload(context.Background(), &api.FileUploadRequest{
		ContentType: storage.ContentTypeISO,
		FileName:    fileName,
		File:        f,
	}, t.TempDir())
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

// checkISOExists verifies an ISO named fileName is still on the node's local datastore.
func checkISOExists(te *test.Environment, nodeName, fileName string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		isos, err := seedISOs(context.Background(), te, nodeName, fileName)
		if err != nil {
			return err
		}

		if len(isos) == 0 {
			return fmt.Errorf("%s was deleted, but it isn't the seed ISO", fileName)
		}

		return nil
	}
}

// deleteISO removes an ISO the test uploaded itself.
func deleteISO(t *testing.T, te *test.Environment, fileName string) {
	t.Helper()

	// A delete of a missing file fails as a task, so skip files a step already removed.
	isos, err := seedISOs(context.Background(), te, te.NodeName, fileName)
	if err != nil {
		t.Errorf("unable to list ISOs to delete %s: %v", fileName, err)

		return
	}

	if len(isos) == 0 {
		return
	}

	nodes := te.NodeClient()
	nodes.NodeName = te.NodeName

	err = nodes.Storage("local").DeleteDatastoreFile(context.Background(), "local:iso/"+fileName).Err()
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		t.Errorf("unable to delete %s: %v", fileName, err)
	}
}

// requireSeedAttached checks a failed bootstrap left its seed ISO attached for debugging.
func requireSeedAttached(t *testing.T, te *test.Environment) {
	t.Helper()

	_, vm := testVM(t, te)

	cfg, err := vm.GetVM(context.Background())
	require.NoError(t, err)

	dev := cfg.StorageDevices["ide0"]
	require.NotNil(t, dev)
	require.Equal(t, "local:iso/"+seedName(vm.VMID), dev.FileVolume, "ide0 must hold the seed ISO")
}

// requireSeedRemoved checks the VM's drive is empty and the seed ISO is gone.
func requireSeedRemoved(t *testing.T, te *test.Environment) {
	t.Helper()

	_, vm := testVM(t, te)

	cfg, err := vm.GetVM(context.Background())
	require.NoError(t, err)

	dev := cfg.StorageDevices["ide0"]
	require.NotNil(t, dev)
	require.Equal(t, "none", dev.FileVolume)

	requireNoSeed(t, te)
}

// checkGuestFile reads a file in the guest through the agent, to prove which seed the VM booted with.
func checkGuestFile(te *test.Environment, nodeName, path, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id, err := vmID(s)
		if err != nil {
			return err
		}

		nodes := te.NodeClient()
		nodes.NodeName = nodeName

		got, err := nodes.VM(id).GetFileFromAgent(context.Background(), path)
		if err != nil {
			return fmt.Errorf("unable to read %s in VM %d: %w", path, id, err)
		}

		if strings.TrimSpace(got) != want {
			return fmt.Errorf("%s in VM %d = %q, want %q", path, id, got, want)
		}

		return nil
	}
}

// requireNoSeed checks the test VM's default seed ISO isn't on the datastore.
func requireNoSeed(t *testing.T, te *test.Environment) {
	t.Helper()

	nodeName, vm := testVM(t, te)

	isos, err := seedISOs(context.Background(), te, nodeName, seedName(vm.VMID))
	require.NoError(t, err)
	require.Empty(t, isos, "seed ISO left on the datastore")
}

// testSeedName is the default seed ISO name of the test VM, for steps that run before it is in state.
func testSeedName(t *testing.T, te *test.Environment) string {
	t.Helper()

	_, vm := testVM(t, te)

	return seedName(vm.VMID)
}

// checkSeedRecorded checks file_name defaults to the VM's seed name and seed_file_id is that file, or unset.
func checkSeedRecorded(recorded bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id, err := vmID(s)
		if err != nil {
			return err
		}

		checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(bootstrapAddr, "file_name", seedName(id))}
		if recorded {
			checks = append(checks, resource.TestCheckResourceAttr(bootstrapAddr, "seed_file_id", "local:iso/"+seedName(id)))
		} else {
			checks = append(checks, resource.TestCheckNoResourceAttr(bootstrapAddr, "seed_file_id"))
		}

		return resource.ComposeTestCheckFunc(checks...)(s)
	}
}

func TestAccResourceCloudInitBootstrap(t *testing.T) {
	type testCase struct {
		name  string
		steps func(t *testing.T, te *test.Environment) []resource.TestStep
	}

	tests := []testCase{
		{"bootstrap ejects and deletes the seed ISO", func(t *testing.T, te *test.Environment) []resource.TestStep {
			t.Helper()

			cfg := vmConfig(te, te.NodeName) + bootstrapConfig("", 0, networkConfig)

			return []resource.TestStep{{
				Config:      vmConfig(te, te.NodeName) + bootstrapConfig("", -1, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			}, {
				Config:      vmConfig(te, te.NodeName) + bootstrapConfig("", 0, `file_name = "seeds/a.iso"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must\s+end\s+in\s+.iso\s+and\s+use\s+only\s+letters`),
			}, {
				// Create the VM alone first, so a stale seed ISO can be planted before the bootstrap.
				Config: vmConfig(te, te.NodeName),
			}, {
				// A file at the seed name that this resource didn't upload fails the bootstrap and is left alone.
				PreConfig: func() {
					name := testSeedName(t, te)
					uploadStaleISO(t, te, name)
					t.Cleanup(func() { deleteISO(t, te, name) })
				},
				Config:      cfg,
				ExpectError: regexp.MustCompile(`already\s+exists\s+in\s+the\s+datastore\s+and\s+this\s+resource\s+didn't\s+upload\s+it`),
			}, {
				// The failure stored nothing (a plain create, not a tainted replace), and the stale file survived it.
				PreConfig: func() {
					name := testSeedName(t, te)
					require.NoError(t, checkISOExists(te, te.NodeName, name)(nil))
					deleteISO(t, te, name)
					uploadStaleISO(t, te, testName(te)+otherSuffix)
					t.Cleanup(func() { deleteISO(t, te, testName(te)+otherSuffix) })
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					checkSeedRecorded(false),
					checkISOExists(te, te.NodeName, testName(te)+otherSuffix),
					resource.TestCheckResourceAttr(bootstrapAddr, "interface", "ide0"),
					resource.TestCheckResourceAttr(bootstrapAddr, "timeout", "1800"),
					resource.TestCheckResourceAttr(bootstrapAddr, "node_name", te.NodeName),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
				),
			}, {
				// A refresh must not report changes: the VM picks up its IPs and the bootstrap keeps its state.
				Config:   cfg,
				PlanOnly: true,
			}, {
				// A content change after the bootstrap can't reach the VM, so the update fails.
				Config: vmConfig(te, te.NodeName) + bootstrapConfig(`- [echo, changed]`, 0, networkConfig),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`user_data\s+changed,\s+but\s+cloud-init\s+applies\s+the\s+seed\s+only`),
			}, {
				// The failed update saved nothing: the original config plans no changes, and the VM is untouched.
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
				),
			}, {
				// ignore_changes keeps the bootstrapped VM as it is.
				Config: vmConfig(te, te.NodeName) + bootstrapConfig(`- [echo, changed]`, 0,
					networkConfig+"\n  lifecycle {\n    ignore_changes = [user_data]\n  }"),
				PlanOnly: true,
			}, {
				// Other attributes update in place and leave the VM alone.
				Config: vmConfig(te, te.NodeName) + bootstrapConfig("", 600, networkConfig),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(vmAddr, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					resource.TestCheckResourceAttr(bootstrapAddr, "timeout", "600"),
					resource.TestCheckResourceAttr(bootstrapAddr, "node_name", te.NodeName),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
				),
			}, {
				// A refresh while the VM is stopped keeps the last known addresses.
				PreConfig:    func() { _, vm := testVM(t, te); require.NoError(t, vm.StopVM(context.Background()).Err()) },
				RefreshState: true,
				Check:        resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
			}, {
				// A VM deleted outside Terraform is created again, and so must its bootstrap be.
				PreConfig: func() {
					_, vm := testVM(t, te)
					require.NoError(t, vm.DeleteVM(context.Background(), true, true).Err())
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(vmAddr, plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
				),
			}, {
				// After a successful bootstrap destroy deletes nothing, not even a file at the seed name.
				PreConfig: func() {
					name := testSeedName(t, te)
					uploadStaleISO(t, te, name)
					t.Cleanup(func() { deleteISO(t, te, name) })
				},
				Config: vmConfig(te, te.NodeName),
				Check: func(s *terraform.State) error {
					id, err := vmID(s)
					if err != nil {
						return err
					}

					return checkISOExists(te, te.NodeName, seedName(id))(s)
				},
			}}
		}},
		{"cloud-init errors fail and a retry removes the seed ISO", func(t *testing.T, te *test.Environment) []resource.TestStep {
			t.Helper()

			cfg := vmConfig(te, te.NodeName) + bootstrapConfig(`- ["false"]`, 0, "")

			return []resource.TestStep{{
				// Without the guest agent the wait could never finish, so the bootstrap fails before uploading anything.
				Config:      vmConfigWith(te, te.NodeName, vmOpts{noAgent: true}) + bootstrapConfig(`- ["false"]`, 0, ""),
				ExpectError: regexp.MustCompile(`guest\s+agent\s+is\s+disabled\s+on\s+VM`),
			}, {
				// Attaching the seed to a slot without an empty CD-ROM would replace what is there, so that fails too.
				PreConfig:   func() { requireNoSeed(t, te) },
				Config:      vmConfigWith(te, te.NodeName, vmOpts{cdrom: "ide2"}) + bootstrapConfig(`- ["false"]`, 0, ""),
				ExpectError: regexp.MustCompile(`ide0\s+on\s+VM\s+\d+\s+must\s+be\s+an\s+empty\s+CD-ROM,\s+but\s+it\s+is\s+not\s+configured`),
			}, {
				// The failed prechecks stored nothing, so this is a plain create.
				PreConfig: func() { requireNoSeed(t, te) },
				Config:    cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionCreate),
					},
				},
				ExpectError: regexp.MustCompile(`cloud-init\s+finished\s+with\s+errors`),
			}, {
				// The failed bootstrap recorded the seed ISO it left attached, and refresh keeps it while it's there. The
				// tainted bootstrap still plans a replace.
				PreConfig:          func() { requireSeedAttached(t, te) },
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: checkSeedRecorded(true),
			}, {
				// Retrying on the same VM can't help: destroying the tainted bootstrap removes exactly the recorded seed
				// ISO, then the new one says to replace the VM.
				PreConfig: func() {
					uploadStaleISO(t, te, testName(te)+otherSuffix)
					t.Cleanup(func() { deleteISO(t, te, testName(te)+otherSuffix) })
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectError: regexp.MustCompile(`already\s+running,\s+and\s+cloud-init\s+applies\s+the\s+seed\s+only`),
			}, {
				PreConfig: func() { requireSeedRemoved(t, te) },
				Config:    vmConfig(te, te.NodeName),
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					checkISOExists(te, te.NodeName, testName(te)+otherSuffix),
				),
			}}
		}},
		{"PVE cloud-init drive fails before the upload", func(t *testing.T, te *test.Environment) []resource.TestStep {
			t.Helper()

			return []resource.TestStep{{
				// A clone of a template with a cloud-init drive would boot with two NoCloud seeds.
				Config:      vmConfigWith(te, te.NodeName, vmOpts{initialization: true}) + bootstrapConfig("", 0, ""),
				ExpectError: regexp.MustCompile(`has\s+a\s+Proxmox\s+cloud-init\s+drive\s+on\s+ide2`),
			}, {
				PreConfig: func() { requireNoSeed(t, te) },
				Config:    vmConfigWith(te, te.NodeName, vmOpts{initialization: true}),
			}}
		}},
		{"replacement bootstraps the new VM with the new content", func(t *testing.T, te *test.Environment) []resource.TestStep {
			t.Helper()

			fileName := testName(te) + "-seed.iso"

			// The seed's version drives the VM's replacement, the way a module rotates a key.
			cfg := func(version string) string {
				return vmConfigWith(te, te.NodeName, vmOpts{cdrom: "sata1", nics: 2, replacedBy: "terraform_data.seed"}) +
					fmt.Sprintf(`
resource "terraform_data" "seed" {
  input = %[1]q
}

resource "proxmox_cloudinit_bootstrap" "test" {
  vm_id        = proxmox_virtual_environment_vm.test.vm_id
  datastore_id = "local"
  file_name    = %[2]q
  interface    = "sata1"
  user_data    = <<-EOT
    #cloud-config
    packages: [qemu-guest-agent]
    write_files:
      - path: /var/lib/bootstrap-test/user-data
        content: %[1]s
    runcmd:
      - [systemctl, enable, --now, qemu-guest-agent]
  EOT
  # user_data's keys replace vendor_data's, so the marker uses a key user_data doesn't set.
  vendor_data = <<-EOT
    #cloud-config
    bootcmd:
      - [sh, -c, "mkdir -p /var/lib/bootstrap-test && echo vendor > /var/lib/bootstrap-test/vendor-data"]
  EOT
  network_config = yamlencode({
    version = 1
    config = [for i, nic in proxmox_virtual_environment_vm.test.network_device : {
      type        = "physical"
      name        = "lan${i}"
      mac_address = lower(nic.mac_address)
      subnets     = [{ type = "dhcp" }]
    }]
  })

  lifecycle {
    replace_triggered_by = [proxmox_virtual_environment_vm.test.id]
  }
}
`, version, fileName)
			}

			check := func(version string) resource.TestCheckFunc {
				return resource.ComposeTestCheckFunc(
					checkCleanedUpOn(te, te.NodeName, "sata1", fileName),
					checkGuestFile(te, te.NodeName, "/var/lib/bootstrap-test/user-data", version),
					checkGuestFile(te, te.NodeName, "/var/lib/bootstrap-test/vendor-data", "vendor"),
					resource.TestCheckResourceAttr(bootstrapAddr, "file_name", fileName),
					resource.TestCheckResourceAttr(bootstrapAddr, "interface", "sata1"),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan1.#", regexp.MustCompile(`^[1-9]`)),
				)
			}

			return []resource.TestStep{{
				Config: cfg("v1"),
				Check:  check("v1"),
			}, {
				Config: cfg("v2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(vmAddr, plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: check("v2"),
			}, {
				Config:   cfg("v2"),
				PlanOnly: true,
			}}
		}},
	}

	if os.Getenv("PROXMOX_VE_ACC_NODE_2_NAME") != "" {
		tests = append(tests, testCase{
			"migrating a bootstrapped VM does not re-run the bootstrap",
			func(t *testing.T, te *test.Environment) []resource.TestStep {
				t.Helper()

				return []resource.TestStep{{
					Config: vmConfig(te, te.NodeName) + bootstrapConfig("", 0, ""),
					Check:  checkCleanedUp(te, te.NodeName),
				}, {
					Config: vmConfig(te, te.Node2Name) + bootstrapConfig("", 0, ""),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(vmAddr, plancheck.ResourceActionUpdate),
							plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionNoop),
						},
					},
					Check: checkCleanedUp(te, te.Node2Name),
				}, {
					// Refresh records the node the VM moved to, and that doesn't show as a change.
					RefreshState: true,
					Check:        resource.TestCheckResourceAttr(bootstrapAddr, "node_name", te.Node2Name),
				}, {
					Config:   vmConfig(te, te.Node2Name) + bootstrapConfig("", 0, ""),
					PlanOnly: true,
				}}
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// One environment per subtest gives each its own test name, so their image files don't collide.
			te := test.InitEnvironment(t)

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps(t, te),
			})
		})
	}
}
