package cpu

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

func UpdateWhenClone(d *schema.ResourceData, api proxmox.Client, updateBody *vms.UpdateRequestBody) {
	cpu := d.Get(MkCPU).([]interface{})
	if len(cpu) > 0 && cpu[0] != nil {
		cpuBlock := cpu[0].(map[string]interface{})

		cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
		cpuCores := cpuBlock[mkCPUCores].(int)
		cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
		cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
		cpuLimit := cpuBlock[mkCPULimit].(int)
		cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
		cpuSockets := cpuBlock[mkCPUSockets].(int)
		cpuType := cpuBlock[mkCPUType].(string)
		cpuUnits := cpuBlock[mkCPUUnits].(int)
		cpuAffinity := cpuBlock[mkCPUAffinity].(string)

		cpuFlagsConverted := make([]string, len(cpuFlags))

		for fi, flag := range cpuFlags {
			cpuFlagsConverted[fi] = flag.(string)
		}

		// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
		if api.API().IsRootTicket() ||
			cpuArchitecture != dvCPUArchitecture {
			updateBody.CPUArchitecture = &cpuArchitecture
		}

		if cpuCores != dvCPUCores && cpuCores > 0 {
			// update only if we have non-default & non-empty value
			updateBody.CPUCores = &cpuCores
		}

		updateBody.NUMAEnabled = &cpuNUMA
		updateBody.CPUSockets = &cpuSockets
		updateBody.CPUUnits = &cpuUnits

		if cpuType != dvCPUType && cpuType != "" {
			// update only if we have non-default & non-empty value
			if updateBody.CPUEmulation == nil {
				updateBody.CPUEmulation = &vms.CustomCPUEmulation{}
			}

			updateBody.CPUEmulation.Type = &cpuType
		}

		if len(cpuFlagsConverted) > 0 {
			if updateBody.CPUEmulation == nil {
				updateBody.CPUEmulation = &vms.CustomCPUEmulation{}
			}

			updateBody.CPUEmulation.Flags = &cpuFlagsConverted
		}

		if cpuAffinity != "" {
			updateBody.CPUAffinity = &cpuAffinity
		}

		if cpuHotplugged > 0 {
			updateBody.VirtualCPUCount = &cpuHotplugged
		}

		if cpuLimit > 0 {
			updateBody.CPULimit = &cpuLimit
		}
	}
}

func Create(resource *schema.Resource, d *schema.ResourceData, api proxmox.Client, createBody *vms.CreateRequestBody) error {
	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{MkCPU},
		0,
		true,
	)
	if err != nil {
		return fmt.Errorf("error reading CPU block: %w", err)
	}

	cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
	cpuCores := cpuBlock[mkCPUCores].(int)
	cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
	cpuLimit := cpuBlock[mkCPULimit].(int)
	cpuSockets := cpuBlock[mkCPUSockets].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
	cpuType := cpuBlock[mkCPUType].(string)
	cpuUnits := cpuBlock[mkCPUUnits].(int)
	cpuAffinity := cpuBlock[mkCPUAffinity].(string)

	if cpuCores != dvCPUCores && cpuCores > 0 {
		// set only if we have non-default & non-empty value
		createBody.CPUCores = &cpuCores
	}

	if cpuSockets != dvCPUSockets && cpuSockets > 0 {
		// set only if we have non-default & non-empty value
		createBody.CPUSockets = &cpuSockets
	}

	if cpuUnits != dvCPUUnits && cpuUnits > 0 {
		// set only if we have non-default & non-empty value
		createBody.CPUUnits = &cpuUnits
	}

	if cpuNUMA != dvCPUNUMA && cpuNUMA {
		// set only if we have non-default & non-empty value
		createBody.NUMAEnabled = &cpuNUMA
	}

	if cpuLimit > 0 {
		createBody.CPULimit = &cpuLimit
	}

	if cpuAffinity != "" {
		createBody.CPUAffinity = &cpuAffinity
	}

	if cpuHotplugged > 0 {
		createBody.VirtualCPUCount = &cpuHotplugged
	}

	if cpuType != dvCPUType && cpuType != "" {
		// set only if we have non-default & non-empty value
		if createBody.CPUEmulation == nil {
			createBody.CPUEmulation = &vms.CustomCPUEmulation{}
		}

		createBody.CPUEmulation.Type = &cpuType
	}

	cpuFlagsConverted := make([]string, len(cpuFlags))
	for fi, flag := range cpuFlags {
		cpuFlagsConverted[fi] = flag.(string)
	}
	if len(cpuFlagsConverted) > 0 {
		// set only if we have non-default & non-empty value
		if createBody.CPUEmulation == nil {
			createBody.CPUEmulation = &vms.CustomCPUEmulation{}
		}

		createBody.CPUEmulation.Flags = &cpuFlagsConverted
	}

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if api.API().IsRootTicket() ||
		cpuArchitecture != dvCPUArchitecture {
		createBody.CPUArchitecture = &cpuArchitecture
	}

	return nil
}

