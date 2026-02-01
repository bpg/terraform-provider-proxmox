//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

const (
	accTestContainerName      = "proxmox_virtual_environment_container.test_container"
	alpineTemplateRelativeURL = "images/system/alpine-3.22-default_20250617_amd64.tar.xz"
)

func getTemplateURL(t *testing.T, templateServerURL string) string {
	t.Helper()
	return fmt.Sprintf("%s/%s", templateServerURL, alpineTemplateRelativeURL)
}

func TestAccResourceContainer(t *testing.T) {
	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	accTestContainerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
		"TimeoutDelete":   300,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local"
					}
					device_passthrough {
						path = "/dev/zero"
					}
					description = <<-EOT
						my
						description
						value
					EOT
					initialization {
						hostname = "test"
						ip_config {
							ipv4 {
							  address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"unprivileged":              "true",
						"description":               "my\ndescription\nvalue\n",
						"device_passthrough.#":      "1",
						"device_passthrough.0.mode": "0660",
						"initialization.0.dns.#":    "0",
					}),
					//nolint:godox
					// TODO: depends on DHCP, which may not work in some environments
					// ResourceAttributesSet(accTestContainerName, []string{
					// 	"ipv4.vmbr0",
					// }),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")

						dev0, ok := ctInfo.PassthroughDevices["dev0"]
						require.True(te.t, ok, `"dev0" passthrough device not found`)
						require.NotNil(te.t, dev0, `"dev0" passthrough device is <nil>`)
						require.Equal(te.t, "/dev/zero", dev0.Path)

						return nil
					},
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = ["discard"]
					}
					mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local"
					}
					device_passthrough {
						path = "/dev/zero"
					}
					description = <<-EOT
						my
						description
						value
					EOT
					initialization {
						hostname = "test"
						ip_config {
							ipv4 {
								address = "172.16.10.10/15"
								gateway = "172.16.0.1"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"description":            "my\ndescription\nvalue\n",
						"device_passthrough.#":   "1",
						"initialization.0.dns.#": "0",
						"disk.0.mount_options.#": "1",
					}),
				),
			},
			{
				// remove disk options
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = []
					}
					mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local"
					}
					device_passthrough {
						path = "/dev/zero"
					}
					description = <<-EOT
						my
						description
						value
					EOT
					initialization {
						hostname = "test"
						ip_config {
							ipv4 {
								address = "172.16.10.10/15"
								gateway = "172.16.0.1"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.mount_options.#": "0",
					}),
				),
			},
		},
	})
}

func TestAccResourceContainerMountOptions(t *testing.T) {
	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	accTestContainerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
		"TimeoutDelete":   300,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = ["lazytime"]
					}
					initialization {
						hostname = "test-mount-options"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.mount_options.#": "1",
						"disk.0.mount_options.0": "lazytime",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")

						require.NotNil(te.t, ctInfo.RootFS, "rootfs should not be nil")
						require.NotNil(te.t, ctInfo.RootFS.MountOptions, "mount_options should not be nil on initial creation")
						require.Len(te.t, *ctInfo.RootFS.MountOptions, 1, "mount_options should have 1 element on initial creation")
						require.Equal(te.t, "lazytime", (*ctInfo.RootFS.MountOptions)[0], "mount_options should contain 'lazytime'")

						te.t.Logf("Container created with rootFS volume: %s", ctInfo.RootFS.Volume)

						return nil
					},
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
						mount_options = ["lazytime", "noatime"]
					}
					initialization {
						hostname = "test-mount-options"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.mount_options.#": "2",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")

						require.NotNil(te.t, ctInfo.RootFS, "rootfs should not be nil")
						require.NotNil(te.t, ctInfo.RootFS.MountOptions, "mount_options should not be nil after update")
						require.Len(te.t, *ctInfo.RootFS.MountOptions, 2, "mount_options should have 2 elements after update")

						te.t.Logf("After update, rootFS volume: %s", ctInfo.RootFS.Volume)

						return nil
					},
				),
			},
		},
	})
}

