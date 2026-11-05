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

// The guest needs network access to install qemu-guest-agent; the bridge defaults to vmbr0.
func networkDevice() string {
	bridge := os.Getenv("PROXMOX_VE_ACC_CLOUDINIT_BRIDGE")
	if bridge == "" {
		bridge = "vmbr0"
	}

	vlan := os.Getenv("PROXMOX_VE_ACC_CLOUDINIT_VLAN_ID")
	if vlan == "" {
		vlan = "0"
	}

	return fmt.Sprintf("network_device {\n  bridge = %q\n  vlan_id = %s\n}", bridge, vlan)
}

// vmConfig renders a stopped VM with an empty ide0 CD-ROM, ready to be bootstrapped.
func vmConfig(te *test.Environment, nodeName string) string {
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
    interface = "ide0"
    file_id   = "none"
  }

  ` + networkDevice() + `

  agent {
    enabled = true
  }

  lifecycle {
    ignore_changes = [started]
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

// seedISOs lists the seed ISOs of a VM on a node's local datastore.
func seedISOs(ctx context.Context, te *test.Environment, nodeName string, id int) ([]string, error) {
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
		if strings.HasPrefix(f.VolumeID, fmt.Sprintf("local:iso/cloudinit-%d-", id)) {
			out = append(out, f.VolumeID)
		}
	}

	return out, nil
}

// checkCleanedUp verifies through the API that the VM is running with an empty drive and its seed ISO is gone.
func checkCleanedUp(te *test.Environment, nodeName string) resource.TestCheckFunc {
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

		if dev := cfg.StorageDevices["ide0"]; dev == nil || dev.FileVolume != "none" {
			return fmt.Errorf("ide0 of VM %d = %+v, want none", id, dev)
		}

		isos, err := seedISOs(ctx, te, nodeName, id)
		if err != nil {
			return err
		}

		if len(isos) > 0 {
			return fmt.Errorf("seed ISOs left for VM %d: %v", id, isos)
		}

		return nil
	}
}

// testVM looks up the test VM by name, for steps that need it before it is in state or after a failed apply.
func testVM(t *testing.T, te *test.Environment) (string, *vms.Client) {
	t.Helper()

	// RenderConfig prepends the provider block; the name is the last line.
	rendered := strings.Split(strings.TrimSpace(te.RenderConfig("{{.TestName}}")), "\n")
	name := strings.TrimSpace(rendered[len(rendered)-1])

	// A just-created VM shows up in the cluster resources without its name until pvestatd's next update.
	for range 10 {
		list, err := te.ClusterClient().GetClusterResourcesVM(context.Background())
		require.NoError(t, err)

		for _, v := range list {
			if v.Name == name {
				nodes := te.NodeClient()
				nodes.NodeName = v.NodeName

				return v.NodeName, nodes.VM(v.VMID)
			}
		}

		time.Sleep(3 * time.Second)
	}

	require.FailNow(t, "test VM not found: "+name)

	return "", nil
}

// uploadStaleISO leaves a seed ISO for the test VM that no state knows about, like a killed run would.
func uploadStaleISO(t *testing.T, te *test.Environment) {
	t.Helper()

	nodeName, vm := testVM(t, te)

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
		FileName:    fmt.Sprintf("cloudinit-%d-stale.iso", vm.VMID),
		File:        f,
	}, t.TempDir())
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

// requireSeedAttached checks a failed bootstrap left its seed ISO attached for debugging.
func requireSeedAttached(t *testing.T, te *test.Environment) {
	t.Helper()

	_, vm := testVM(t, te)

	cfg, err := vm.GetVM(context.Background())
	require.NoError(t, err)

	dev := cfg.StorageDevices["ide0"]
	require.NotNil(t, dev)
	require.True(t, strings.HasPrefix(dev.FileVolume, fmt.Sprintf("local:iso/cloudinit-%d-", vm.VMID)),
		"ide0 = %s, want the seed ISO", dev.FileVolume)
}

// requireSeedRemoved checks the VM's drive is empty and none of its seed ISOs is left.
func requireSeedRemoved(t *testing.T, te *test.Environment) {
	t.Helper()

	nodeName, vm := testVM(t, te)

	cfg, err := vm.GetVM(context.Background())
	require.NoError(t, err)

	dev := cfg.StorageDevices["ide0"]
	require.NotNil(t, dev)
	require.Equal(t, "none", dev.FileVolume)

	isos, err := seedISOs(context.Background(), te, nodeName, vm.VMID)
	require.NoError(t, err)
	require.Empty(t, isos)
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
				// Create the VM alone first, so a stale seed ISO can be planted before the bootstrap.
				Config: vmConfig(te, te.NodeName),
			}, {
				// The bootstrap must also remove a seed ISO that no state knows about.
				PreConfig: func() { uploadStaleISO(t, te) },
				Config:    cfg,
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
					resource.TestCheckResourceAttr(bootstrapAddr, "interface", "ide0"),
					resource.TestCheckResourceAttr(bootstrapAddr, "timeout", "1800"),
					resource.TestMatchResourceAttr(bootstrapAddr, "ipv4_addresses.lan0.#", regexp.MustCompile(`^[1-9]`)),
				),
			}, {
				// A refresh must not report changes: the VM picks up its IPs and the bootstrap keeps its state.
				Config:   cfg,
				PlanOnly: true,
			}, {
				// Content changes after bootstrap are stored in place and must not touch the running VM.
				Config: vmConfig(te, te.NodeName) + bootstrapConfig(`- [echo, changed]`, 0, networkConfig),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(bootstrapAddr, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					checkCleanedUp(te, te.NodeName),
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
			}}
		}},
		{"cloud-init errors fail and a retry removes the seed ISO", func(t *testing.T, te *test.Environment) []resource.TestStep {
			t.Helper()

			return []resource.TestStep{{
				Config:      vmConfig(te, te.NodeName) + bootstrapConfig(`- ["false"]`, 0, ""),
				ExpectError: regexp.MustCompile(`cloud-init finished with errors`),
			}, {
				// Retrying on the same VM can't help: it removes the seed ISO and says to replace the VM.
				PreConfig:   func() { requireSeedAttached(t, te) },
				Config:      vmConfig(te, te.NodeName) + bootstrapConfig(`- ["false"]`, 0, ""),
				ExpectError: regexp.MustCompile(`already\s+running\s+and\s+cloud-init\s+only\s+runs\s+on\s+first\s+boot`),
			}, {
				PreConfig: func() { requireSeedRemoved(t, te) },
				Config:    vmConfig(te, te.NodeName),
				Check:     checkCleanedUp(te, te.NodeName),
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
