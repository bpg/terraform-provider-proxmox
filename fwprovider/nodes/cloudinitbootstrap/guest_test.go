/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

func TestCheckCloudInitResult(t *testing.T) {
	t.Parallel()

	require.NoError(t, checkCloudInitResult(`{"v1": {"datasource": "DataSourceNoCloud [seed=/dev/sr0]", "errors": []}}`))

	err := checkCloudInitResult(`{"v1": {"errors": ["module a failed", "module b failed"]}}`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "module a failed; module b failed")

	require.Error(t, checkCloudInitResult(`not json`))
}

func TestIsRetryableUploadError(t *testing.T) {
	t.Parallel()

	transport := fmt.Errorf("error uploading file to datastore local: failed to perform HTTP POST request - Reason: %w",
		&url.Error{Op: "Post", URL: "https://pve:8006/api2/json/nodes/n/storage/local/upload", Err: io.EOF})

	assert.True(t, isRetryableUploadError(transport))
	assert.True(t, isRetryableUploadError(fmt.Errorf("wrapped: %w", syscall.ECONNRESET)))
	assert.False(t, isRetryableUploadError(errors.New("403 Forbidden: permission check failed")))
}

func TestIsPermanentAgentError(t *testing.T) {
	t.Parallel()

	wrap := func(err error) error { return fmt.Errorf("error reading file from VM 100 via agent: %w", err) }

	assert.True(t, isPermanentAgentError(wrap(&api.HTTPError{Code: 403, Message: "Permission check failed"})))
	assert.True(t, isPermanentAgentError(wrap(&api.HTTPError{Code: 500, Message: "No QEMU guest agent configured"})))
	assert.False(t, isPermanentAgentError(wrap(&api.HTTPError{Code: 500, Message: "QEMU guest agent is not running"})))
	// A missing file is how PVE reports "not written yet"; it must keep polling.
	notFound := errors.Join(api.ErrResourceDoesNotExist, &api.HTTPError{Code: 500, Message: "No such file or directory"})
	assert.False(t, isPermanentAgentError(wrap(notFound)))
}

func TestIPv4ByInterface(t *testing.T) {
	t.Parallel()

	ip := func(addr, typ string) vms.GetQEMUNetworkInterfacesResponseResultIPAddress {
		return vms.GetQEMUNetworkInterfacesResponseResultIPAddress{Address: addr, Type: typ}
	}
	nic := func(name, mac string, addrs ...vms.GetQEMUNetworkInterfacesResponseResultIPAddress) vms.GetQEMUNetworkInterfacesResponseResult {
		return vms.GetQEMUNetworkInterfacesResponseResult{Name: name, MACAddress: mac, IPAddresses: &addrs}
	}
	macs := map[string]bool{"bc:24:11:00:00:01": true, "bc:24:11:00:00:02": true}

	out, complete := ipv4ByInterface(macs, []vms.GetQEMUNetworkInterfacesResponseResult{
		nic("lo", "00:00:00:00:00:00", ip("127.0.0.1", "ipv4")),
		nic("eth0", "BC:24:11:00:00:01", ip("10.0.0.5", "ipv4"), ip("fe80::1", "ipv6")),
		nic("eth1", "bc:24:11:00:00:02", ip("192.0.2.7", "ipv4")),
		nic("kube-ipvs0", "aa:bb:cc:dd:ee:ff", ip("10.96.0.1", "ipv4")),
	})
	assert.True(t, complete)
	assert.Len(t, out, 2, "only Proxmox network devices are kept")
	assert.Equal(t, `["10.0.0.5"]`, out["eth0"].String())

	// DHCP not done on eth1 yet: incomplete, so Read keeps the last known addresses.
	_, complete = ipv4ByInterface(macs, []vms.GetQEMUNetworkInterfacesResponseResult{
		nic("eth0", "bc:24:11:00:00:01", ip("10.0.0.5", "ipv4")),
		nic("eth1", "bc:24:11:00:00:02", ip("fe80::2", "ipv6")),
	})
	assert.False(t, complete)

	// VLAN subinterface and bridge in the guest share the NIC's MAC; one of them having an address covers it.
	out, complete = ipv4ByInterface(map[string]bool{"bc:24:11:00:00:01": true}, []vms.GetQEMUNetworkInterfacesResponseResult{
		nic("eth0", "bc:24:11:00:00:01"),
		nic("eth0.10", "bc:24:11:00:00:01", ip("10.10.0.5", "ipv4")),
		nic("br0", "bc:24:11:00:00:01", ip("10.0.0.5", "ipv4")),
	})
	assert.True(t, complete)
	assert.Len(t, out, 2, "an interface without an address is dropped when another covers its MAC")
	assert.NotContains(t, out, "eth0")

	// A VLAN covering one NIC doesn't cover another NIC that has no address.
	out, complete = ipv4ByInterface(macs, []vms.GetQEMUNetworkInterfacesResponseResult{
		nic("eth0", "bc:24:11:00:00:01"),
		nic("eth0.10", "bc:24:11:00:00:01", ip("10.10.0.5", "ipv4")),
		{Name: "eth1", MACAddress: "bc:24:11:00:00:02"}, // the agent omits ip-addresses for an unaddressed interface
	})
	assert.False(t, complete)
	assert.Equal(t, `[]`, out["eth1"].String(), "an uncovered NIC is listed without addresses")
	assert.NotContains(t, out, "eth0")

	// The agent doesn't report a device at all yet.
	_, complete = ipv4ByInterface(macs, []vms.GetQEMUNetworkInterfacesResponseResult{
		nic("eth0", "bc:24:11:00:00:01", ip("10.0.0.5", "ipv4")),
	})
	assert.False(t, complete)
}
