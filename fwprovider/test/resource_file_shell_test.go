//go:build acceptance || all

//testacc:tier=heavy
//testacc:resource=file

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceFileNameNotShellInterpreted(t *testing.T) {
	te := InitEnvironment(t)

	// file_name rejects path separators and PVE cannot route a DELETE for a name with spaces, so the payload uses
	// brace expansion ("{touch,x}" expands to "touch x") and a relative path in the SSH user's home directory
	markerName := SafeResourceName("file-name-marker")
	marker := fmt.Sprintf("\"$(getent passwd %s | cut -d: -f6)\"/%s", te.SSHClient().Username(), markerName)
	probe := SafeResourceName("file-name-probe")
	fileName := probe + ";{touch," + markerName + "}"

	t.Cleanup(func() {
		nodeExecRootCleanup(te, "rm -f "+marker)

		for _, name := range []string{probe, fileName} {
			_ = te.NodeStorageClient().DeleteDatastoreFile(context.Background(), "snippets/"+name).Err()
		}
	})

	te.AddTemplateVars(map[string]any{"FileName": fileName})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "probe" {
					content_type = "snippets"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					source_raw {
						data      = "probe"
						file_name = "{{.FileName}}"
					}
				}`),
			Check: resource.ComposeTestCheckFunc(
				ResourceAttributes("proxmox_virtual_environment_file.probe", map[string]string{
					"file_name": "^" + regexp.QuoteMeta(fileName) + "$",
				}),
				assertNodePathAbsent(te, marker),
			),
		}},
	})
}

func TestAccResourceFileNameRejectsPathSeparators(t *testing.T) {
	te := InitEnvironment(t)

	steps := make([]resource.TestStep, 0, 2)

	for _, mode := range []string{"stream", "sftp"} {
		steps = append(steps, resource.TestStep{
			Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_file" "escape" {
					content_type = "snippets"
					datastore_id = "local"
					node_name    = "{{.NodeName}}"
					upload_mode  = "` + mode + `"
					source_raw {
						data      = "escape"
						file_name = "../escape.txt"
					}
				}`),
			ExpectError: regexp.MustCompile(`(?s)file_name.*plain file name without path separators`),
		})
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps:                    steps,
	})
}

func assertNodePathAbsent(te *Environment, path string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		out, err := te.ExecuteNodeCommandsContext(
			context.Background(), []string{fmt.Sprintf("test -e %s && echo PRESENT || echo ABSENT", path)},
		)
		if err != nil {
			return err
		}

		if strings.TrimSpace(out) != "ABSENT" {
			return fmt.Errorf("%s exists on the node: an injected command ran", path)
		}

		return nil
	}
}

func nodeExecRootCleanup(te *Environment, command string) {
	if _, err := te.ExecuteNodeCommandsContext(context.Background(), []string{command}); err != nil {
		te.t.Logf("cleanup: %q failed: %v", command, err)
	}
}