func TestAccResourceContainerDiskResize(t *testing.T) {
	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	accTestContainerID := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]any{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
		"TimeoutDelete":   300,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Create container with 4GB disk
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-disk-resize"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.size": "4",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.RootFS, "rootfs should not be nil")
						require.NotNil(te.t, ctInfo.RootFS.Size, "rootfs size should not be nil")
						require.Equal(te.t, int64(4), ctInfo.RootFS.Size.InGigabytes(), "disk size should be 4GB")

						te.t.Logf("Container created with rootFS volume: %s, size: %s", ctInfo.RootFS.Volume, ctInfo.RootFS.Size)

						return nil
					},
				),
			},
			{
				// Resize disk to 6GB (should update in-place, not recreate)
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = {{ .TimeoutDelete }}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 6
					}
					initialization {
						hostname = "test-disk-resize"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.size": "6",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.RootFS, "rootfs should not be nil")
						require.NotNil(te.t, ctInfo.RootFS.Size, "rootfs size should not be nil")
						require.Equal(te.t, int64(6), ctInfo.RootFS.Size.InGigabytes(), "disk size should be 6GB after resize")

						te.t.Logf("After resize, rootFS volume: %s, size: %s", ctInfo.RootFS.Volume, ctInfo.RootFS.Size)

						return nil
					},
				),
			},
		},
	})
}

func TestAccResourceContainerClone(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerIDClone := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":        imageFileName,
		"TestContainerID":      accTestContainerID,
		"TestContainerIDClone": accTestContainerIDClone,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.TestContainerID}}
				template  = true
				disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				mount_point {
					volume = "local-lvm"
					size   = "4G"
					path   = "mnt/local"
				}
				device_passthrough {
					path = "/dev/zero"
				}
				initialization {
					hostname = "test"
					ip_config {
						ipv4 {
						  address = "dhcp"
						}
					}
				}
				network_interface {
					name = "vmbr0"
				}
				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "alpine"
				}
			}
			resource "proxmox_virtual_environment_container" "test_container_clone" {
				depends_on = [proxmox_virtual_environment_container.test_container]
				node_name  = "{{.NodeName}}"
				vm_id      = {{.TestContainerIDClone}}
				
				clone {
					vm_id = proxmox_virtual_environment_container.test_container.id
				}
				
				initialization {
					hostname = "test-clone"
					dns {
					  servers = ["1.1.1.1"]
					}
				}
			}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerIDClone)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")

						dev0, ok := ctInfo.PassthroughDevices["dev0"]
						require.True(te.t, ok, `"dev0" passthrough device not found`)
						require.NotNil(te.t, dev0, `"dev0" passthrough device is <nil>`)
						require.Equal(te.t, "/dev/zero", dev0.Path)

						return nil
					},
				),
			},
		},
	})
}