func Read(vmConfig *vms.GetResponseData, api proxmox.Client, d *schema.ResourceData, clone bool) error {
	// Compare the CPU configuration to the one stored in the state.
	cpu := map[string]interface{}{}

	if vmConfig.CPUArchitecture != nil {
		cpu[mkCPUArchitecture] = *vmConfig.CPUArchitecture
	} else {
		// Default value of "arch" is "" according to the API documentation.
		// However, assume the provider's default value as a workaround when the root account is not being used.
		if !api.API().IsRootTicket() {
			cpu[mkCPUArchitecture] = dvCPUArchitecture
		} else {
			cpu[mkCPUArchitecture] = ""
		}
	}

	if vmConfig.CPUCores != nil {
		cpu[mkCPUCores] = *vmConfig.CPUCores
		//} else {
		//	// Default value of "cores" is "1" according to the API documentation.
		//	cpu[mkCPUCores] = 1
	}

	if vmConfig.VirtualCPUCount != nil {
		cpu[mkCPUHotplugged] = *vmConfig.VirtualCPUCount
	} else {
		// Default value of "vcpus" is "1" according to the API documentation.
		cpu[mkCPUHotplugged] = 0
	}

	if vmConfig.CPULimit != nil {
		cpu[mkCPULimit] = *vmConfig.CPULimit
	} else {
		// Default value of "cpulimit" is "0" according to the API documentation.
		cpu[mkCPULimit] = 0
	}

	if vmConfig.NUMAEnabled != nil {
		cpu[mkCPUNUMA] = *vmConfig.NUMAEnabled
	} else {
		// Default value of "numa" is "false" according to the API documentation.
		cpu[mkCPUNUMA] = false
	}

	if vmConfig.CPUSockets != nil {
		cpu[mkCPUSockets] = *vmConfig.CPUSockets
	} else {
		// Default value of "sockets" is "1" according to the API documentation.
		cpu[mkCPUSockets] = 1
	}

	if vmConfig.CPUEmulation != nil {
		if vmConfig.CPUEmulation.Flags != nil {
			convertedFlags := make([]interface{}, len(*vmConfig.CPUEmulation.Flags))

			for fi, fv := range *vmConfig.CPUEmulation.Flags {
				convertedFlags[fi] = fv
			}

			cpu[mkCPUFlags] = convertedFlags
			//} else {
			//	cpu[mkCPUFlags] = []interface{}{}
		}

		cpu[mkCPUType] = vmConfig.CPUEmulation.Type
		//} else {
		//	cpu[mkCPUFlags] = []interface{}{}
		//	// Default value of "cputype" is "qemu64" according to the QEMU documentation.
		//	cpu[mkCPUType] = dvCPUType
	}

	if vmConfig.CPUUnits != nil {
		cpu[mkCPUUnits] = *vmConfig.CPUUnits
	} else {
		// Default value of "cpuunits" is "1024" according to the API documentation.
		cpu[mkCPUUnits] = 1024
	}

	if vmConfig.CPUAffinity != nil {
		cpu[mkCPUAffinity] = *vmConfig.CPUAffinity
	} else {
		cpu[mkCPUAffinity] = ""
	}

	currentCPU := d.Get(MkCPU).([]interface{})

	//if len(clone) > 0 {
	//	if len(currentCPU) > 0 {
	//		err := d.Set(mkCPU, []interface{}{cpu})
	//		diags = append(diags, diag.FromErr(err)...)
	//	}
	//} else if len(currentCPU) > 0 ||
	//	cpu[mkCPUArchitecture] != dvCPUArchitecture ||
	//	cpu[mkCPUCores] != dvCPUCores ||
	//	len(cpu[mkCPUFlags].([]interface{})) > 0 ||
	//	cpu[mkCPUHotplugged] != dvCPUHotplugged ||
	//	cpu[mkCPULimit] != dvCPULimit ||
	//	cpu[mkCPUSockets] != dvCPUSockets ||
	//	cpu[mkCPUType] != dvCPUType ||
	//	cpu[mkCPUUnits] != dvCPUUnits {
	//	err := d.Set(mkCPU, []interface{}{cpu})
	//	diags = append(diags, diag.FromErr(err)...)
	//}

	if clone || len(currentCPU) > 0 ||
		cpu[mkCPUArchitecture] != dvCPUArchitecture ||
		//cpu[mkCPUCores] != dvCPUCores ||
		//len(cpu[mkCPUFlags].([]interface{})) > 0 ||
		cpu[mkCPUHotplugged] != dvCPUHotplugged ||
		cpu[mkCPULimit] != dvCPULimit ||
		cpu[mkCPUSockets] != dvCPUSockets ||
		//cpu[mkCPUType] != dvCPUType ||
		cpu[mkCPUUnits] != dvCPUUnits {
		err := d.Set(MkCPU, []interface{}{cpu})
		if err != nil {
			return fmt.Errorf("error setting CPU: %w", err)
		}
	}

	return nil
}

