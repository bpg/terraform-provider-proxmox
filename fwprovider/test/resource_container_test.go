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

	imageFileName := gofakeit.Word() + "-ubuntu-23.04-standard_23.04-1_amd64.tar.zst"
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
		URL:      ptr.Ptr(fmt.Sprintf("%s/images/system/ubuntu-23.04-standard_23.04-1_amd64.tar.zst", te.ContainerImagesServer)),
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
						mount_options = ["discard"]
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
					ResourceAttributesSet(accTestContainerName, []string{
						"ipv4.vmbr0",
					}),
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
