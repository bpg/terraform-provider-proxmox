package cpu

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
)

const (
	dvCPUArchitecture = "x86_64"
	dvCPUCores        = 1
	dvCPUHotplugged   = 0
	dvCPULimit        = 0
	dvCPUNUMA         = false
	dvCPUSockets      = 1
	dvCPUType         = "qemu64"
	dvCPUUnits        = 1024
	dvCPUAffinity     = ""

	MkCPU             = "cpu"
	mkCPUArchitecture = "architecture"
	mkCPUCores        = "cores"
	mkCPUFlags        = "flags"
	mkCPUHotplugged   = "hotplugged"
	mkCPULimit        = "limit"
	mkCPUNUMA         = "numa"
	mkCPUSockets      = "sockets"
	mkCPUType         = "type"
	mkCPUUnits        = "units"
	mkCPUAffinity     = "affinity"
)

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MkCPU: {
			Type:        schema.TypeList,
			Description: "The CPU allocation",
			Optional:    true,
			Computed:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{
					map[string]interface{}{
						mkCPUArchitecture: dvCPUArchitecture,
						mkCPUCores:        dvCPUCores,
						mkCPUFlags:        []interface{}{},
						mkCPUHotplugged:   dvCPUHotplugged,
						mkCPULimit:        dvCPULimit,
						mkCPUNUMA:         dvCPUNUMA,
						mkCPUSockets:      dvCPUSockets,
						mkCPUType:         dvCPUType,
						mkCPUUnits:        dvCPUUnits,
						mkCPUAffinity:     dvCPUAffinity,
					},
				}, nil
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkCPUArchitecture: {
						Type:             schema.TypeString,
						Description:      "The CPU architecture",
						Optional:         true,
						Default:          dvCPUArchitecture,
						ValidateDiagFunc: validators.CPUArchitectureValidator(),
					},
					mkCPUCores: {
						Type:             schema.TypeInt,
						Description:      "The number of CPU cores",
						Optional:         true,
						Computed:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 2304)),
					},
					mkCPUFlags: {
						Type:        schema.TypeList,
						Description: "The CPU flags",
						Optional:    true,
						Computed:    true,
						//DefaultFunc: func() (interface{}, error) {
						//	return []interface{}{}, nil
						//},
						Elem: &schema.Schema{Type: schema.TypeString},
					},
					mkCPUHotplugged: {
						Type:             schema.TypeInt,
						Description:      "The number of hotplugged vCPUs",
						Optional:         true,
						Default:          dvCPUHotplugged,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 2304)),
					},
					mkCPULimit: {
						Type:        schema.TypeInt,
						Description: "Limit of CPU usage",
						Optional:    true,
						Default:     dvCPULimit,
						ValidateDiagFunc: validation.ToDiagFunc(
							validation.IntBetween(0, 128),
						),
					},
					mkCPUNUMA: {
						Type:        schema.TypeBool,
						Description: "Enable/disable NUMA.",
						Optional:    true,
						Default:     dvCPUNUMA,
					},
					mkCPUSockets: {
						Type:             schema.TypeInt,
						Description:      "The number of CPU sockets",
						Optional:         true,
						Default:          dvCPUSockets,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 16)),
					},
					mkCPUType: {
						Type:             schema.TypeString,
						Description:      "The emulated CPU type",
						Optional:         true,
						Computed:         true,
						ValidateDiagFunc: validators.CPUTypeValidator(),
					},
					mkCPUUnits: {
						Type:        schema.TypeInt,
						Description: "The CPU units",
						Optional:    true,
						Default:     dvCPUUnits,
						ValidateDiagFunc: validation.ToDiagFunc(
							validation.IntBetween(2, 262144),
						),
					},
					mkCPUAffinity: {
						Type:             schema.TypeString,
						Description:      "The CPU affinity",
						Optional:         true,
						Default:          dvCPUAffinity,
						ValidateDiagFunc: validators.CPUAffinityValidator(),
					},
				},
			},
			MaxItems: 1,
			MinItems: 0,
		},
	}
}
