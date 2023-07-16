/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/bpg/proxmox-api/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/proxmox-api/nodes"
	"github.com/bpg/proxmox-api/ssh"

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

	var restClient rest.Client

	var sshClient ssh.Client

	var creds *rest.Credentials

	var conn *rest.Connection

	// Check environment variables
	apiToken := utils.GetAnyStringEnv("PROXMOX_VE_API_TOKEN", "PM_VE_API_TOKEN")
	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT", "PM_VE_ENDPOINT")
	insecure := utils.GetAnyBoolEnv("PROXMOX_VE_INSECURE", "PM_VE_INSECURE")
	username := utils.GetAnyStringEnv("PROXMOX_VE_USERNAME", "PM_VE_USERNAME")
	password := utils.GetAnyStringEnv("PROXMOX_VE_PASSWORD", "PM_VE_PASSWORD")
	otp := utils.GetAnyStringEnv("PROXMOX_VE_OTP", "PM_VE_OTP")

	if v, ok := d.GetOk(mkProviderAPIToken); ok {
		apiToken = v.(string)
	}

	if v, ok := d.GetOk(mkProviderEndpoint); ok {
		endpoint = v.(string)
	}

	if v, ok := d.GetOk(mkProviderInsecure); ok {
		insecure = v.(bool)
	}

	if v, ok := d.GetOk(mkProviderUsername); ok {
		username = v.(string)
	}

	if v, ok := d.GetOk(mkProviderPassword); ok {
		password = v.(string)
	}

	if v, ok := d.GetOk(mkProviderOTP); ok {
		otp = v.(string)
	}

	creds, err = rest.NewCredentials(username, password, otp, apiToken)
	diags = append(diags, diag.FromErr(err)...)

	conn, err = rest.NewConnection(endpoint, insecure)
	diags = append(diags, diag.FromErr(err)...)

	if diags.HasError() {
		return nil, diags
	}

	restClient, err = rest.NewClient(creds, conn)
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

	if v, ok := sshConf[mkProviderSSHUsername]; !ok || v.(string) == "" {
		if sshUsername != "" {
			sshConf[mkProviderSSHUsername] = sshUsername
		} else {
			sshConf[mkProviderSSHUsername] = strings.Split(creds.Username, "@")[0]
		}
	}

	if v, ok := sshConf[mkProviderSSHPassword]; !ok || v.(string) == "" {
		if sshPassword != "" {
			sshConf[mkProviderSSHPassword] = sshPassword
		} else {
			sshConf[mkProviderSSHPassword] = creds.Password
		}
	}

	if _, ok := sshConf[mkProviderSSHAgent]; !ok {
		sshConf[mkProviderSSHAgent] = sshAgent
	}

	if _, ok := sshConf[mkProviderSSHAgentSocket]; !ok {
		sshConf[mkProviderSSHAgentSocket] = sshAgentSocket
	}

	nodeOverrides := map[string]string{}

	if ns, ok := sshConf[mkProviderSSHNode]; ok {
		for _, n := range ns.([]interface{}) {
			node := n.(map[string]interface{})
			nodeOverrides[node[mkProviderSSHNodeName].(string)] = node[mkProviderSSHNodeAddress].(string)
		}
	}

	sshClient, err = ssh.NewClient(
		sshConf[mkProviderSSHUsername].(string),
		sshConf[mkProviderSSHPassword].(string),
		sshConf[mkProviderSSHAgent].(bool),
		sshConf[mkProviderSSHAgentSocket].(string),
		&apiResolverWithOverrides{
			ar:        apiResolver{c: restClient},
			overrides: nodeOverrides,
		},
	)
	if err != nil {
		return nil, diag.Errorf("error creating SSH client: %s", err)
	}

	config := proxmoxtf.NewProviderConfiguration(restClient, sshClient)

	return config, nil
}

type apiResolver struct {
	c rest.Client
}

func (r *apiResolver) Resolve(ctx context.Context, nodeName string) (string, error) {
	nc := &nodes.Client{Client: r.c, NodeName: nodeName}

	networkDevices, err := nc.ListNetworkInterfaces(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list network devices of node \"%s\": %w", nc.NodeName, err)
	}

	nodeAddress := ""

	for _, d := range networkDevices {
		if d.Address != nil {
			nodeAddress = *d.Address
			break
		}
	}

	if nodeAddress == "" {
		return "", fmt.Errorf("failed to determine the IP address of node \"%s\"", nc.NodeName)
	}

	nodeAddressParts := strings.Split(nodeAddress, "/")

	return nodeAddressParts[0], nil
}

type apiResolverWithOverrides struct {
	ar        apiResolver
	overrides map[string]string
}

func (r *apiResolverWithOverrides) Resolve(ctx context.Context, nodeName string) (string, error) {
	if ip, ok := r.overrides[nodeName]; ok {
		return ip, nil
	}

	return r.ar.Resolve(ctx, nodeName)
}
