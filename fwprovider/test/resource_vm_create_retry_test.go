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

			cleanup := fmt.Sprintf("qm unlock %d && qm destroy %d", vmID, vmID)

			t.Cleanup(func() {
				nodeExecRootCleanup(te, fmt.Sprintf("if [ -e /etc/pve/qemu-server/%d.conf ]; then %s; fi", vmID, cleanup))
			})

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps: []resource.TestStep{
					{
						PreConfig:   func() { forceNoWorkerUPID(te) },
						Config:      te.RenderConfig(tt.config),
						ExpectError: regexp.MustCompile(`(?s)got\s+no\s+worker\s+upid.*qm\s+unlock\s+` + fmt.Sprint(vmID)),
					},
					{
						PreConfig: func() { te.ExecuteNodeCommands([]string{stopNoWorkerUPID, cleanup}) },
						Config:    te.RenderConfig(tt.config),
					},
				},
			})
		})
	}
}

// stopNoWorkerUPID stops the signals forceNoWorkerUPID started; [s] keeps pkill from matching the shell that runs it.
const stopNoWorkerUPID = "pkill -f 'perl -MTime::HiRe[s]=usleep' || true"

// forceNoWorkerUPID makes the node's next task start fail with "got no worker upid": SIGCHLD interrupts pvedaemon's
// wait for a new worker's UPID. A create can run helper commands before it forks its worker, so every new child of a
// pvedaemon worker in the 2s after the first one gets the signals. Signals sent before any fork fail the request
// before it reaches the task start.
func forceNoWorkerUPID(te *Environment) {
	te.t.Helper()

	te.t.Cleanup(func() { nodeExecRootCleanup(te, stopNoWorkerUPID) })

	te.ExecuteNodeCommands([]string{
		`(nohup timeout 60 perl -MTime::HiRes=usleep,time -e '` +
			`sub kids { open(my $f, "<", "/proc/$_[0]/task/$_[0]/children") or return; split " ", <$f> // "" } ` +
			`my %old = map { $_ => 1 } map { kids($_) } @ARGV; my $stop; ` +
			`while (!$stop || time < $stop) { for my $w (@ARGV) { my @new = grep { !$old{$_} } kids($w) or next; ` +
			`$old{$_} = 1 for @new; $stop //= time + 2; ` +
			`my $end = time + 0.02; while (time < $end) { kill "CHLD", $w; usleep 20 } } }' ` +
			`$(pgrep -f '^pvedaemon worker') >/dev/null 2>&1 </dev/null &)`,
		"sleep 1",
	})
}
