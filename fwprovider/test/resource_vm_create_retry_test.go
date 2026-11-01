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
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVMCreateNoWorkerUPID(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{
			name: "sdk vm",
			config: `
			resource "proxmox_virtual_environment_vm" "test" {
				node_name = "{{.NodeName}}"
				vm_id     = {{.TestVMID}}
				name      = "test-create-no-worker"
				started   = false
			}`,
		},
		{
			name: "framework vm",
			config: `
			resource "proxmox_vm" "test" {
				node_name = "{{.NodeName}}"
				id        = {{.TestVMID}}
				name      = "test-create-no-worker"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := InitEnvironment(t)
			vmID := 100000 + rand.Intn(99999)
			te.AddTemplateVars(map[string]any{"TestVMID": vmID})

			// PVE leaves the config of the failed create locked, and only root can unlock it.
			t.Cleanup(func() {
				nodeExecRootCleanup(te, fmt.Sprintf("qm unlock %d; qm destroy %d --purge", vmID, vmID))
			})

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							// SIGCHLD interrupts pvedaemon's wait for a new worker's UPID, so the next task start fails
							// with "got no worker upid". The signals start once a pvedaemon worker forks a new child,
							// since earlier ones fail the request before it reaches the create.
							te.ExecuteNodeCommands([]string{
								`(nohup timeout 60 perl -MTime::HiRes=usleep,time -e '` +
									`sub kids { open(my $f, "<", "/proc/$_[0]/task/$_[0]/children") or return; split " ", <$f> // "" } ` +
									`my %old = map { $_ => 1 } map { kids($_) } @ARGV; ` +
									`while (1) { for my $w (@ARGV) { next unless grep { !$old{$_} } kids($w); ` +
									`my $end = time + 0.02; while (time < $end) { kill "CHLD", $w; usleep 20 } exit } }' ` +
									`$(pgrep -f '^pvedaemon worker') >/dev/null 2>&1 </dev/null &)`,
								"sleep 1",
							})
						},
						Config:      te.RenderConfig(tt.config),
						ExpectError: regexp.MustCompile(`got no worker\s+upid`),
					},
				},
			})
		})
	}
}
