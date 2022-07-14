/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentHostsAddresses        = "addresses"
	mkDataSourceVirtualEnvironmentHostsDigest           = "digest"
	mkDataSourceVirtualEnvironmentHostsEntries          = "entries"
	mkDataSourceVirtualEnvironmentHostsEntriesAddress   = "address"
	mkDataSourceVirtualEnvironmentHostsEntriesHostnames = "hostnames"
	mkDataSourceVirtualEnvironmentHostsHostnames        = "hostnames"
	mkDataSourceVirtualEnvironmentHostsNodeName         = "node_name"
)

func dataSourceVirtualEnvironmentHosts() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentHostsAddresses: {
				Type:        schema.TypeList,
				Description: "The addresses",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentHostsDigest: {
				Type:        schema.TypeString,
				Description: "The SHA1 digest",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentHostsEntries: {
				Type:        schema.TypeList,
				Description: "The host entries",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentHostsEntriesAddress: {
							Type:        schema.TypeString,
							Description: "The address",
							Computed:    true,
						},
						mkDataSourceVirtualEnvironmentHostsEntriesHostnames: {
							Type:        schema.TypeList,
							Description: "The hostnames",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentHostsHostnames: {
				Type:        schema.TypeList,
				Description: "The hostnames",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkDataSourceVirtualEnvironmentHostsNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
		},
		ReadContext: dataSourceVirtualEnvironmentHostsRead,
	}
}

func dataSourceVirtualEnvironmentHostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentHostsNodeName).(string)
	hosts, err := veClient.GetHosts(ctx, nodeName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s_hosts", nodeName))

	// Parse the entries in the hosts file.
	var addresses []interface{}
	var entries []interface{}
	var hostnames []interface{}
	lines := strings.Split(hosts.Data, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.ReplaceAll(line, "\t", " ")
		values := strings.Split(line, " ")

		if values[0] == "" {
			continue
		}

		addresses = append(addresses, values[0])
		entry := map[string]interface{}{}
		var hostnamesForAddress []interface{}

		for _, hostname := range values[1:] {
			if hostname != "" {
				hostnamesForAddress = append(hostnamesForAddress, hostname)
			}
		}

		entry[mkDataSourceVirtualEnvironmentHostsEntriesAddress] = values[0]
		entry[mkDataSourceVirtualEnvironmentHostsEntriesHostnames] = hostnamesForAddress

		entries = append(entries, entry)
		hostnames = append(hostnames, hostnamesForAddress)
	}

	err = d.Set(mkDataSourceVirtualEnvironmentHostsAddresses, addresses)
	diags = append(diags, diag.FromErr(err)...)

	if hosts.Digest != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentHostsDigest, *hosts.Digest)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentHostsDigest, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentHostsEntries, entries)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentHostsHostnames, hostnames)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
