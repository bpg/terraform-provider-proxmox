//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
)

// TestAccResourceDownloadFileUpstreamChange tests that when the upstream URL
// reports a different Content-Length than the local file, replacement is triggered.
//
// This test uses a controllable HTTP server running on the test machine.
//
// REQUIREMENTS:
// Set PROXMOX_VE_ACC_TEST_FILE_SERVER_IP to the IP that Proxmox can use to reach
// the test machine. For example:
//
//	export PROXMOX_VE_ACC_TEST_FILE_SERVER_IP=192.168.1.100
//
// The test machine must:
//   - Have the specified IP reachable from Proxmox
//   - Allow incoming connections on a random high port
//
// This test is skipped if PROXMOX_VE_ACC_TEST_FILE_SERVER_IP is not set.
func TestAccResourceDownloadFileUpstreamChange(t *testing.T) {
	fileServer := test.NewTestFileServer(t)
	if fileServer == nil {
		t.Skip("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP not set - skipping upstream change test")
	}

	te := test.InitEnvironment(t)

	// set initial content (15 bytes)
	fileServer.SetContent([]byte("initial content"))
	fileServer.SetFilename("upstream_test.iso")

	te.AddTemplateVars(map[string]interface{}{
		"FileURL": fileServer.FileURL(),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Download file - local size matches URL size (15 bytes)
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "upstream_test" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "upstream_change_test.iso"
						overwrite           = true
						overwrite_unmanaged = true
						verify              = false
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_download_file.upstream_test", map[string]string{
						"size":      "15",
						"overwrite": "true",
					}),
				),
			},
			{
				// Step 2: Change server's Content-Length to simulate upstream update
				// Local file is still 15 bytes, but URL now reports 100 bytes
				PreConfig: func() {
					// simulate upstream releasing a new version with different size
					fileServer.SetReportedSize(100)
					t.Log("Changed server to report Content-Length: 100 (was 15)")
				},
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "upstream_test" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "upstream_change_test.iso"
						overwrite           = true
						overwrite_unmanaged = true
						verify              = false
					}`),
				// with overwrite=true, size mismatch should trigger replacement
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"proxmox_virtual_environment_download_file.upstream_test",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
				// after apply, size should be 100 (the new reported size from URL)
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_download_file.upstream_test", map[string]string{
						"size": "100",
					}),
				),
			},
		},
	})
}

// TestAccResourceDownloadFileUpstreamChangeIgnored tests that with overwrite=false,
// upstream size changes are NOT detected (no replacement triggered).
func TestAccResourceDownloadFileUpstreamChangeIgnored(t *testing.T) {
	fileServer := test.NewTestFileServer(t)
	if fileServer == nil {
		t.Skip("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP not set - skipping test")
	}

	te := test.InitEnvironment(t)

	fileServer.SetContent([]byte("initial content"))
	fileServer.SetFilename("no_upstream_check.iso")

	te.AddTemplateVars(map[string]interface{}{
		"FileURL": fileServer.FileURL(),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Download file with overwrite=false
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "no_check" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "no_upstream_check.iso"
						overwrite           = false
						overwrite_unmanaged = true
						verify              = false
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_download_file.no_check", map[string]string{
						"size":      "15",
						"overwrite": "false",
					}),
				),
			},
			{
				// Step 2: Change server's Content-Length
				// With overwrite=false, this should be IGNORED
				PreConfig: func() {
					fileServer.SetReportedSize(200)
					t.Log("Changed server to report Content-Length: 200 - but overwrite=false, so should be ignored")
				},
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "no_check" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "no_upstream_check.iso"
						overwrite           = false
						overwrite_unmanaged = true
						verify              = false
					}`),
				// with overwrite=false, URL is not checked - expect empty plan
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccResourceDownloadFileUpstreamNoChange tests that when URL size matches
// local file size, no replacement is triggered (steady state).
func TestAccResourceDownloadFileUpstreamNoChange(t *testing.T) {
	fileServer := test.NewTestFileServer(t)
	if fileServer == nil {
		t.Skip("PROXMOX_VE_ACC_TEST_FILE_SERVER_IP not set - skipping test")
	}

	te := test.InitEnvironment(t)

	fileServer.SetContent([]byte("stable content"))
	fileServer.SetFilename("stable.iso")

	te.AddTemplateVars(map[string]interface{}{
		"FileURL": fileServer.FileURL(),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Download file
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "stable" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "stable_test.iso"
						overwrite           = true
						overwrite_unmanaged = true
						verify              = false
					}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_download_file.stable", map[string]string{
						"size": "14",
					}),
				),
			},
			{
				// Step 2: Same config, no changes - should be empty plan
				// URL is checked (overwrite=true) but sizes match
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_download_file" "stable" {
						content_type        = "iso"
						node_name           = "{{.NodeName}}"
						datastore_id        = "{{.DatastoreID}}"
						url                 = "{{.FileURL}}"
						file_name           = "stable_test.iso"
						overwrite           = true
						overwrite_unmanaged = true
						verify              = false
					}`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
