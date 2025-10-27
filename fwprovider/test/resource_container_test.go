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

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

const (
	accTestContainerName = "proxmox_virtual_environment_container.test_container"
)

func TestAccResourceContainer(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	imageFileName := gofakeit.Word() + "-ubuntu-24.04-standard_24.04-2_amd64.tar.zst"
	accTestContainerID := 100000 + rand.Intn(99999)
	accTestContainerIDClone := 100000 + rand.Intn(99999)

	te.AddTemplateVars(map[string]interface{}{
		"ImageFileName":        imageFileName,
		"TestContainerID":      accTestContainerID,
		"TestContainerIDClone": accTestContainerIDClone,
	})

	err := te.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  ptr.Ptr("vztmpl"),
		FileName: ptr.Ptr(imageFileName),
		Node:     ptr.Ptr(te.NodeName),
		Storage:  ptr.Ptr(te.DatastoreID),
		URL:      ptr.Ptr(fmt.Sprintf("%s/images/system/ubuntu-24.04-standard_24.04-2_amd64.tar.zst", te.ContainerImagesServer)),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		e := te.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("vztmpl/%s", imageFileName))
		require.NoError(t, e)
	})

	tests := []struct {
		name string
		step []resource.TestStep
	}{
		{"create, start and update container", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = 10
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
						type             = "ubuntu"
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
					// TODO: depends on DHCP, which may not work in some environments
					// ResourceAttributesSet(accTestContainerName, []string{
					// 	"ipv4.vmbr0",
					// }),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						ct := te.NodeClient().Container(accTestContainerID)
						err := ct.WaitForContainerStatus(ctx, "running")
						require.NoError(te.t, err, "container did not start")

						ctInfo, err := ct.GetContainer(ctx)
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
					timeout_delete = 10
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
						type             = "ubuntu"
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
					timeout_delete = 10
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"disk.0.mount_options.#": "0",
					}),
				),
			},
		}},
		{"update mount points", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
				    node_name = "{{.NodeName}}"
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
						type             = "ubuntu"
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
						type             = "ubuntu"
				    }
				}`),
				Check: ResourceAttributes("proxmox_virtual_environment_container.test_container", map[string]string{
					"mount_point.#": "2",
				}),
			},
		}},
		{"ipv4 and ipv6", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
				started   = false
				disk {
					datastore_id = "local-lvm"
					size         = 4
				}
				initialization {
					hostname = "test"
					ip_config {
						ipv4 {
							address = "10.0.0.100/24"
							gateway = "10.0.0.1"
						}
					}
					ip_config {
						ipv6 {
							address = "2001:db8::100/64"	
							gateway = "2001:db8::1"
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
					type             = "ubuntu"
				}
			}`),
		}}},
		{"clone container", []resource.TestStep{{
			Config: te.RenderConfig(`
			resource "proxmox_virtual_environment_container" "test_container" {
				node_name = "{{.NodeName}}"
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
					type             = "ubuntu"
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
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					ct := te.NodeClient().Container(accTestContainerIDClone)
					err := ct.WaitForContainerStatus(ctx, "running")
					require.NoError(te.t, err, "container did not start")

					ctInfo, err := ct.GetContainer(ctx)
					require.NoError(te.t, err, "failed to get container")
					dev0, ok := ctInfo.PassthroughDevices["dev0"]
					require.True(te.t, ok, `"dev0" passthrough device not found`)
					require.NotNil(te.t, dev0, `"dev0" passthrough device is <nil>`)
					require.Equal(te.t, "/dev/zero", dev0.Path)

					return nil
				},
			),
		}}},
		{"hostname update", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = 10
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-hostname-1",
					}),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						ct := te.NodeClient().Container(accTestContainerID)
						err := ct.WaitForContainerStatus(ctx, "running")
						require.NoError(te.t, err, "container did not start")

						return nil
					},
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
					vm_id     = {{.TestContainerID}}
					timeout_delete = 10
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-hostname-2",
					}),
					func(*terraform.State) error {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						ct := te.NodeClient().Container(accTestContainerID)
						err := ct.WaitForContainerStatus(ctx, "running")
						require.NoError(te.t, err, "container did not start after hostname change")

						// Verify the hostname was actually updated in the container config
						ctInfo, err := ct.GetContainer(ctx)
						require.NoError(te.t, err, "failed to get container")
						require.NotNil(te.t, ctInfo.Hostname, "hostname should not be nil")
						require.Equal(te.t, "test-hostname-2", *ctInfo.Hostname, "hostname should be updated")

						return nil
					},
				),
			},
		}},
		{"empty dns block on update", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
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
						type             = "ubuntu"
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-issue",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
		}},
		{"empty dns block on create", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-create",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
		}},
		{"dns block with null values on create", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-create",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
		}},
		{"dns block with null values on update", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_container" "test_container" {
					node_name = "{{.NodeName}}"
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
					type             = "ubuntu"
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
						type             = "ubuntu"
					}
				}`, WithRootUser()),
				Check: resource.ComposeTestCheckFunc(
					ResourceAttributes(accTestContainerName, map[string]string{
						"initialization.0.hostname": "test-dns-update",
						"initialization.0.dns.#":    "0",
					}),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.step,
			})
		})
	}
}