func Update(d *schema.ResourceData, resource *schema.Resource, api proxmox.Client, updateBody *vms.UpdateRequestBody) ([]string, bool, error) {
	// Prepare the new CPU configuration.

	if !d.HasChange(MkCPU) {
		return nil, false, nil
	}

	del := []string{}
	cpuBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{MkCPU},
		0,
		true,
	)
	if err != nil {
		return nil, false, err
	}

	cpuArchitecture := cpuBlock[mkCPUArchitecture].(string)
	cpuCores := cpuBlock[mkCPUCores].(int)
	cpuFlags := cpuBlock[mkCPUFlags].([]interface{})
	cpuHotplugged := cpuBlock[mkCPUHotplugged].(int)
	cpuLimit := cpuBlock[mkCPULimit].(int)
	cpuNUMA := types.CustomBool(cpuBlock[mkCPUNUMA].(bool))
	cpuSockets := cpuBlock[mkCPUSockets].(int)
	cpuType := cpuBlock[mkCPUType].(string)
	cpuUnits := cpuBlock[mkCPUUnits].(int)
	cpuAffinity := cpuBlock[mkCPUAffinity].(string)

	// Only the root account is allowed to change the CPU architecture, which makes this check necessary.
	if api.API().IsRootTicket() ||
		cpuArchitecture != dvCPUArchitecture {
		updateBody.CPUArchitecture = &cpuArchitecture
	}

	if cpuCores != dvCPUCores && cpuCores > 0 {
		// set only if we have non-default & non-empty value
		updateBody.CPUCores = &cpuCores
	}

	updateBody.CPUSockets = &cpuSockets
	updateBody.CPUUnits = &cpuUnits
	updateBody.NUMAEnabled = &cpuNUMA

	// CPU affinity is a special case, only root can change it.
	// we can't even have it in the delete list, as PVE will return an error for non-root.
	// Hence, checking explicitly if it has changed.
	if d.HasChange(MkCPU + ".0." + mkCPUAffinity) {
		if cpuAffinity != "" {
			updateBody.CPUAffinity = &cpuAffinity
		} else {
			del = append(del, "affinity")
		}
	}

	if cpuHotplugged > 0 {
		updateBody.VirtualCPUCount = &cpuHotplugged
	} else {
		del = append(del, "vcpus")
	}

	if cpuLimit > 0 {
		updateBody.CPULimit = &cpuLimit
	} else {
		del = append(del, "cpulimit")
	}

	cpuFlagsConverted := make([]string, len(cpuFlags))

	for fi, flag := range cpuFlags {
		cpuFlagsConverted[fi] = flag.(string)
	}

	if cpuType != dvCPUType && cpuType != "" {
		// set only if we have non-default & non-empty value
		if updateBody.CPUEmulation == nil {
			updateBody.CPUEmulation = &vms.CustomCPUEmulation{}
		}

		updateBody.CPUEmulation.Type = &cpuType
	}

	if len(cpuFlagsConverted) > 0 {
		// set only if we have non-default & non-empty value
		if updateBody.CPUEmulation == nil {
			updateBody.CPUEmulation = &vms.CustomCPUEmulation{}
		}

		updateBody.CPUEmulation.Flags = &cpuFlagsConverted
	}

	// TODO: this may not be true
	rebootRequired := true

	return del, rebootRequired, nil
}
