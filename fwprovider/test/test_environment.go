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
	"maps"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

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
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	sdkV2provider "github.com/bpg/terraform-provider-proxmox/proxmoxtf/provider"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Environment is a test environment for acceptance tests.
type Environment struct {
	t              *testing.T
	templateVars   map[string]any
	NodeName       string
	Node2Name      string
	DatastoreID    string
	ZfsDatastoreID string
	ZfsDisk        string

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

	node2Name := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_2_NAME")

	const datastoreID = "local"

	zfsDatastoreID := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ZFS_DATASTORE_ID")
	zfsDisk := utils.GetAnyStringEnv("PROXMOX_VE_ACC_ZFS_DISK")

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
			"Node2Name":             node2Name,
			"DatastoreID":           datastoreID,
			"CloudImagesServer":     cloudImagesServer,
			"ContainerImagesServer": containerImagesServer,
			"TestName":              sanitizeTemplateName(t.Name()),
			"ZfsDatastoreID":        zfsDatastoreID,
			"ZfsDisk":               zfsDisk,
		},
		NodeName:              nodeName,
		Node2Name:             node2Name,
		DatastoreID:           datastoreID,
		ZfsDatastoreID:        zfsDatastoreID,
		ZfsDisk:               zfsDisk,
		CloudImagesServer:     cloudImagesServer,
		ContainerImagesServer: containerImagesServer,

		AccProviders: muxProviders(t),
	}
}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func sanitizeTemplateName(name string) string {
	sanitized := strings.Trim(nonAlnum.ReplaceAllString(name, "-"), "-")
	if sanitized == "" {
		return "test"
	}

	if len(sanitized) > 48 {
		return sanitized[:48]
	}

	return sanitized
}

// AddTemplateVars adds the given variables to the template variables of the current test environment.
// Please note that NodeName and ProviderConfig are reserved keys, they are set by the test environment
// and cannot be overridden.
func (e *Environment) AddTemplateVars(vars map[string]any) {
	maps.Copy(e.templateVars, vars)
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

// SSHClient returns an SSH client configured identically to the provider's SSH client:
// PROXMOX_VE_SSH_USERNAME / PROXMOX_VE_SSH_PASSWORD are used first, falling back to the
// stripped PROXMOX_VE_USERNAME / PROXMOX_VE_PASSWORD when the SSH-specific vars are unset.
// All other SSH settings (agent, private key, socks5) mirror the provider's env var lookup.
// The node resolver maps every node name to the configured test node SSH address and port.
func (e *Environment) SSHClient() ssh.Client {
	e.t.Helper()

	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME")
	if sshUsername == "" {
		sshUsername = strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	}

	sshPassword := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PASSWORD")
	if sshPassword == "" {
		sshPassword = utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")
	}

	address := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if address == "" {
		u, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
		require.NoError(e.t, err)

		address = u.Hostname()
	}

	port := int32(22)

	if p := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT"); p != "" {
		v, err := strconv.ParseInt(p, 10, 32)
		require.NoError(e.t, err)

		port = int32(v)
	}

	client, err := ssh.NewClient(
		sshUsername,
		sshPassword,
		utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT"),
		utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK"),
		utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT_FORWARDING"),
		utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY"),
		utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_SERVER"),
		utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_USERNAME"),
		utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_PASSWORD"),
		staticNodeResolver{node: ssh.ProxmoxNode{Address: address, Port: port}},
	)
	require.NoError(e.t, err)

	return client
}

// staticNodeResolver resolves any node name to a fixed SSH address/port.
type staticNodeResolver struct {
	node ssh.ProxmoxNode
}

func (r staticNodeResolver) Resolve(context.Context, string) (ssh.ProxmoxNode, error) {
	return r.node, nil
}

// ExecuteNodeCommands runs shell commands on the test node over SSH as the root API user
// (PROXMOX_VE_USERNAME / PROXMOX_VE_PASSWORD) and returns the combined output. Connecting as
// root avoids the restricted sudoers allowlist of the provider's SSH user, so privileged
// commands such as `pct exec` work without sudo.
func (e *Environment) ExecuteNodeCommands(commands []string) string {
	e.t.Helper()

	address := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_ADDRESS")
	if address == "" {
		u, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
		require.NoError(e.t, err)

		address = u.Hostname()
	}

	port := int32(22)

	if p := utils.GetAnyStringEnv("PROXMOX_VE_ACC_NODE_SSH_PORT"); p != "" {
		v, err := strconv.ParseInt(p, 10, 32)
		require.NoError(e.t, err)

		port = int32(v)
	}

	// Strip the realm from "root@pam" to get the SSH login name.
	username := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]

	client, err := ssh.NewClient(
		username,
		utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD"),
		false, "", false,
		"",
		"", "", "",
		staticNodeResolver{node: ssh.ProxmoxNode{Address: address, Port: port}},
	)
	require.NoError(e.t, err)

	out, err := client.ExecuteNodeCommands(e.t.Context(), e.NodeName, commands)
	require.NoError(e.t, err)

	return string(out)
}

// DownloadCloudImage downloads a cloud image with a unique filename for use in VM tests.
// The image is automatically cleaned up after the test completes.
// Returns the file ID in format "local:iso/filename".
func (e *Environment) DownloadCloudImage() string {
	e.t.Helper()

	fileName := "ubuntu-24.04-minimal-cloudimg-amd64.img"
	imageFileName := fmt.Sprintf("%d-%s", time.Now().UnixMicro(), fileName)
	err := e.NodeStorageClient().DownloadFileByURL(context.Background(), &storage.DownloadURLPostRequestBody{
		Content:  new("iso"),
		FileName: new(imageFileName),
		Node:     new(e.NodeName),
		Storage:  new("local"),
		URL:      new(fmt.Sprintf("%s/minimal/releases/noble/release/%s", e.CloudImagesServer, fileName)),
	})
	require.NoError(e.t, err)

	e.t.Cleanup(func() {
		// Best effort cleanup - the file may already be deleted by Proxmox
		err = e.NodeStorageClient().DeleteDatastoreFile(context.Background(), fmt.Sprintf("iso/%s", imageFileName))
		if err != nil {
			e.t.Logf("cleanup: failed to delete cloud image %s: %v", imageFileName, err)
		}
	})

	return fmt.Sprintf("local:iso/%s", imageFileName)
}

// ClusterClient returns a new cluster client for the test environment.
func (e *Environment) ClusterClient() *cluster.Client {
	return &cluster.Client{Client: e.Client()}
}

// PoolsClient returns a new pools client for the test environment.
func (e *Environment) PoolsClient() *pools.Client {
	return &pools.Client{Client: e.Client()}
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