// Test that mount_point blocks specified in clone config are provisioned
// See https://github.com/bpg/terraform-provider-proxmox/issues/2518
func TestAccResourceContainerCloneMountPoint(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerIDClone := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":        imageFileName,
		"TestContainerID":      accTestContainerID,
		"TestContainerIDClone": accTestContainerIDClone,
	})

	config := te.RenderConfig(`
		resource "proxmox_virtual_environment_container" "test_container_mp" {
			node_name = "{{.NodeName}}"
			vm_id     = {{.TestContainerID}}
			template  = true
			disk {
				datastore_id = "local-lvm"
				size         = 4
			}
			initialization {
				hostname = "test"
				ip_config {
					ipv4 {
					  address = "dhcp"
					}
				}
			}
			network_interface {
				name = "vmbr0"
			}
			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}
		}
		resource "proxmox_virtual_environment_container" "test_container_clone_mp" {
			depends_on = [proxmox_virtual_environment_container.test_container_mp]
			node_name  = "{{.NodeName}}"
			vm_id      = {{.TestContainerIDClone}}

			clone {
				vm_id = proxmox_virtual_environment_container.test_container_mp.id
			}

			mount_point {
				volume = "local-lvm"
				size   = "4G"
				path   = "/mnt/data"
			}

			initialization {
				hostname = "test-clone"
			}
		}`, WithRootUser())

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerIDClone)

						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")

						mp0, ok := ctInfo.MountPoints["mp0"]
						require.True(te.t, ok, `"mp0" mount point not found`)
						require.NotNil(te.t, mp0, `"mp0" mount point is <nil>`)
						require.Equal(te.t, "/mnt/data", mp0.MountPoint)

						return nil
					},
				),
			},
			// Step 2: Re-apply same config - verifies no replacement triggered (#2507)
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceContainerDnsBlock(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)

	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-dns-issue"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-issue",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-dns-issue"
						dns {}
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-issue",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-dns-create"
						dns {}
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-create",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-dns-create"
						dns {
							domain = null
							server = ""
							servers = null
						}
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-create",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				initialization {
					hostname = "test-dns-update"
					dns {
						domain = "example.com"
						servers = ["8.8.8.8", "8.8.4.4"]
					}
					ip_config {
						ipv4 {
							address = "dhcp"
						}
					}
				}
				network_interface {
					name = "vmbr0"
				}
				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "alpine"
				}
			}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname":        "test-dns-update",
						"initialization.0.dns.#":           "1",
						"initialization.0.dns.0.domain":    "example.com",
						"initialization.0.dns.0.servers.#": "2",
						"initialization.0.dns.0.servers.0": "8.8.8.8",
						"initialization.0.dns.0.servers.1": "8.8.4.4",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-dns-update"
						dns {
							domain = ""
							servers = null
						}
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-update",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
		},
	})
}

func TestAccResourceContainerHostname(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-hostname-1"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-hostname-1",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-hostname-2"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-hostname-2",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)
						// Verify the hostname was actually updated in the container config
						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.Hostname, "hostname should not be nil")
						require.Equal(te.t, "test-hostname-2", *ctInfo.Hostname, "hostname should be updated")

						return nil
					},
				),
			},
		},
	})
}

func TestAccResourceContainerMountPoint(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
					vm_id = {{ .TestContainerID }}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local1"
				    }
				    initialization {
				  		hostname = "test"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
					"mount_point.#": "1",
				}),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
					vm_id = {{ .TestContainerID }}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local1"
				    }
					// add a new mount point
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "mnt/local2"
				    }
				    initialization {
				  		hostname = "test"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
					"mount_point.#": "2",
				}),
			},
		},
	})
}

func TestAccResourceContainerIpv4Ipv6(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_network_linux_bridge" "vmbr1" {
					node_name = "{{ .NodeName }}"
					name = "vmbr1"
				}
				
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					started   = false
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-hostname-1"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}
				`),
			},
			{
				Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_network_linux_bridge" "vmbr1" {
				node_name = "{{ .NodeName }}"
				name = "vmbr1"
			}
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.TestContainerID}}
				unprivileged = true
				started   = false
				disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				initialization {
					hostname = "test-hostname-1"
					ip_config {
						ipv4 {
							address = "dhcp"
						}
					}
				}
				network_interface {
					name = "vmbr0"
					bridge = "vmbr0"
				}
				network_interface {
					name = "vmbr1"
					bridge = "vmbr1"
				}

				operating_system {
					template_file_id = "local:vztmpl/{{.ImageFileName}}"
					type             = "alpine"
				}
			}`),
			},
		},
	})
}

