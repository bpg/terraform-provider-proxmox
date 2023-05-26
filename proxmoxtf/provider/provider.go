/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvProviderOTP                = ""
	mkProviderVirtualEnvironment = "virtual_environment"
	mkProviderEndpoint           = "endpoint"
	mkProviderInsecure           = "insecure"
	mkProviderOTP                = "otp"
	mkProviderPassword           = "password"
	mkProviderUsername           = "username"
	mkProviderSSH                = "ssh"
	mkProviderSSHUsername        = "username"
	mkProviderSSHPassword        = "password"
	mkProviderSSHAgent           = "agent"
	mkProviderSSHAgentSocket     = "agent_socket"
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

	var apiClient api.Client

	var sshClient ssh.Client

	var username, password string

	// Legacy configuration, wrapped in the deprecated `virtual_environment` block
	veConfigBlock := d.Get(mkProviderVirtualEnvironment).([]interface{})
	if len(veConfigBlock) > 0 {
		veConfig := veConfigBlock[0].(map[string]interface{})

		username = veConfig[mkProviderUsername].(string)
		password = veConfig[mkProviderPassword].(string)

		apiClient, err = api.NewClient(
			veConfig[mkProviderEndpoint].(string),
			username,
			password,
			veConfig[mkProviderOTP].(string),
			veConfig[mkProviderInsecure].(bool),
		)
	} else {
		username = d.Get(mkProviderUsername).(string)
		password = d.Get(mkProviderPassword).(string)

		apiClient, err = api.NewClient(
			d.Get(mkProviderEndpoint).(string),
			username,
			password,
			d.Get(mkProviderOTP).(string),
			d.Get(mkProviderInsecure).(bool),
		)
	}

	if err != nil {
		return nil, diag.Errorf("error creating virtual environment client: %s", err)
	}

	sshConf := map[string]interface{}{}

	sshBlock := d.Get(mkProviderSSH).([]interface{})
	if len(sshBlock) > 0 {
		sshConf = sshBlock[0].(map[string]interface{})
	}

	if v, ok := sshConf[mkProviderSSHUsername]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHUsername] = strings.Split(username, "@")[0]
	}

	if v, ok := sshConf[mkProviderSSHPassword]; !ok || v.(string) == "" {
		sshConf[mkProviderSSHPassword] = password
	}

	if _, ok := sshConf[mkProviderSSHAgent]; !ok {
		sshConf[mkProviderSSHAgent] = false
	}

	if _, ok := sshConf[mkProviderSSHAgentSocket]; !ok {
		sshConf[mkProviderSSHAgentSocket] = ""
	}

	sshClient, err = ssh.NewClient(
		sshConf[mkProviderSSHUsername].(string),
		sshConf[mkProviderSSHPassword].(string),
		sshConf[mkProviderSSHAgent].(bool),
		sshConf[mkProviderSSHAgentSocket].(string),
	)
	if err != nil {
		return nil, diag.Errorf("error creating SSH client: %s", err)
	}

	config := proxmoxtf.NewProviderConfiguration(apiClient, sshClient)

	return config, nil
}
