/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/access"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/acme"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/ha"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/hardwaremapping"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/network"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/apt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &proxmoxProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &proxmoxProvider{
			version: version,
		}
	}
}

type proxmoxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// proxmoxProviderModel maps provider schema data.
type proxmoxProviderModel struct {
	APIToken types.String `tfsdk:"api_token"`
	Endpoint types.String `tfsdk:"endpoint"`
	Insecure types.Bool   `tfsdk:"insecure"`
	MinTLS   types.String `tfsdk:"min_tls"`
	OTP      types.String `tfsdk:"otp"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	SSH      []struct {
		Agent          types.Bool   `tfsdk:"agent"`
		AgentSocket    types.String `tfsdk:"agent_socket"`
		PrivateKey     types.String `tfsdk:"private_key"`
		Password       types.String `tfsdk:"password"`
		Username       types.String `tfsdk:"username"`
		Socks5Server   types.String `tfsdk:"socks5_server"`
		Socks5Username types.String `tfsdk:"socks5_username"`
		Socks5Password types.String `tfsdk:"socks5_password"`

		Nodes []struct {
			Name    types.String `tfsdk:"name"`
			Address types.String `tfsdk:"address"`
			Port    types.Int64  `tfsdk:"port"`
		} `tfsdk:"node"`
	} `tfsdk:"ssh"`
	TmpDir types.String `tfsdk:"tmp_dir"`
}

func (p *proxmoxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	// resp.TypeName = "proxmox"
	resp.TypeName = "proxmox_virtual_environment"
	resp.Version = p.version
}

func (p *proxmoxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		// Attributes specified in alphabetical order.
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: "The API token for the Proxmox VE API.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\S+@\S+!\S+=([a-zA-Z0-9-]+)$`),
						`must be a valid API token, e.g. 'USER@REALM!TOKENID=UUID'`,
					),
				},
			},
			"endpoint": schema.StringAttribute{
				Description: "The endpoint for the Proxmox VE API.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to skip the TLS verification step.",
				Optional:    true,
			},
			"min_tls": schema.StringAttribute{
				Description: "The minimum required TLS version for API calls." +
					"Supported values: `1.0|1.1|1.2|1.3`. Defaults to `1.3`.",
				Optional: true,
			},
			"otp": schema.StringAttribute{
				Description: "The one-time password for the Proxmox VE API.",
				Optional:    true,
				DeprecationMessage: "The `otp` attribute is deprecated and will be removed in a future release. " +
					"Please use the `api_token` attribute instead.",
			},
			"password": schema.StringAttribute{
				Description: "The password for the Proxmox VE API.",
				Optional:    true,
				Sensitive:   true,
			},
			"username": schema.StringAttribute{
				Description: "The username for the Proxmox VE API.",
				Optional:    true,
			},
			"tmp_dir": schema.StringAttribute{
				Description: "The alternative temporary directory.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			// have to define it as a list due to backwards compatibility
			"ssh": schema.ListNestedBlock{
				Description: "The SSH configuration for the Proxmox nodes.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"agent": schema.BoolAttribute{
							Description: "Whether to use the SSH agent for authentication. Takes precedence over " +
								"the `private_key` and `password` fields. Defaults to the value of the " +
								"`PROXMOX_VE_SSH_AGENT` environment variable, or `false` if not set.",
							Optional: true,
						},
						"agent_socket": schema.StringAttribute{
							Description: "The path to the SSH agent socket. " +
								"Defaults to the value of the `SSH_AUTH_SOCK` " +
								"environment variable.",
							Optional: true,
						},
						"private_key": schema.StringAttribute{
							Description: "The unencrypted private key (in PEM format) used for the SSH connection. " +
								"Defaults to the value of the `PROXMOX_VE_SSH_PRIVATE_KEY` environment variable.",
							Optional:  true,
							Sensitive: true,
						},
						"password": schema.StringAttribute{
							Description: "The password used for the SSH connection. " +
								"Defaults to the value of the `password` field of the " +
								"`provider` block.",
							Optional:  true,
							Sensitive: true,
						},
						"username": schema.StringAttribute{
							Description: "The username used for the SSH connection. " +
								"Defaults to the value of the `username` field of the " +
								"`provider` block.",
							Optional: true,
						},
						"socks5_server": schema.StringAttribute{
							Description: "The address:port of the SOCKS5 proxy server. " +
								"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_SERVER` environment variable.",
							Optional: true,
						},
						"socks5_username": schema.StringAttribute{
							Description: "The username for the SOCKS5 proxy server. " +
								"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_USERNAME` environment variable.",
							Optional: true,
						},
						"socks5_password": schema.StringAttribute{
							Description: "The password for the SOCKS5 proxy server. " +
								"Defaults to the value of the `PROXMOX_VE_SSH_SOCKS5_PASSWORD` environment variable.",
							Optional:  true,
							Sensitive: true,
						},
					},
					Blocks: map[string]schema.Block{
						"node": schema.ListNestedBlock{
							Description: "Overrides for SSH connection configuration for a Proxmox VE node.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the Proxmox VE node.",
										Required:    true,
									},
									"address": schema.StringAttribute{
										Description: "The address of the Proxmox VE node.",
										Required:    true,
									},
									"port": schema.Int64Attribute{
										Description: "The port of the Proxmox VE node.",
										Optional:    true,
										Validators:  []validator.Int64{int64validator.Between(1, 65535)},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (p *proxmoxProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	tflog.Info(ctx, "Configuring the Proxmox provider...")

	// Retrieve provider data from configuration
	var config proxmoxProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Proxmox VE API Endpoint",
			"The provider cannot create the Proxmox VE API client as there is an unknown configuration value "+
				"for the API endpoint. Either target apply the source of the value first, set the value statically in "+
				"the configuration, or use the PROXMOX_VE_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	// Check environment variables
	apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN")
	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	insecure := utils.GetAnyBoolEnv("PROXMOX_VE_INSECURE")
	minTLS := utils.GetAnyStringEnv("PROXMOX_VE_MIN_TLS")
	username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME")
	password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD")

	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	if !config.MinTLS.IsNull() {
		minTLS = config.MinTLS.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Proxmox VE API Endpoint",
			"The provider cannot create the Proxmox VE API client as there is a missing or empty value for the API endpoint. "+
				"Set the host value in the configuration or use the PROXMOX_VE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Proxmox VE API client

	creds, err := api.NewCredentials(username, password, "", apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Proxmox VE API credentials",
			err.Error(),
		)
	}

	conn, err := api.NewConnection(
		endpoint,
		insecure,
		minTLS,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Proxmox VE API connection",
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := api.NewClient(creds, conn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Proxmox VE API client",
			err.Error(),
		)
	}

	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME")
	sshPassword := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PASSWORD")
	sshAgent := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT")
	sshPrivateKey := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY")
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK")
	sshSocks5Server := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_SERVER")
	sshSocks5Username := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_USERNAME")
	sshSocks5Password := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_PASSWORD")
	nodeOverrides := map[string]ssh.ProxmoxNode{}

	//nolint: nestif
	if len(config.SSH) > 0 {
		if !config.SSH[0].Username.IsNull() {
			sshUsername = config.SSH[0].Username.ValueString()
		}

		if !config.SSH[0].Password.IsNull() {
			sshPassword = config.SSH[0].Password.ValueString()
		}

		if !config.SSH[0].Agent.IsNull() {
			sshAgent = config.SSH[0].Agent.ValueBool()
		}

		if !config.SSH[0].AgentSocket.IsNull() {
			sshAgentSocket = config.SSH[0].AgentSocket.ValueString()
		}

		if !config.SSH[0].PrivateKey.IsNull() {
			sshPrivateKey = config.SSH[0].PrivateKey.ValueString()
		}

		if !config.SSH[0].Socks5Server.IsNull() {
			sshSocks5Server = config.SSH[0].Socks5Server.ValueString()
		}

		if !config.SSH[0].Socks5Username.IsNull() {
			sshSocks5Username = config.SSH[0].Socks5Username.ValueString()
		}

		if !config.SSH[0].Socks5Password.IsNull() {
			sshSocks5Password = config.SSH[0].Socks5Password.ValueString()
		}

		for _, n := range config.SSH[0].Nodes {
			nodePort := int32(n.Port.ValueInt64())
			if nodePort == 0 {
				nodePort = 22
			}

			nodeOverrides[n.Name.ValueString()] = ssh.ProxmoxNode{
				Address: n.Address.ValueString(),
				Port:    nodePort,
			}
		}
	}

	if sshUsername == "" {
		sshUsername = strings.Split(creds.Username, "@")[0]
	}

	if sshPassword == "" {
		sshPassword = creds.Password
	}

	sshClient, err := ssh.NewClient(
		sshUsername, sshPassword, sshAgent, sshAgentSocket, sshPrivateKey,
		sshSocks5Server, sshSocks5Username, sshSocks5Password,
		&apiResolverWithOverrides{
			ar:        apiResolver{c: apiClient},
			overrides: nodeOverrides,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Proxmox VE SSH client",
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Intentionally use 'PROXMOX_VE_TMPDIR' with 'TMP' instead of 'TEMP', to match os.TempDir's use of $TMPDIR
	tmpDirOverride := utils.GetAnyStringEnv("PROXMOX_VE_TMPDIR", "PM_VE_TMPDIR")

	if !config.TmpDir.IsNull() {
		tmpDirOverride = config.TmpDir.ValueString()
	}

	client := proxmox.NewClient(apiClient, sshClient, tmpDirOverride)

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *proxmoxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterOptionsResource,
		NewDownloadFileResource,
		acme.NewACMEAccountResource,
		acme.NewACMEPluginResource,
		apt.NewRepositoryResource,
		apt.NewStandardRepositoryResource,
		access.NewACLResource,
		access.NewUserTokenResource,
		ha.NewHAGroupResource,
		ha.NewHAResourceResource,
		hardwaremapping.NewPCIResource,
		hardwaremapping.NewUSBResource,
		network.NewLinuxBridgeResource,
		network.NewLinuxVLANResource,
		vm.NewResource,
	}
}

func (p *proxmoxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVersionDataSource,
		acme.NewACMEAccountsDataSource,
		acme.NewACMEAccountDataSource,
		acme.NewACMEPluginsDataSource,
		acme.NewACMEPluginDataSource,
		apt.NewRepositoryDataSource,
		apt.NewStandardRepositoryDataSource,
		ha.NewHAGroupDataSource,
		ha.NewHAGroupsDataSource,
		ha.NewHAResourceDataSource,
		ha.NewHAResourcesDataSource,
		hardwaremapping.NewDataSource,
		hardwaremapping.NewPCIDataSource,
		hardwaremapping.NewUSBDataSource,
		vm.NewDataSource,
	}
}

type apiResolver struct {
	c api.Client
}

func (r *apiResolver) Resolve(ctx context.Context, nodeName string) (ssh.ProxmoxNode, error) {
	nc := &nodes.Client{Client: r.c, NodeName: nodeName}

	networkDevices, err := nc.ListNetworkInterfaces(ctx)
	if err != nil {
		return ssh.ProxmoxNode{}, fmt.Errorf("failed to list network devices of node %q: %w", nc.NodeName, err)
	}

	nodeAddress := ""

	// try IPv4 address on the interface with IPv4 gateway
	tflog.Debug(ctx, "Attempting to find interfaces with both a static IPV4 address and gateway.")

	for _, d := range networkDevices {
		if d.Gateway != nil && d.Address != nil {
			nodeAddress = *d.Address
			break
		}
	}

	if nodeAddress == "" {
		// fallback 1: try IPv6 address on the interface with IPv6 gateway
		tflog.Debug(ctx, "Attempting to find interfaces with both a static IPV6 address and gateway.")

		for _, d := range networkDevices {
			if d.Gateway6 != nil && d.Address6 != nil {
				nodeAddress = *d.Address6
				break
			}
		}
	}

	if nodeAddress == "" {
		// fallback 2: use first interface with any IPv4 address
		tflog.Debug(ctx, "Attempting to find interfaces with at least a static IPV4 address.")

		for _, d := range networkDevices {
			if d.Address != nil {
				nodeAddress = *d.Address
				break
			}
		}
	}

	if nodeAddress == "" {
		// fallback 3: do a good old DNS lookup
		tflog.Debug(ctx, fmt.Sprintf("Attempting a DNS lookup of node %q.", nc.NodeName))

		ips, err := net.LookupIP(nodeName)
		if err == nil {
			for _, ip := range ips {
				if ipv4 := ip.To4(); ipv4 != nil {
					nodeAddress = ipv4.String()
					break
				}
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Failed to do a DNS lookup of the node: %s", err.Error()))
		}
	}

	if nodeAddress == "" {
		return ssh.ProxmoxNode{}, fmt.Errorf("failed to determine the IP address of node \"%s\"", nc.NodeName)
	}

	nodeAddressParts := strings.Split(nodeAddress, "/")
	node := ssh.ProxmoxNode{Address: nodeAddressParts[0], Port: 22}

	return node, nil
}

type apiResolverWithOverrides struct {
	ar        apiResolver
	overrides map[string]ssh.ProxmoxNode
}

func (r *apiResolverWithOverrides) Resolve(ctx context.Context, nodeName string) (ssh.ProxmoxNode, error) {
	if node, ok := r.overrides[nodeName]; ok {
		return node, nil
	}

	return r.ar.Resolve(ctx, nodeName)
}