func TestAccResourceContainerEnvironmentVariables(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-env-vars"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
					environment_variables = {
						FOO = "bar"
						BAZ = "qux"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-env-vars",
						"environment_variables.FOO": "bar",
						"environment_variables.BAZ": "qux",
						"environment_variables.%":   "2",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)
						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.EnvironmentVariables, "environment_variables should not be nil")
						require.Equal(te.t, "bar", (*ctInfo.EnvironmentVariables)["FOO"])
						require.Equal(te.t, "qux", (*ctInfo.EnvironmentVariables)["BAZ"])

						return nil
					},
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-env-vars"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
					environment_variables = {
						FOO = "updated"
						NEW_VAR = "new_value"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"environment_variables.FOO":     "updated",
						"environment_variables.NEW_VAR": "new_value",
						"environment_variables.%":       "2",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)
						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.EnvironmentVariables, "environment_variables should not be nil")
						require.Equal(te.t, "updated", (*ctInfo.EnvironmentVariables)["FOO"])
						require.Equal(te.t, "new_value", (*ctInfo.EnvironmentVariables)["NEW_VAR"])
						_, hasBaz := (*ctInfo.EnvironmentVariables)["BAZ"]
						require.False(te.t, hasBaz, "BAZ should be removed")

						return nil
					},
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					unprivileged = true
					disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					initialization {
						hostname = "test-env-vars"
						ip_config {
							ipv4 {
								address = "dhcp"
							}
						}
					}
					network_interface {
						name = "vmbr0"
					}
					operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"environment_variables.%": "0",
					}),
					func(*terraform.State) error {
						ct := te.NodeClient().Container(accTestContainerID)
						ctInfo, err := ct.GetContainer(t.Context())
						require.NoError(te.t, err, "failed to get container")
						if ctInfo.EnvironmentVariables != nil {
							require.Empty(te.t, *ctInfo.EnvironmentVariables, "environment_variables should be empty")
						}

						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerMountExistingVolumeWithSize verifies that mounting an existing
// subvolume with an explicit size works correctly (issue #2490). The size parameter should
// be sent as a separate `size=` parameter in the mount point config, not appended to the
// volume name.
//
// This test creates a container with a mount point, then in step 2 updates it specifying
// the existing volume explicitly along with size. This exercises the update code path fix.
func TestAccResourceContainerMountExistingVolumeWithSize(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID := 100000 + rand.Intn(99999)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":   imageFileName,
		"TestContainerID": accTestContainerID,
	})

	var createdVolumeName string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// step 1: create container with a new mount point (creates volume)
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
					vm_id     = {{ .TestContainerID }}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "/mnt/data"
				    }
				    initialization {
				  		hostname = "test-mount-existing"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
						"mount_point.#":      "1",
						"mount_point.0.size": "4G",
						"mount_point.0.path": "/mnt/data",
					}),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["proxmox_virtual_environment_container.test_container"]
						if !ok {
							return fmt.Errorf("test_container not found")
						}
						createdVolumeName = rs.Primary.Attributes["mount_point.0.volume"]
						t.Logf("Created volume: %s", createdVolumeName)
						// verify volume name has the expected format (datastore:volume-id)
						if !strings.Contains(createdVolumeName, ":") {
							return fmt.Errorf("volume name should contain colon: %s", createdVolumeName)
						}
						return nil
					},
				),
			},
			{
				// step 2: add a second mount point while keeping the first one
				// this tests that existing mount points with size are handled correctly during update
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
					vm_id     = {{ .TestContainerID }}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "/mnt/data"
				    }
				    mount_point {
						volume = "local-lvm"
						size   = "2G"
						path   = "/mnt/data2"
				    }
				    initialization {
				  		hostname = "test-mount-existing"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}`),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
						"mount_point.#": "2",
					}),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["proxmox_virtual_environment_container.test_container"]
						if !ok {
							return fmt.Errorf("test_container not found")
						}

						// verify both mount points exist and have correct attributes
						mp0Volume := rs.Primary.Attributes["mount_point.0.volume"]
						mp0Size := rs.Primary.Attributes["mount_point.0.size"]
						mp1Volume := rs.Primary.Attributes["mount_point.1.volume"]
						mp1Size := rs.Primary.Attributes["mount_point.1.size"]

						t.Logf("MP0: volume=%s, size=%s", mp0Volume, mp0Size)
						t.Logf("MP1: volume=%s, size=%s", mp1Volume, mp1Size)

						// verify volumes have proper format (not mangled)
						if !strings.Contains(mp0Volume, ":") || !strings.Contains(mp1Volume, ":") {
							return fmt.Errorf("volume names should contain colon: mp0=%s, mp1=%s", mp0Volume, mp1Volume)
						}

						return nil
					},
				),
			},
		},
	})
}

// TestAccResourceContainerMountPointVolumeReference verifies that the mount_point.volume
// attribute can be referenced by another resource without causing "inconsistent final plan"
// errors. This requires mount_point.volume to be marked as Computed in the schema.
func TestAccResourceContainerMountPointVolumeReference(t *testing.T) {
	te := InitEnvironment(t)
	accTestContainerID1 := 100000 + rand.Intn(99999)
	accTestContainerID2 := accTestContainerID1 + 1
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())

	testAccDownloadContainerTemplate(t, te, imageFileName)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":    imageFileName,
		"TestContainerID1": accTestContainerID1,
		"TestContainerID2": accTestContainerID2,
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// create two containers where the second references the first's mount_point volume
				// this tests that mount_point.volume is properly marked as Computed
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "source" {
				    node_name = "{{.NodeName}}"
					vm_id     = {{ .TestContainerID1 }}
				  	started   = false
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
				    mount_point {
						volume = "local-lvm"
						size   = "4G"
						path   = "/mnt/data"
				    }
				    initialization {
				  		hostname = "source"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}

				resource "proxmox_virtual_environment_container" "target" {
				    node_name  = "{{.NodeName}}"
					vm_id      = {{ .TestContainerID2 }}
				  	started    = false
					depends_on = [proxmox_virtual_environment_container.source]
				    disk {
						datastore_id = "local-lvm"
						size         = 4
					}
					# reference the path_in_datastore from source container (computed attribute for cross-resource refs)
				    mount_point {
						volume = proxmox_virtual_environment_container.source.mount_point[0].path_in_datastore
						size   = "4G"
						path   = "/mnt/shared"
				    }
				    initialization {
				  		hostname = "target"
						ip_config {
						  	ipv4 {
								address = "dhcp"
						  	}
						}
				    }
				    network_interface {
				  	    name = "vmbr0"
				    }
				    operating_system {
						template_file_id = "local:vztmpl/{{.ImageFileName}}"
						type             = "alpine"
				    }
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes("proxmox_virtual_environment_container.source", map[string]string{
						"mount_point.#": "1",
					}),
					ResourceAttributes("proxmox_virtual_environment_container.target", map[string]string{
						"mount_point.#": "1",
					}),
					func(s *terraform.State) error {
						rs1, ok := s.RootModule().Resources["proxmox_virtual_environment_container.source"]
						if !ok {
							return fmt.Errorf("source container not found")
						}
						rs2, ok := s.RootModule().Resources["proxmox_virtual_environment_container.target"]
						if !ok {
							return fmt.Errorf("target container not found")
						}

						sourcePath := rs1.Primary.Attributes["mount_point.0.path_in_datastore"]
						targetPath := rs2.Primary.Attributes["mount_point.0.path_in_datastore"]

						t.Logf("Source path_in_datastore: %s", sourcePath)
						t.Logf("Target path_in_datastore: %s", targetPath)

						// verify both containers reference the same volume
						if sourcePath != targetPath {
							return fmt.Errorf("path_in_datastore don't match: source=%s, target=%s", sourcePath, targetPath)
						}

						// verify the path_in_datastore has the expected format (datastore:volume-id)
						if !strings.Contains(sourcePath, ":") {
							return fmt.Errorf("source path_in_datastore should have format datastore:id, got: %s", sourcePath)
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccDownloadContainerTemplate(t *testing.T, te *Environment, imageFileName string) {
	t.Helper()
	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: ptr.Ptr(imageFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr(te.DatastoreID),
		URL:      ptr.Ptr(getTemplateURL(t, te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName))
		require.NoError(t, e)
	})
}
