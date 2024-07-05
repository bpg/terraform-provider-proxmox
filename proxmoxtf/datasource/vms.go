/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentVMs = "vms"
	mkDataSourceFilter                = "filter"
	mkDataSourceFilterName            = "name"
	mkDataSourceFilterValues          = "values"
	mkDataSourceFilterRegex           = "regex"
)

// VMs returns a resource for the Proxmox VMs.
func VMs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVMNodeName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The node name. All cluster nodes will be queried in case this is omitted",
			},
			mkDataSourceVirtualEnvironmentVMTags: {
				Type:        schema.TypeList,
				Description: "Tags of the VM to match",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceFilter: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter blocks",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceFilterName: {
							Type:        schema.TypeString,
							Description: "Attribute to filter on. One of [name, template, status, node_name]",
							Required:    true,
						},
						mkDataSourceFilterValues: {
							Type:        schema.TypeList,
							Description: "List of values to pass the filter (OR logic)",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						mkDataSourceFilterRegex: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Treat values as regex patterns",
						},
					},
				},
			},
			mkDataSourceVirtualEnvironmentVMs: {
				Type:        schema.TypeList,
				Description: "VMs",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: VM().Schema,
				},
			},
		},
		ReadContext: vmsRead,
	}
}

// vmRead reads the data of a VM by ID.
func vmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeNames, err := getNodeNames(ctx, d, api)
	if err != nil {
		return diag.FromErr(err)
	}

	var filterTags []string

	tagsData := d.Get(mkDataSourceVirtualEnvironmentVMTags).([]interface{})
	for _, tagData := range tagsData {
		tag := strings.TrimSpace(tagData.(string))
		if len(tag) > 0 {
			filterTags = append(filterTags, tag)
		}
	}

	sort.Strings(filterTags)

	filters := d.Get(mkDataSourceFilter).([]interface{})

	var vms []interface{}

	for _, nodeName := range nodeNames {
		listData, e := api.Node(nodeName).VM(0).ListVMs(ctx)
		if e != nil {
			diags = append(diags, diag.FromErr(e)...)
		}

		sort.Slice(listData, func(i, j int) bool {
			return listData[i].VMID < listData[j].VMID
		})

		for _, data := range listData {
			vm := map[string]interface{}{
				mkDataSourceVirtualEnvironmentVMNodeName: nodeName,
				mkDataSourceVirtualEnvironmentVMVMID:     data.VMID,
			}

			if data.Name != nil {
				vm[mkDataSourceVirtualEnvironmentVMName] = *data.Name
			} else {
				vm[mkDataSourceVirtualEnvironmentVMName] = ""
			}

			var tags []string
			if data.Tags != nil && *data.Tags != "" {
				tags = strings.Split(*data.Tags, ";")
				sort.Strings(tags)
				vm[mkDataSourceVirtualEnvironmentVMTags] = tags
			}

			if len(filterTags) > 0 {
				match := true

				for _, tag := range filterTags {
					if !slices.Contains(tags, tag) {
						match = false
						break
					}
				}

				if !match {
					continue
				}
			}

			if data.Template != (*types.CustomBool)(nil) && *data.Template {
				vm[mkDataSourceVirtualEnvironmentVMTemplate] = true
			} else {
				vm[mkDataSourceVirtualEnvironmentVMTemplate] = false
			}

			vm[mkDataSourceVirtualEnvironmentVMStatus] = *data.Status

			if len(filters) > 0 {
				allFiltersMatched, err := checkVMMatchFilters(vm, filters)
				diags = append(diags, diag.FromErr(err)...)

				if !allFiltersMatched {
					continue
				}
			}

			vms = append(vms, vm)
		}
	}

	err = d.Set(mkDataSourceVirtualEnvironmentVMs, vms)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(uuid.New().String())

	return diags
}

func getNodeNames(ctx context.Context, d *schema.ResourceData, api proxmox.Client) ([]string, error) {
	var nodeNames []string

	nodeName := d.Get(mkDataSourceVirtualEnvironmentVMNodeName).(string)
	if nodeName != "" {
		nodeNames = append(nodeNames, nodeName)
	} else {
		nodes, err := api.Node(nodeName).ListNodes(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing nodes: %w", err)
		}

		for _, node := range nodes {
			nodeNames = append(nodeNames, node.Name)
		}
	}

	sort.Strings(nodeNames)

	return nodeNames, nil
}

func checkVMMatchFilters(vm map[string]interface{}, filters []interface{}) (bool, error) {
	for _, v := range filters {
		filter := v.(map[string]interface{})
		filterName := filter["name"]
		filterValues := filter["values"].([]interface{})
		filterRegex := filter["regex"].(bool)

		var vmValue string

		switch filterName {
		case "template":
			vmValue = strconv.FormatBool(vm[mkDataSourceVirtualEnvironmentVMTemplate].(bool))
		case "status":
			vmValue = vm[mkDataSourceVirtualEnvironmentVMStatus].(string)
		case "name":
			vmValue = vm[mkDataSourceVirtualEnvironmentVMName].(string)
		case "node_name":
			vmValue = vm[mkDataSourceVirtualEnvironmentVMNodeName].(string)
		default:
			return false, fmt.Errorf(
				"unknown filter name '%s', should be one of [name, template, status, node_name]",
				filterName,
			)
		}

		atLeastOneValueMatched := false

		for _, filterValue := range filterValues {
			if filterRegex {
				r, err := regexp.Compile(filterValue.(string))
				if err != nil {
					return false, fmt.Errorf("error interpreting filter '%s' value '%s' as regex: %w", filterName, filterValue, err)
				}

				if r.MatchString(vmValue) {
					atLeastOneValueMatched = true
					break
				}
			} else {
				if filterValue == vmValue {
					atLeastOneValueMatched = true
					break
				}
			}
		}

		if !atLeastOneValueMatched {
			return false, nil
		}
	}

	return true, nil
}
