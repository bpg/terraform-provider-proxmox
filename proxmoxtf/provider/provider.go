/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
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
	var veClient *proxmox.VirtualEnvironmentClient

	// Initialize the client for the Virtual Environment, if required.
	veConfigBlock := d.Get(mkProviderVirtualEnvironment).([]interface{})

	tflog.Info(ctx, fmt.Sprintf("veConfigBlock is  %v", veConfigBlock))

	if len(veConfigBlock) > 0 {
		veConfig := veConfigBlock[0].(map[string]interface{})
		veSSHConfig := veConfig[mkProviderSSH].(map[string]interface{})

		veClient, err = proxmox.NewVirtualEnvironmentClient(
			veConfig[mkProviderEndpoint].(string),
			veConfig[mkProviderUsername].(string),
			veConfig[mkProviderSSH].(map[string]interface{})[mkProviderSSHUsername].(string),
			veConfig[mkProviderPassword].(string),
			veConfig[mkProviderInsecure].(bool),
			veSSHConfig[mkProviderSSHUsername].(string),
			veSSHConfig[mkProviderSSHPassword].(string),
			veSSHConfig[mkProviderSSHAgent].(bool),
		)
	} else {

		sshconf := map[string]interface{}{
			mkProviderSSHUsername : "",
			mkProviderSSHPassword : "",
			mkProviderSSHAgent : false,
		}

		sshBlock, sshSet := d.GetOk(mkProviderSSH)
		
		if sshSet {
			sshconf = sshBlock.(*schema.Set).List()[0].(map[string]interface{})
		}

		tflog.Info(ctx, fmt.Sprintf("sshconf is  %v", sshconf))

		veClient, err = proxmox.NewVirtualEnvironmentClient(
			d.Get(mkProviderEndpoint).(string),
			d.Get(mkProviderUsername).(string),
			d.Get(mkProviderPassword).(string),
			d.Get(mkProviderOTP).(string),
			d.Get(mkProviderInsecure).(bool),
			sshconf[mkProviderSSHUsername].(string),
			sshconf[mkProviderSSHPassword].(string),
			sshconf[mkProviderSSHAgent].(bool),
		)
	}

	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := proxmoxtf.NewProviderConfiguration(veClient)

	return config, nil
}
