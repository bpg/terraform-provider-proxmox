//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=container

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceContainerCreateNoWorkerUPID(t *testing.T) {
	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	ctID := 100000 + rand.Intn(99999)
	te.AddTemplateVars(map[string]any{"ImageFileName": imageFileName, "TestContainerID": ctID})

	cleanup := fmt.Sprintf("pct unlock %d && pct destroy %d", ctID, ctID)

	t.Cleanup(func() {
		nodeExecRootCleanup(te, fmt.Sprintf("if [ -e /etc/pve/lxc/%d.conf ]; then %s; fi", ctID, cleanup))
	})

	config := te.RenderConfig(`
		resource "proxmox_virtual_environment_container" "test" {
			node_name    = "{{.NodeName}}"
			vm_id        = {{.TestContainerID}}
			started      = false
			unprivileged = true
			disk {
				datastore_id = "{{.ContainerDatastoreID}}"
				size         = 4
			}
			initialization {
				hostname = "test-create-no-worker"
			}
			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}
		}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { forceNoWorkerUPID(te) },
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)got\s+no\s+worker\s+upid.*pct\s+unlock\s+` + fmt.Sprint(ctID)),
			},
			{
				PreConfig: func() { te.ExecuteNodeCommands([]string{stopNoWorkerUPID, cleanup}) },
				Config:    config,
			},
		},
	})
}

func TestAccResourceContainerCloneNoWorkerUPID(t *testing.T) {
	te := InitEnvironment(t)
	imageFileName := fmt.Sprintf("%d-alpine-3.22-default_20250617_amd64.tar.xz", time.Now().UnixMicro())
	testAccDownloadContainerTemplate(t, te, imageFileName)

	sourceID := 100000 + rand.Intn(99999)
	cloneID := sourceID + 1
	te.AddTemplateVars(map[string]any{"ImageFileName": imageFileName, "SourceID": sourceID, "CloneID": cloneID})

	cleanup := fmt.Sprintf("pct unlock %d && pct unlock %d && pct destroy %d", sourceID, cloneID, cloneID)

	t.Cleanup(func() {
		nodeExecRootCleanup(te, fmt.Sprintf("if [ -e /etc/pve/lxc/%d.conf ]; then %s; fi", cloneID, cleanup))
	})

	source := `
		resource "proxmox_virtual_environment_container" "source" {
			node_name    = "{{.NodeName}}"
			vm_id        = {{.SourceID}}
			started      = false
			unprivileged = true
			disk {
				datastore_id = "{{.ContainerDatastoreID}}"
				size         = 4
			}
			initialization {
				hostname = "test-clone-no-worker-source"
			}
			operating_system {
				template_file_id = "local:vztmpl/{{.ImageFileName}}"
				type             = "alpine"
			}
		}`
	withClone := source + `
		resource "proxmox_virtual_environment_container" "clone" {
			node_name = "{{.NodeName}}"
			vm_id     = {{.CloneID}}
			started   = false
			clone {
				vm_id = proxmox_virtual_environment_container.source.vm_id
			}
		}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(source),
			},
			{
				PreConfig: func() { forceNoWorkerUPID(te) },
				Config:    te.RenderConfig(withClone),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(`(?s)got\s+no\s+worker\s+upid.*pct\s+unlock\s+%d\s+&&\s+pct\s+unlock\s+%d`, sourceID, cloneID),
				),
			},
			{
				PreConfig: func() { te.ExecuteNodeCommands([]string{stopNoWorkerUPID, cleanup}) },
				Config:    te.RenderConfig(withClone),
			},
		},
	})
}
