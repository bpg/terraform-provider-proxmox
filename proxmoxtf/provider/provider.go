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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
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

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error

	var diags diag.Diagnostics

	var apiClient api.Client

	var sshClient ssh.Client

	var creds *api.Credentials

	var conn *api.Connection

	// Legacy configuration, wrapped in the deprecated `virtual_environment` block
	veConfigBlock := d.Get(mkProviderVirtualEnvironment).([]interface{})
	if len(veConfigBlock) > 0 {
		veConfig := veConfigBlock[0].(map[string]interface{})
		creds, err = api.NewCredentials(
			veConfig[mkProviderUsername].(string),
			veConfig[mkProviderPassword].(string),
			veConfig[mkProviderOTP].(string),
			"",
		)
		diags = append(diags, diag.FromErr(err)...)

		conn, err = api.NewConnection(
			veConfig[mkProviderEndpoint].(string),
			veConfig[mkProviderInsecure].(bool),
		)
		diags = append(diags, diag.FromErr(err)...)
	} else {
		creds, err = api.NewCredentials(
			d.Get(mkProviderUsername).(string),
			d.Get(mkProviderPassword).(string),
			d.Get(mkProviderOTP).(string),
			d.Get(mkProviderAPIToken).(string),
		)
		diags = append(diags, diag.FromErr(err)...)

		conn, err = api.NewConnection(
			d.Get(mkProviderEndpoint).(string),
			d.Get(mkProviderInsecure).(bool),
		)
		diags = append(diags, diag.FromErr(err)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	apiClient, err = api.NewClient(ctx, creds, conn)
	if err != nil {
		return nil, diag.Errorf("error creating virtual environment client: %s", err)
	}

	// ////////////////////////////////////////////////////////////////////////////////////

	sshConf := map[string]interface{}{}

	sshBlock := d.Get(mkProviderSSH).([]interface{})
	if len(sshBlock) > 0 {
		sshConf = sshBlock[0].(map[string]interface{})
	}

	if v, ok := sshConf[mkProviderSSHUsername]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHUsername] = strings.Split(creds.Username, "@")[0]
	}

	if v, ok := sshConf[mkProviderSSHPassword]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHPassword] = creds.Password
	}

	if _, ok := sshConf[mkProviderSSHAgent]; !ok {
		sshConf[mkProviderSSHAgent] = false
	}

	if _, ok := sshConf[mkProviderSSHAgentSocket]; !ok {
		sshConf[mkProviderSSHAgentSocket] = ""
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
			ar:        apiResolver{c: apiClient},
			overrides: nodeOverrides,
		},
	)
	if err != nil {
		return nil, diag.Errorf("error creating SSH client: %s", err)
	}

	config := proxmoxtf.NewProviderConfiguration(apiClient, sshClient)

	return config, nil
}

type apiResolver struct {
	c api.Client
}

func (r *apiResolver) Resolve(ctx context.Context, nodeName string) (string, error) {
	nc := &nodes.Client{Client: r.c, NodeName: nodeName}

	ip, err := nc.GetIP(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get node IP: %w", err)
	}

	return ip, nil
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
