/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
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

func resourceVirtualEnvironmentHosts() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentHostsAddresses: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The addresses",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkResourceVirtualEnvironmentHostsDigest: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The SHA1 digest",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentHostsEntries: &schema.Schema{
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
						mkResourceVirtualEnvironmentHostsEntriesHostnames: &schema.Schema{
							Type:        schema.TypeList,
							Description: "The hostnames",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			mkResourceVirtualEnvironmentHostsEntry: &schema.Schema{
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
						mkResourceVirtualEnvironmentHostsEntryHostnames: &schema.Schema{
							Type:        schema.TypeList,
							Description: "The hostnames",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MinItems:    1,
						},
					},
				},
			},
			mkResourceVirtualEnvironmentHostsHostnames: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The hostnames",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkResourceVirtualEnvironmentHostsNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
		},
		Create: resourceVirtualEnvironmentHostsCreate,
		Read:   resourceVirtualEnvironmentHostsRead,
		Update: resourceVirtualEnvironmentHostsUpdate,
		Delete: resourceVirtualEnvironmentHostsDelete,
	}
}

func resourceVirtualEnvironmentHostsCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceVirtualEnvironmentHostsUpdate(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)

	d.SetId(fmt.Sprintf("%s_hosts", nodeName))

	return nil
}

func resourceVirtualEnvironmentHostsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)
	hosts, err := veClient.GetHosts(nodeName)

	if err != nil {
		return err
	}

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

		entry[mkResourceVirtualEnvironmentHostsEntriesAddress] = values[0]
		entry[mkResourceVirtualEnvironmentHostsEntriesHostnames] = hostnamesForAddress

		entries = append(entries, entry)
		hostnames = append(hostnames, hostnamesForAddress)
	}

	d.Set(mkResourceVirtualEnvironmentHostsAddresses, addresses)

	if hosts.Digest != nil {
		d.Set(mkResourceVirtualEnvironmentHostsDigest, *hosts.Digest)
	} else {
		d.Set(mkResourceVirtualEnvironmentHostsDigest, "")
	}

	d.Set(mkResourceVirtualEnvironmentHostsEntries, entries)
	d.Set(mkResourceVirtualEnvironmentHostsEntry, entries)
	d.Set(mkResourceVirtualEnvironmentHostsHostnames, hostnames)

	return nil
}

func resourceVirtualEnvironmentHostsUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	entry := d.Get(mkResourceVirtualEnvironmentHostsEntry).([]interface{})
	nodeName := d.Get(mkResourceVirtualEnvironmentHostsNodeName).(string)

	// Generate the data for the hosts file based on the specified entries.
	body := proxmox.VirtualEnvironmentHostsUpdateRequestBody{
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

	err = veClient.UpdateHosts(nodeName, &body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentHostsRead(d, m)
}

func resourceVirtualEnvironmentHostsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")

	return nil
}
