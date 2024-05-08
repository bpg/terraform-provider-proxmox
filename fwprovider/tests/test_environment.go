package tests

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

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

type testEnvironment struct {
	t              *testing.T
	templateVars   map[string]any
	providerConfig string
	nodeName       string
	datastoreID    string

	accProviders map[string]func() (tfprotov6.ProviderServer, error)
	once         sync.Once
	c            api.Client
}

func initTestEnvironment(t *testing.T) *testEnvironment {
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
}
`, nodeName, nodeAddress, nodePort)

	const datastoreID = "local"

	return &testEnvironment{
		t: t,
		templateVars: map[string]any{
			"ProviderConfig": pc,
			"NodeName":       nodeName,
			"DatastoreID":    datastoreID,
		},
		providerConfig: pc,
		nodeName:       nodeName,
		datastoreID:    datastoreID,
		accProviders:   muxProviders(t),
	}
}

// addTemplateVars adds the given variables to the template variables of the current test environment.
// Please note that NodeName and ProviderConfig are reserved keys, they are set by the test environment
// and cannot be overridden.
func (e *testEnvironment) addTemplateVars(vars map[string]any) {
	for k, v := range vars {
		e.templateVars[k] = v
	}
}

// renderConfig renders the given configuration with for the current test environment using template engine.
func (e *testEnvironment) renderConfig(cfg string) string {
	tmpl, err := template.New("config").Parse("{{.ProviderConfig}}" + cfg)
	require.NoError(e.t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, e.templateVars)
	require.NoError(e.t, err)

	return buf.String()
}

func (e *testEnvironment) client() api.Client {
	if e.c == nil {
		e.once.Do(
			func() {
				username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME")
				password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")
				endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
				apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")

				creds, err := api.NewCredentials(username, password, "", apiToken)
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

func (e *testEnvironment) accessClient() *access.Client {
	return &access.Client{Client: e.client()}
}

func (e *testEnvironment) nodeClient() *nodes.Client {
	return &nodes.Client{Client: e.client(), NodeName: e.nodeName}
}

func (e *testEnvironment) nodeStorageClient() *storage.Client {
	return &storage.Client{Client: e.nodeClient(), StorageName: e.datastoreID}
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
