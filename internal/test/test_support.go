/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

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

	fwprovider "github.com/bpg/terraform-provider-proxmox/internal/provider"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
)

const (
	// ProviderConfig is a shared configuration to combine with the actual
	// test configuration so the Proxmox VE client is properly configured.
	// It is also possible to use the PROXMOX_VE_ environment variables instead.
	ProviderConfig = `
provider "proxmox" {
  username = "root@pam"
  password = "password"
  insecure = true
  ssh {
    agent = true
  }
}
`
)

// AccMuxProviders returns a map of mux servers for the acceptance tests.
func AccMuxProviders(ctx context.Context, t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
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
			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			return muxServer, fmt.Errorf("failed to create mux server: %w", err)
		},
	}

	return muxServers
}
