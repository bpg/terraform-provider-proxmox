/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
)

// TestMuxServerProviderSchemaParity ensures the Framework and SDK provider schemas stay identical.
// tf6muxserver compares them when the server is created, so any attribute added to only one of the
// two providers - or described differently in each - breaks the provider at startup. Acceptance
// tests would catch it, but they only run on demand, so this guards the invariant in PR CI.
func TestMuxServerProviderSchemaParity(t *testing.T) {
	t.Parallel()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		t.Context(),
		func() tfprotov5.ProviderServer {
			return schema.NewGRPCProviderServer(provider.ProxmoxVirtualEnvironment())
		},
	)
	require.NoError(t, err)

	muxServer, err := tf6muxserver.NewMuxServer(
		t.Context(),
		providerserver.NewProtocol6(fwprovider.New("test")()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	)
	require.NoError(t, err)

	// NewMuxServer defers the comparison to the first GetProviderSchema call, which reports a
	// mismatch as an error diagnostic rather than as a Go error.
	resp, err := muxServer.ProviderServer().GetProviderSchema(t.Context(), &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)

	for _, d := range resp.Diagnostics {
		require.NotEqualf(t, tfprotov6.DiagnosticSeverityError, d.Severity, "%s: %s", d.Summary, d.Detail)
	}
}
