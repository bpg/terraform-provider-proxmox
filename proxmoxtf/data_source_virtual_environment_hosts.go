/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
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
		Read: dataSourceVirtualEnvironmentHostsRead,
	}
}

func dataSourceVirtualEnvironmentHostsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentHostsNodeName).(string)
	hosts, err := veClient.GetHosts(nodeName)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_hosts", nodeName))

	// Parse the entries in the hosts file.
	addresses := []interface{}{}
	entries := []interface{}{}
	hostnames := []interface{}{}
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
		hostnamesForAddress := []interface{}{}

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

	d.Set(mkDataSourceVirtualEnvironmentHostsAddresses, addresses)

	if hosts.Digest != nil {
		d.Set(mkDataSourceVirtualEnvironmentHostsDigest, *hosts.Digest)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentHostsDigest, "")
	}

	d.Set(mkDataSourceVirtualEnvironmentHostsEntries, entries)
	d.Set(mkDataSourceVirtualEnvironmentHostsHostnames, hostnames)

	return nil
}
