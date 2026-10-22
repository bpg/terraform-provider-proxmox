//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=vm

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVMDiskFileIDNotShellInterpreted(t *testing.T) {
	te := InitEnvironment(t)
	imageFileID := te.DownloadCloudImage()

	marker := "/tmp/" + SafeResourceName("file-id-marker")

	t.Cleanup(func() {
		nodeExecRootCleanup(te, "rm -f "+marker)
	})

	te.AddTemplateVars(map[string]any{"FileID": imageFileID + "; touch " + marker})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_vm" "probe" {
					node_name = "{{.NodeName}}"
					started   = false
					name      = "test-disk-file-id-probe"
					disk {
						datastore_id = "local-lvm"
						file_id      = "{{.FileID}}"
						interface    = "virtio0"
						size         = 20
					}
				}`),
			ExpectError: regexp.MustCompile(`creating custom disk`),
		}},
		CheckDestroy: assertNodePathAbsent(te, marker),
	})
}
