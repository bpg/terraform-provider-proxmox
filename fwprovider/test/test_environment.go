/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Environment is a test environment for acceptance tests.
type Environment struct {
	t              *testing.T
	templateVars   map[string]any
	providerConfig string
	NodeName       string
	DatastoreID    string

	AccProviders          map[string]func() (tfprotov6.ProviderServer, error)
	once                  sync.Once
	c                     api.Client
	CloudImagesServer     string
	ContainerImagesServer string
}

// InitEnvironment initializes a new test environment for acceptance tests.
func InitEnvironment(t *testing.T) *Environment {
	t.Helper()

	nodeName := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_NAME")
	if nodeName == "" {
		nodeName = "pve"
	}

	nodeAddress := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if nodeAddress == "" {
		endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
		u, err := url.Parse(endpoint)
		require.NoError(t, err)

		nodeAddress = u.Hostname()
	}

	nodePort := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	if nodePort == "" {
		nodePort = "22"
	}

	pc := fmt.Sprintf(`	
provider "proxmox" {
  ssh {
	node {
	  name    = "%s"
	  address = "%s"
	  port    = %s
	}
  }
  //random_vm_ids = true
}
`, nodeName, nodeAddress, nodePort)

	const datastoreID = "local"

	cloudImagesServer := utils.GetAnyStringEnv("PROXMOX_VE_ACC_CLOUD_IMAGES_SERVER")
	if cloudImagesServer == "" {
		cloudImagesServer = "https://cloud-images.ubuntu.com"
	}

	containerImagesServer := utils.GetAnyStringEnv("PROXMOX_VE_ACC_CONTAINER_IMAGES_SERVER")
	if containerImagesServer == "" {
		containerImagesServer = "http://download.proxmox.com"
	}

	return &Environment{
		t: t,
		templateVars: map[string]any{
			"ProviderConfig":        pc,
			"NodeName":              nodeName,
			"DatastoreID":           datastoreID,
			"CloudImagesServer":     cloudImagesServer,
			"ContainerImagesServer": containerImagesServer,
		},
		providerConfig:        pc,
		NodeName:              nodeName,
		DatastoreID:           datastoreID,
		AccProviders:          muxProviders(t),
		CloudImagesServer:     cloudImagesServer,
		ContainerImagesServer: containerImagesServer,
	}
}

// AddTemplateVars adds the given variables to the template variables of the current test environment.
// Please note that NodeName and ProviderConfig are reserved keys, they are set by the test environment
// and cannot be overridden.
func (e *Environment) AddTemplateVars(vars map[string]any) {
	for k, v := range vars {
		e.templateVars[k] = v
	}
}

// RenderConfig renders the given configuration with for the current test environment using template engine.
func (e *Environment) RenderConfig(cfg string) string {
	tmpl, err := template.New("config").Parse("{{.ProviderConfig}}" + cfg)
	require.NoError(e.t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, e.templateVars)
	require.NoError(e.t, err)

	return buf.String()
}

// Client returns a new API client for the test environment.
func (e *Environment) Client() api.Client {
	if e.c == nil {
		e.once.Do(
			func() {
				endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
				authTicket := utils.GetAnyStringEnv("PROXMOX_VE_AUTH_TICKET")
				csrfPreventionToken := utils.GetAnyStringEnv("PROXMOX_VE_CSRF_PREVENTION_TOKEN")
				apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")
				username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME")
				password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")

				creds, err := api.NewCredentials(username, password, "", apiToken, authTicket, csrfPreventionToken)
				if err != nil {
					panic(err)
				}

				conn, err := api.NewConnection(endpoint, true, "")
				if err != nil {
					panic(err)
				}

				e.c, err = api.NewClient(creds, conn)
				if err != nil {
					panic(err)
				}
			})
	}

	return e.c
}

// AccessClient returns a new access client for the test environment.
func (e *Environment) AccessClient() *access.Client {
	return &access.Client{Client: e.Client()}
}

// NodeClient returns a new nodes client for the test environment.
func (e *Environment) NodeClient() *nodes.Client {
	return &nodes.Client{Client: e.Client(), NodeName: e.NodeName}
}

// NodeStorageClient returns a new storage client for the test environment.
func (e *Environment) NodeStorageClient() *storage.Client {
	return &storage.Client{Client: e.NodeClient(), StorageName: e.DatastoreID}
}

// ClusterClient returns a new cluster client for the test environment.
func (e *Environment) ClusterClient() *cluster.Client {
	return &cluster.Client{Client: e.Client()}
}

// testAccMuxProviders returns a map of mux servers for the acceptance tests.
func muxProviders(t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()

	ctx := context.Background()

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
