package vm

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createAgent(d *schema.ResourceData, updateBody *vms.UpdateRequestBody) {
	agent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})
	if len(agent) > 0 {
		agentBlock := agent[0].(map[string]interface{})

		agentEnabled := types.CustomBool(
			agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
		agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

		updateBody.Agent = &vms.CustomAgent{
			Enabled:         &agentEnabled,
			TrimClonedDisks: &agentTrim,
			Type:            &agentType,
		}
	}
}

func customAgent(d *schema.ResourceData, resource *schema.Resource) (*vms.CustomAgent, error) {
	agentBlock, err := structure.GetSchemaBlock(
		resource,
		d,
		[]string{mkResourceVirtualEnvironmentVMAgent},
		0,
		true,
	)
	if err != nil {
		return nil, err
	}

	agentEnabled := types.CustomBool(
		agentBlock[mkResourceVirtualEnvironmentVMAgentEnabled].(bool),
	)
	agentTrim := types.CustomBool(agentBlock[mkResourceVirtualEnvironmentVMAgentTrim].(bool))
	agentType := agentBlock[mkResourceVirtualEnvironmentVMAgentType].(string)

	return &vms.CustomAgent{
		Enabled:         &agentEnabled,
		TrimClonedDisks: &agentTrim,
		Type:            &agentType,
	}, nil
}

func setAgent(d *schema.ResourceData, clone bool, vmConfig *vms.GetResponseData) error {
	// Compare the agent configuration to the one stored in the state.
	currentAgent := d.Get(mkResourceVirtualEnvironmentVMAgent).([]interface{})

	if !clone || len(currentAgent) > 0 {
		if vmConfig.Agent != nil {
			agent := map[string]interface{}{}

			if vmConfig.Agent.Enabled != nil {
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] = bool(*vmConfig.Agent.Enabled)
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] = false
			}

			if vmConfig.Agent.TrimClonedDisks != nil {
				agent[mkResourceVirtualEnvironmentVMAgentTrim] = bool(
					*vmConfig.Agent.TrimClonedDisks,
				)
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentTrim] = false
			}

			if len(currentAgent) > 0 {
				currentAgentBlock := currentAgent[0].(map[string]interface{})
				currentAgentTimeout := currentAgentBlock[mkResourceVirtualEnvironmentVMAgentTimeout].(string)

				if currentAgentTimeout != "" {
					agent[mkResourceVirtualEnvironmentVMAgentTimeout] = currentAgentTimeout
				} else {
					agent[mkResourceVirtualEnvironmentVMAgentTimeout] = dvResourceVirtualEnvironmentVMAgentTimeout
				}
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentTimeout] = dvResourceVirtualEnvironmentVMAgentTimeout
			}

			if vmConfig.Agent.Type != nil {
				agent[mkResourceVirtualEnvironmentVMAgentType] = *vmConfig.Agent.Type
			} else {
				agent[mkResourceVirtualEnvironmentVMAgentType] = ""
			}

			if clone {
				if len(currentAgent) > 0 {
					return d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})
				}
			} else if len(currentAgent) > 0 ||
				agent[mkResourceVirtualEnvironmentVMAgentEnabled] != dvResourceVirtualEnvironmentVMAgentEnabled ||
				agent[mkResourceVirtualEnvironmentVMAgentTimeout] != dvResourceVirtualEnvironmentVMAgentTimeout ||
				agent[mkResourceVirtualEnvironmentVMAgentTrim] != dvResourceVirtualEnvironmentVMAgentTrim ||
				agent[mkResourceVirtualEnvironmentVMAgentType] != dvResourceVirtualEnvironmentVMAgentType {
				return d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{agent})

			}
		} else if clone {
			if len(currentAgent) > 0 {
				return d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
			}
		} else {
			return d.Set(mkResourceVirtualEnvironmentVMAgent, []interface{}{})
		}
	}
	return nil
}
