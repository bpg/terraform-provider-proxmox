package vm

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createAgent(d *schema.ResourceData, updateBody *vms.UpdateRequestBody) {
	agent := d.Get(mkAgent).([]interface{})
	if len(agent) > 0 {
		agentBlock := agent[0].(map[string]interface{})

		agentEnabled := types.CustomBool(
			agentBlock[mkAgentEnabled].(bool),
		)
		agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
		agentType := agentBlock[mkAgentType].(string)

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
		[]string{mkAgent},
		0,
		true,
	)
	if err != nil {
		return nil, err
	}

	agentEnabled := types.CustomBool(
		agentBlock[mkAgentEnabled].(bool),
	)
	agentTrim := types.CustomBool(agentBlock[mkAgentTrim].(bool))
	agentType := agentBlock[mkAgentType].(string)

	return &vms.CustomAgent{
		Enabled:         &agentEnabled,
		TrimClonedDisks: &agentTrim,
		Type:            &agentType,
	}, nil
}

func setAgent(d *schema.ResourceData, clone bool, vmConfig *vms.GetResponseData) error {
	// Compare the agent configuration to the one stored in the state.
	currentAgent := d.Get(mkAgent).([]interface{})

	if !clone || len(currentAgent) > 0 {
		if vmConfig.Agent != nil {
			agent := map[string]interface{}{}

			if vmConfig.Agent.Enabled != nil {
				agent[mkAgentEnabled] = bool(*vmConfig.Agent.Enabled)
			} else {
				agent[mkAgentEnabled] = false
			}

			if vmConfig.Agent.TrimClonedDisks != nil {
				agent[mkAgentTrim] = bool(
					*vmConfig.Agent.TrimClonedDisks,
				)
			} else {
				agent[mkAgentTrim] = false
			}

			if len(currentAgent) > 0 {
				currentAgentBlock := currentAgent[0].(map[string]interface{})
				currentAgentTimeout := currentAgentBlock[mkAgentTimeout].(string)

				if currentAgentTimeout != "" {
					agent[mkAgentTimeout] = currentAgentTimeout
				} else {
					agent[mkAgentTimeout] = dvAgentTimeout
				}
			} else {
				agent[mkAgentTimeout] = dvAgentTimeout
			}

			if vmConfig.Agent.Type != nil {
				agent[mkAgentType] = *vmConfig.Agent.Type
			} else {
				agent[mkAgentType] = ""
			}

			if clone {
				if len(currentAgent) > 0 {
					return d.Set(mkAgent, []interface{}{agent})
				}
			} else if len(currentAgent) > 0 ||
				agent[mkAgentEnabled] != dvAgentEnabled ||
				agent[mkAgentTimeout] != dvAgentTimeout ||
				agent[mkAgentTrim] != dvAgentTrim ||
				agent[mkAgentType] != dvAgentType {
				return d.Set(mkAgent, []interface{}{agent})
			}
		} else if clone {
			if len(currentAgent) > 0 {
				return d.Set(mkAgent, []interface{}{})
			}
		} else {
			return d.Set(mkAgent, []interface{}{})
		}
	}

	return nil
}
