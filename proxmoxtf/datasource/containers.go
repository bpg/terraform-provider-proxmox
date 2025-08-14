/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	proxmoxapi "github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentContainers = "containers"
)

// Containers returns a resource for the Proxmox Containers.
//
//nolint:dupl // TODO: refactor to avoid duplication
func Containers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentContainerNodeName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The node name. All cluster nodes will be queried in case this is omitted",
			},
			mkDataSourceVirtualEnvironmentContainerTags: {
				Type:        schema.TypeList,
				Description: "Tags of the Container to match",
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
			mkDataSourceVirtualEnvironmentContainers: {
				Type:        schema.TypeList,
				Description: "Containers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: Container().Schema,
				},
			},
		},
		ReadContext: containersRead,
	}
}

// containersRead reads the Containers.
func containersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	tagsData := d.Get(mkDataSourceVirtualEnvironmentContainerTags).([]interface{})
	for _, tagData := range tagsData {
		tag := strings.TrimSpace(tagData.(string))
		if len(tag) > 0 {
			filterTags = append(filterTags, tag)
		}
	}

	sort.Strings(filterTags)

	filters := d.Get(mkDataSourceFilter).([]interface{})

	var containers []interface{}

	for _, nodeName := range nodeNames {
		listData, e := api.Node(nodeName).Container(0).ListContainers(ctx)
		if e != nil {
			var httpError *proxmoxapi.HTTPError
			if errors.As(e, &httpError) && httpError.Code == 595 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("node %q is not available - Container list may be incomplete", nodeName),
				})

				continue
			}

			diags = append(diags, diag.FromErr(e)...)
		}

		sort.Slice(listData, func(i, j int) bool {
			return listData[i].VMID < listData[j].VMID
		})

		for _, data := range listData {
			container := map[string]interface{}{
				mkDataSourceVirtualEnvironmentContainerNodeName: nodeName,
				mkDataSourceVirtualEnvironmentContainerVMID:     data.VMID,
			}

			if data.Name != nil {
				container[mkDataSourceVirtualEnvironmentContainerName] = *data.Name
			} else {
				container[mkDataSourceVirtualEnvironmentContainerName] = ""
			}

			var tags []string
			if data.Tags != nil && *data.Tags != "" {
				tags = strings.Split(*data.Tags, ";")
				sort.Strings(tags)
				container[mkDataSourceVirtualEnvironmentContainerTags] = tags
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
				container[mkDataSourceVirtualEnvironmentContainerTemplate] = true
			} else {
				container[mkDataSourceVirtualEnvironmentContainerTemplate] = false
			}

			container[mkDataSourceVirtualEnvironmentContainerStatus] = *data.Status

			if len(filters) > 0 {
				allFiltersMatched, err := checkContainerMatchFilters(container, filters)
				diags = append(diags, diag.FromErr(err)...)

				if !allFiltersMatched {
					continue
				}
			}

			containers = append(containers, container)
		}
	}

	err = d.Set(mkDataSourceVirtualEnvironmentContainers, containers)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(uuid.New().String())

	return diags
}

//nolint:dupl // TODO: refactor to avoid duplication
func checkContainerMatchFilters(container map[string]interface{}, filters []interface{}) (bool, error) {
	for _, v := range filters {
		filter := v.(map[string]interface{})
		filterName := filter["name"]
		filterValues := filter["values"].([]interface{})
		filterRegex := filter["regex"].(bool)

		var containerValue string

		switch filterName {
		case "template":
			containerValue = strconv.FormatBool(container[mkDataSourceVirtualEnvironmentContainerTemplate].(bool))
		case "status":
			containerValue = container[mkDataSourceVirtualEnvironmentContainerStatus].(string)
		case "name":
			containerValue = container[mkDataSourceVirtualEnvironmentContainerName].(string)
		case "node_name":
			containerValue = container[mkDataSourceVirtualEnvironmentContainerNodeName].(string)
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

				if r.MatchString(containerValue) {
					atLeastOneValueMatched = true
					break
				}
			} else if filterValue == containerValue {
				atLeastOneValueMatched = true
				break
			}
		}

		if !atLeastOneValueMatched {
			return false, nil
		}
	}

	return true, nil
}
