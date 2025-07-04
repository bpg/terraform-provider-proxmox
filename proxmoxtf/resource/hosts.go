/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkResourceVirtualEnvironmentHostsAddresses        = "addresses"
	mkResourceVirtualEnvironmentHostsDigest           = "digest"
	mkResourceVirtualEnvironmentHostsEntries          = "entries"
	mkResourceVirtualEnvironmentHostsEntriesAddress   = "address"
	mkResourceVirtualEnvironmentHostsEntriesHostnames = "hostnames"
	mkResourceVirtualEnvironmentHostsEntry            = "entry"
	mkResourceVirtualEnvironmentHostsEntryAddress     = "address"
	mkResourceVirtualEnvironmentHostsEntryHostnames   = "hostnames"
	mkResourceVirtualEnvironmentHostsHostnames        = "hostnames"
	mkResourceVirtualEnvironmentHostsNodeName         = "node_name"
)

// Hosts returns a resource that manages hosts settings for a node.
func Hosts() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentHostsAddresses: {
				Type:        schema.TypeList,
				Description: "The addresses",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentHostsDigest: {
				Type:        schema.TypeString,
				Description: "The SHA1 digest",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentHostsEntries: {
				Type:        schema.TypeList,
				Description: "The host entries",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentHostsEntriesAddress: {
							Type:        schema.TypeString,
							Description: "The address",
							Computed:    true,
						},
						mkResourceVirtualEnvironmentHostsEntriesHostnames: {
							Type:        schema.TypeList,
							Description: "The hostnames",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			mkResourceVirtualEnvironmentHostsEntry: {
				Type:        schema.TypeList,
				Description: "The host entries",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentHostsEntryAddress: {
							Type:        schema.TypeString,
							Description: "The address",
							Required:    true,
						},
						mkResourceVirtualEnvironmentHostsEntryHostnames: {
							Type:        schema.TypeList,
							Description: "The hostnames",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MinItems:    1,
						},
					},
				},
				MinItems: 1,
			},
			mkResourceVirtualEnvironmentHostsHostnames: {
				Type:        schema.TypeList,
				Description: "The hostnames",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkResourceVirtualEnvironmentHostsNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
		},
		CreateContext: hostsCreate,
		ReadContext:   hostsRead,
		UpdateContext: hostsUpdate,
		DeleteContext: hostsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				nodeName := d.Id()

				err := d.Set(mkResourceVirtualEnvironmentHostsNodeName, nodeName)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				d.SetId(fmt.Sprintf("%s_hosts", nodeName))

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func hostsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := hostsUpdate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)

	d.SetId(fmt.Sprintf("%s_hosts", nodeName))

	return diags
}

func hostsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)

	hosts, err := api.Node(nodeName).GetHosts(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

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

		entry[mkResourceVirtualEnvironmentHostsEntriesAddress] = values[0]
		entry[mkResourceVirtualEnvironmentHostsEntriesHostnames] = hostnamesForAddress

		entries = append(entries, entry)
		hostnames = append(hostnames, hostnamesForAddress)
	}

	err = d.Set(mkResourceVirtualEnvironmentHostsAddresses, addresses)
	diags = append(diags, diag.FromErr(err)...)

	if hosts.Digest != nil {
		err = d.Set(mkResourceVirtualEnvironmentHostsDigest, *hosts.Digest)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentHostsDigest, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkResourceVirtualEnvironmentHostsEntries, entries)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkResourceVirtualEnvironmentHostsEntry, entries)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkResourceVirtualEnvironmentHostsHostnames, hostnames)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func hostsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	entry := d.Get(mkResourceVirtualEnvironmentHostsEntry).([]interface{})
	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)

	// Generate the data for the hosts file based on the specified entries.
	body := nodes.HostsUpdateRequestBody{
		Data: "",
	}

	for _, e := range entry {
		eMap := e.(map[string]interface{})

		address := eMap[mkResourceVirtualEnvironmentHostsEntryAddress].(string)
		hostnames := eMap[mkResourceVirtualEnvironmentHostsEntryHostnames].([]interface{})

		body.Data += address

		for _, h := range hostnames {
			hostname := h.(string)
			body.Data += " " + hostname
		}

		body.Data += "\n"
	}

	err = api.Node(nodeName).UpdateHosts(ctx, &body)
	if err != nil {
		return diag.FromErr(err)
	}

	return hostsRead(ctx, d, m)
}

func hostsDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
