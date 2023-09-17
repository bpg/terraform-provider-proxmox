/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
)

const (
	accTestNodeName = "pve"
)

// testAccMuxProviders returns a map of mux servers for the acceptance tests.
func testAccMuxProviders(ctx context.Context, t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()

	// Init sdkV2 provider
	sdkV2Provider, err := tf5to6server.UpgradeServer(
		ctx,
		func() tfprotov5.ProviderServer {
			return schema.NewGRPCProviderServer(
				sdkV2provider.ProxmoxVirtualEnvironment(),
			)
		},
	)
	require.NoError(t, err)

	// Init framework provider
	frameworkProvider := fwprovider.New("test")()

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(frameworkProvider),
		func() tfprotov6.ProviderServer {
			return sdkV2Provider
		},
	}

	// Init mux servers
	muxServers := map[string]func() (tfprotov6.ProviderServer, error){
		"proxmox": func() (tfprotov6.ProviderServer, error) {
			muxServer, e := tf6muxserver.NewMuxServer(ctx, providers...)
			if e != nil {
				return nil, fmt.Errorf("failed to create mux server: %w", e)
			}
			return muxServer, nil
		},
	}

	return muxServers
}
