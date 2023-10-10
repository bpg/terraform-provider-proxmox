/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
	"github.com/bpg/terraform-provider-proxmox/utils"
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

//nolint:gochecknoglobals
var (
	once        sync.Once
	nodesClient *nodes.Client
)

func getNodesClient() *nodes.Client {
	if nodesClient == nil {
		once.Do(
			func() {
				username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME")
				password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")
				endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
				apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")

				creds, err := api.NewCredentials(username, password, "", apiToken)
				if err != nil {
					panic(err)
				}

				conn, err := api.NewConnection(endpoint, true)
				if err != nil {
					panic(err)
				}

				client, err := api.NewClient(creds, conn)
				if err != nil {
					panic(err)
				}

				nodesClient = &nodes.Client{Client: client, NodeName: accTestNodeName}
			})
	}

	return nodesClient
}
