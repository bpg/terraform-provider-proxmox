/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	dvIPSetCIDRComment = ""
	dvIPSetCIDRNoMatch = false

	mkIPSetName        = "name"
	mkIPSetCIDR        = "cidr"
	mkIPSetCIDRName    = "name"
	mkIPSetCIDRComment = "comment"
	mkIPSetCIDRNoMatch = "nomatch"
)

// IPSet returns a resource to manage IP sets.
func IPSet() *schema.Resource {
	s := map[string]*schema.Schema{
		mkIPSetName: {
			Type:        schema.TypeString,
			Description: "IPSet name",
			Required:    true,
			ForceNew:    false,
		},
		mkIPSetCIDRComment: {
			Type:        schema.TypeString,
			Description: "IPSet comment",
			Optional:    true,
			Default:     dvIPSetCIDRComment,
		},
		mkIPSetCIDR: {
			Type:        schema.TypeList,
			Description: "List of IP or Networks",
			Optional:    true,
			ForceNew:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{}, nil
			},
			DiffSuppressFunc: structure.SuppressIfListsOfMapsAreEqualIgnoringOrderByKey(mkIPSetCIDRName),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkIPSetCIDRName: {
						Type:        schema.TypeString,
						Description: "Network/IP specification in CIDR format",
						Required:    true,
						ForceNew:    true,
					},
					mkIPSetCIDRNoMatch: {
						Type:        schema.TypeBool,
						Description: "No match this IP/CIDR",
						Optional:    true,
						Default:     dvIPSetCIDRNoMatch,
						ForceNew:    true,
					},
					mkIPSetCIDRComment: {
						Type:        schema.TypeString,
						Description: "IP/CIDR comment",
						Optional:    true,
						Default:     dvIPSetCIDRComment,
						ForceNew:    true,
					},
				},
			},
		},
	}

	structure.MergeSchema(s, selectorSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: selectFirewallAPI(ipSetCreate),
		ReadContext:   selectFirewallAPI(ipSetRead),
		UpdateContext: selectFirewallAPI(ipSetUpdate),
		DeleteContext: selectFirewallAPI(ipSetDelete),
		Importer: &schema.ResourceImporter{
			StateContext: ipSetImport,
		},
	}
}

// ipSetImport imports IP sets.
func ipSetImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	switch {
	case strings.HasPrefix(id, "cluster/"):
		name := strings.TrimPrefix(id, "cluster/")
		d.SetId(name)

		err := d.Set(mkIPSetName, name)
		if err != nil {
			return nil, fmt.Errorf("failed setting IPSet name during import: %w", err)
		}
	case strings.HasPrefix(id, "vm/"):
		parts := strings.SplitN(id, "/", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid VM import ID format: %s (expected: vm/<node_name>/<vm_id>/<ipset_name>)", id)
		}

		nodeName := parts[1]

		err := d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting node name during import: %w", err)
		}

		vmID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid VM import ID: VM ID must be a number in %s: %w", id, err)
		}

		err = d.Set(mkSelectorVMID, vmID)
		if err != nil {
			return nil, fmt.Errorf("failed setting VM ID during import: %w", err)
		}

		name := parts[3]
		d.SetId(name)

		err = d.Set(mkIPSetName, name)
		if err != nil {
			return nil, fmt.Errorf("failed setting IPSet name during import: %w", err)
		}
	case strings.HasPrefix(id, "container/"):
		parts := strings.SplitN(id, "/", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid container import ID format: %s (expected: container/<node_name>/<container_id>/<ipset_name>)", id)
		}

		nodeName := parts[1]

		err := d.Set(mkSelectorNodeName, nodeName)
		if err != nil {
			return nil, fmt.Errorf("failed setting node name during import: %w", err)
		}

		containerID, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid container import ID: container ID must be a number in %s: %w", id, err)
		}

		err = d.Set(mkSelectorContainerID, containerID)
		if err != nil {
			return nil, fmt.Errorf("failed setting container ID during import: %w", err)
		}

		name := parts[3]
		d.SetId(name)

		err = d.Set(mkIPSetName, name)
		if err != nil {
			return nil, fmt.Errorf("failed setting IPSet name during import: %w", err)
		}
	default:
		//nolint:lll
		return nil, fmt.Errorf("invalid import ID: %s (expected: 'cluster/<ipset_name>', 'vm/<node_name>/<vm_id>/<ipset_name>', or 'container/<node_name>/<container_id>/<ipset_name>')", id)
	}

	return []*schema.ResourceData{d}, nil
}

func ipSetCreate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkIPSetCIDRComment).(string)
	name := d.Get(mkIPSetName).(string)

	ipSets := d.Get(mkIPSetCIDR).([]interface{})
	ipSetsArray := make(firewall.IPSetContent, len(ipSets))

	for i, v := range ipSets {
		ipSetMap := v.(map[string]interface{})
		ipSetObject := firewall.IPSetGetResponseData{}

		cidr := ipSetMap[mkIPSetCIDRName].(string)
		noMatch := ipSetMap[mkIPSetCIDRNoMatch].(bool)

		comm := ipSetMap[mkIPSetCIDRComment].(string)
		if comm != "" {
			ipSetObject.Comment = &comm
		}

		ipSetObject.CIDR = cidr

		if noMatch {
			noMatchBool := types.CustomBool(true)
			ipSetObject.NoMatch = &noMatchBool
		}

		ipSetsArray[i] = ipSetObject
	}

	body := &firewall.IPSetCreateRequestBody{
		Comment: comment,
		Name:    name,
	}

	err := api.CreateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range ipSetsArray {
		err = api.AddCIDRToIPSet(ctx, name, v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return ipSetRead(ctx, api, d)
}

func ipSetRead(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Id()

	allIPSets, err := api.ListIPSets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range allIPSets {
		if v.Name == name {
			err = d.Set(mkIPSetName, v.Name)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkIPSetCIDRComment, v.Comment)
			diags = append(diags, diag.FromErr(err)...)

			break
		}
	}

	ipSet, err := api.GetIPSetContent(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such IPSet") {
			d.SetId("")
			return nil
		}

		diags = append(diags, diag.FromErr(err)...)

		return diags
	}

	//nolint:prealloc
	var entries []interface{}

	for key := range ipSet {
		entry := map[string]interface{}{}

		entry[mkIPSetCIDRName] = ipSet[key].CIDR
		entry[mkIPSetCIDRNoMatch] = ipSet[key].NoMatch
		entry[mkIPSetCIDRComment] = ipSet[key].Comment

		entries = append(entries, entry)
	}

	err = d.Set(mkIPSetCIDR, entries)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func ipSetUpdate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkIPSetCIDRComment).(string)
	newName := d.Get(mkIPSetName).(string)
	previousName := d.Id()

	body := &firewall.IPSetUpdateRequestBody{
		ReName:  previousName,
		Name:    newName,
		Comment: &comment,
	}

	err := api.UpdateIPSet(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return ipSetRead(ctx, api, d)
}

func ipSetDelete(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	diags := diag.Diagnostics{}

	name := d.Id()

	IPSetContent, err := api.GetIPSetContent(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// PVE requires content of IPSet be cleared before removal
	if len(IPSetContent) > 0 {
		for _, ipSet := range IPSetContent {
			err = api.DeleteIPSetContent(ctx, name, ipSet.CIDR)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	err = api.DeleteIPSet(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such IPSet") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
