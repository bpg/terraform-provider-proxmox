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
	"testing"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
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
