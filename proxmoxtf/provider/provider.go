/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// ProxmoxVirtualEnvironment returns the object for this provider.
func ProxmoxVirtualEnvironment() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap:       createDatasourceMap(),
		ResourcesMap:         createResourceMap(),
		Schema:               createSchema(),
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error

	var diags diag.Diagnostics

	var apiClient api.Client

	var sshClient ssh.Client

	var creds api.Credentials

	var conn *api.Connection

	// Check environment variables
	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT", "PM_VE_ENDPOINT")
	insecure := utils.GetAnyBoolEnv("PROXMOX_VE_INSECURE", "PM_VE_INSECURE")
	minTLS := utils.GetAnyStringEnv("PROXMOX_VE_MIN_TLS", "PM_VE_MIN_TLS")
	authTicket := utils.GetAnyStringEnv("PROXMOX_VE_AUTH_TICKET", "PM_VE_AUTH_TICKET")
	csrfPreventionToken := utils.GetAnyStringEnv("PROXMOX_VE_CSRF_PREVENTION_TOKEN", "PM_VE_CSRF_PREVENTION_TOKEN")
	apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN", "PM_VE_API_TOKEN")
	otp := utils.GetAnyStringEnv("PROXMOX_VE_OTP", "PM_VE_OTP")
	username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME", "PM_VE_USERNAME")
	password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD", "PM_VE_PASSWORD")

	if v, ok := d.GetOk(mkProviderEndpoint); ok {
		endpoint = v.(string)
	}

	if v, ok := d.GetOk(mkProviderInsecure); ok {
		insecure = v.(bool)
	}

	if v, ok := d.GetOk(mkProviderMinTLS); ok {
		minTLS = v.(string)
	}

	if v, ok := d.GetOk(mkProviderAuthTicket); ok {
		authTicket = v.(string)
	}

	if v, ok := d.GetOk(mkProviderCSRFPreventionToken); ok {
		csrfPreventionToken = v.(string)
	}

	if v, ok := d.GetOk(mkProviderAPIToken); ok {
		apiToken = v.(string)
	}

	if v, ok := d.GetOk(mkProviderOTP); ok {
		otp = v.(string)
	}

	if v, ok := d.GetOk(mkProviderUsername); ok {
		username = v.(string)
	}

	if v, ok := d.GetOk(mkProviderPassword); ok {
		password = v.(string)
	}

	creds, err = api.NewCredentials(username, password, otp, apiToken, authTicket, csrfPreventionToken)
	diags = append(diags, diag.FromErr(err)...)

	conn, err = api.NewConnection(endpoint, insecure, minTLS)
	diags = append(diags, diag.FromErr(err)...)

	if diags.HasError() {
		return nil, diags
	}

	apiClient, err = api.NewClient(creds, conn)
	if err != nil {
		return nil, diag.Errorf("error creating virtual environment client: %s", err)
	}

	// ////////////////////////////////////////////////////////////////////////////////////

	sshConf := map[string]interface{}{}

	sshBlock := d.Get(mkProviderSSH).([]interface{})
	if len(sshBlock) > 0 {
		sshConf = sshBlock[0].(map[string]interface{})
	}

	sshUsername := utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME", "PM_VE_SSH_USERNAME")
	sshPassword := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PASSWORD", "PM_VE_SSH_PASSWORD")
	sshAgent := utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT", "PM_VE_SSH_AGENT")
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")
	sshPrivateKey := utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY")
	sshSocks5Server := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_SERVER")
	sshSocks5Username := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_USERNAME")
	sshSocks5Password := utils.GetAnyStringEnv("PROXMOX_VE_SSH_SOCKS5_PASSWORD")

	if v, ok := sshConf[mkProviderSSHUsername]; !ok || v.(string) == "" {
		if sshUsername != "" {
			sshConf[mkProviderSSHUsername] = sshUsername
		} else if creds.UserCredentials != nil {
			sshConf[mkProviderSSHUsername] = strings.Split(creds.UserCredentials.Username, "@")[0]
		}
	}

	if v, ok := sshConf[mkProviderSSHPassword]; !ok || v.(string) == "" {
		if sshPassword != "" {
			sshConf[mkProviderSSHPassword] = sshPassword
		} else if creds.UserCredentials != nil {
			sshConf[mkProviderSSHPassword] = creds.UserCredentials.Password
		}
	}

	if _, ok := sshConf[mkProviderSSHAgent]; !ok {
		sshConf[mkProviderSSHAgent] = sshAgent
	}

	if v, ok := sshConf[mkProviderSSHAgentSocket]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHAgentSocket] = sshAgentSocket
	}

	if v, ok := sshConf[mkProviderSSHPrivateKey]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHPrivateKey] = sshPrivateKey
	}

	if v, ok := sshConf[mkProviderSSHSocks5Server]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHSocks5Server] = sshSocks5Server
	}

	if v, ok := sshConf[mkProviderSSHSocks5Username]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHSocks5Username] = sshSocks5Username
	}

	if v, ok := sshConf[mkProviderSSHSocks5Password]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHSocks5Password] = sshSocks5Password
	}

	nodeOverrides := map[string]ssh.ProxmoxNode{}

	if ns, ok := sshConf[mkProviderSSHNode]; ok {
		for _, n := range ns.([]interface{}) {
			node := n.(map[string]interface{})
			nodeOverrides[node[mkProviderSSHNodeName].(string)] = ssh.ProxmoxNode{
				Address: node[mkProviderSSHNodeAddress].(string),

				Port: int32(node[mkProviderSSHNodePort].(int)),
			}
		}
	}

	sshClient, err = ssh.NewClient(
		sshConf[mkProviderSSHUsername].(string),
		sshConf[mkProviderSSHPassword].(string),
		sshConf[mkProviderSSHAgent].(bool),
		sshConf[mkProviderSSHAgentSocket].(string),
		sshConf[mkProviderSSHPrivateKey].(string),
		sshConf[mkProviderSSHSocks5Server].(string),
		sshConf[mkProviderSSHSocks5Username].(string),
		sshConf[mkProviderSSHSocks5Password].(string),
		&apiResolverWithOverrides{
			ar:        apiResolver{c: apiClient},
			overrides: nodeOverrides,
		},
	)
	if err != nil {
		return nil, diag.Errorf("error creating SSH client: %s", err)
	}

	// Intentionally use 'PROXMOX_VE_TMPDIR' with 'TMP' instead of 'TEMP', to match os.TempDir's use of $TMPDIR
	tmpDirOverride := utils.GetAnyStringEnv("PROXMOX_VE_TMPDIR", "PM_VE_TMPDIR")

	if v, ok := d.GetOk(mkProviderTmpDir); ok {
		tmpDirOverride = v.(string)
	}

	config := proxmoxtf.NewProviderConfiguration(apiClient, sshClient, tmpDirOverride)

	return config, nil
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
