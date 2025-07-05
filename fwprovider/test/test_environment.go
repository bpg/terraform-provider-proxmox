/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"bytes"
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
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Environment is a test environment for acceptance tests.
type Environment struct {
	t            *testing.T
	templateVars map[string]any
	NodeName     string
	DatastoreID  string

	AccProviders          map[string]func() (tfprotov6.ProviderServer, error)
	once                  sync.Once
	c                     api.Client
	CloudImagesServer     string
	ContainerImagesServer string
}

// RenderConfigOption is a configuration option for rendering the provider configuration.
type RenderConfigOption interface {
	apply(rc *renderConfig) error
}

type renderConfig struct {
	providerConfig string
}

// returns the ssh configuration section of the provider config.
func (r *renderConfig) ssh() string {
	nodeName := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_NAME")
	if nodeName == "" {
		nodeName = "pve"
	}

	nodeAddress := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if nodeAddress == "" {
		endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")

		u, err := url.Parse(endpoint)
		if err != nil {
			panic(err)
		}

		nodeAddress = u.Hostname()
	}

	nodePort := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT")
	if nodePort == "" {
		nodePort = "22"
	}

	// one indent level
	return fmt.Sprintf(`
	ssh {
		node {
			name    = "%s"
			address = "%s"
			port    = %s
		}
  	}`, nodeName, nodeAddress, nodePort)
}

// WithRootUser returns a configuration option that sets the root user in the provider configuration.
func WithRootUser() RenderConfigOption {
	return &rootUserConfigOption{}
}

type rootUserConfigOption struct{}

func (o *rootUserConfigOption) apply(rc *renderConfig) error {
	if utils.GetAnyStringEnv("PROXMOX_VE_USERNAME") == "" || utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD") == "" {
		return fmt.Errorf("PROXMOX_VE_USERNAME and PROXMOX_VE_PASSWORD must be set")
	}

	rootUser := fmt.Sprintf("\tusername = \"%s\"\n\tpassword = \"%s\"\n\tapi_token = \"\"",
		utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"),
		utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD"),
	)

	rc.providerConfig = fmt.Sprintf("provider \"proxmox\" {\n%s\n%s\n}", rootUser, rc.ssh())

	return nil
}

// WithAPIToken returns a configuration option that sets the API token in the provider configuration.
func WithAPIToken() RenderConfigOption {
	return &apiTokenConfigOption{}
}

type apiTokenConfigOption struct{}

func (o *apiTokenConfigOption) apply(rc *renderConfig) error {
	if utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN") == "" {
		return fmt.Errorf("PROXMOX_VE_API_TOKEN must be set")
	}

	apiToken := fmt.Sprintf("\tapi_token = \"%s\"\n\tusername = \"\"\n\tpassword = \"\"",
		utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN"))

	rc.providerConfig = fmt.Sprintf("provider \"proxmox\" {\n%s\n%s\n}", apiToken, rc.ssh())

	return nil
}

// InitEnvironment initializes a new test environment for acceptance tests.
func InitEnvironment(t *testing.T) *Environment {
	t.Helper()

	nodeName := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_NAME")
	if nodeName == "" {
		nodeName = "pve"
	}

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
			"NodeName":              nodeName,
			"DatastoreID":           datastoreID,
			"CloudImagesServer":     cloudImagesServer,
			"ContainerImagesServer": containerImagesServer,
		},
		NodeName:              nodeName,
		DatastoreID:           datastoreID,
		CloudImagesServer:     cloudImagesServer,
		ContainerImagesServer: containerImagesServer,

		AccProviders: muxProviders(t),
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
func (e *Environment) RenderConfig(cfg string, opt ...RenderConfigOption) string {
	if len(opt) == 0 {
		opt = append(opt, WithAPIToken())
	}

	rc := &renderConfig{}
	for _, o := range opt {
		err := o.apply(rc)
		require.NoError(e.t, err, "configuration error")
	}

	tmpl, err := template.New("config").Parse(cfg)
	require.NoError(e.t, err)

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, e.templateVars)
	require.NoError(e.t, err)

	return rc.providerConfig + "\n" + buf.String()
}

// Client returns a new API client for the test environment.
// The client will be using the credentials from the environment variables, in precedence order:
// 1. API token
// 2. Ticket
// 3. User credentials.
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

// muxProviders returns a map of mux servers for the acceptance tests.
func muxProviders(t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()

	// Init mux servers
	return map[string]func() (tfprotov6.ProviderServer, error){
		"proxmox": func() (tfprotov6.ProviderServer, error) {
			return tf6muxserver.NewMuxServer(t.Context(),
				providerserver.NewProtocol6(fwprovider.New("test")()),
				func() tfprotov6.ProviderServer {
					sdkV2Provider, err := tf5to6server.UpgradeServer(
						t.Context(),
						func() tfprotov5.ProviderServer {
							return schema.NewGRPCProviderServer(
								sdkV2provider.ProxmoxVirtualEnvironment(),
							)
						},
					)
					require.NoError(t, err)

					return sdkV2Provider
				},
			)
		},
	}
}
