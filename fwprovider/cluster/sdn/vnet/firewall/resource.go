/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	proxmoxfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithConfigure      = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
)

// Resource manages the complete ordered firewall ruleset for an SDN VNet.
type Resource struct {
	client *cluster.Client
}

// NewResource creates an SDN VNet firewall rules resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata defines the resource type name.
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_vnet_firewall_rules"
}

// Configure adds the provider-configured client to the resource.
func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client.Cluster()
}

// Schema defines the Terraform schema for SDN VNet firewall rules.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the complete ordered firewall ruleset for a Proxmox VE SDN VNet.",
		MarkdownDescription: "Manages the complete ordered firewall ruleset for a Proxmox VE SDN VNet. " +
			"The resource has exclusive ownership of the VNet's rules; import existing rules before managing them.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"vnet": schema.StringAttribute{
				Description: "The identifier of the SDN VNet whose firewall rules are managed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: validators.SDNID(),
			},
			"rules": schema.ListNestedAttribute{
				Description: "The ordered firewall rules. List order is the rule evaluation order.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description: "The action for a forwarded packet: `ACCEPT` or `DROP`. Exactly one of `action` or `security_group` must be set.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("ACCEPT", "DROP"),
							},
						},
						"comment": schema.StringAttribute{
							Description: "A descriptive comment for the rule.",
							Optional:    true,
						},
						"dest": schema.StringAttribute{
							Description: "A destination address, IP set, alias, range, or comma-separated address list.",
							Optional:    true,
						},
						"dport": schema.StringAttribute{
							Description: "A destination port, service, range, or comma-separated port list.",
							Optional:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the rule is enabled. Defaults to `true`.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
						},
						"icmp_type": schema.StringAttribute{
							Description: "The ICMP type. Valid only when `proto` is `icmp`, `icmpv6`, or `ipv6-icmp`.",
							Optional:    true,
						},
						"log": schema.StringAttribute{
							Description: "The log level for the rule.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"),
							},
						},
						"macro": schema.StringAttribute{
							Description: "A predefined Proxmox firewall macro.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(128),
							},
						},
						"proto": schema.StringAttribute{
							Description: "An IP protocol name or number.",
							Optional:    true,
						},
						"security_group": schema.StringAttribute{
							Description: "A cluster firewall security group to insert. Exactly one of `action` or `security_group` must be set.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]+$`),
									"must start with a letter and contain only letters, digits, hyphens, and underscores",
								),
								stringvalidator.LengthBetween(2, 18),
							},
						},
						"source": schema.StringAttribute{
							Description: "A source address, IP set, alias, range, or comma-separated address list.",
							Optional:    true,
						},
						"sport": schema.StringAttribute{
							Description: "A source port, service, range, or comma-separated port list.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// ValidateConfig validates rule modes that depend on multiple nested attributes.
func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg model

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	for i, rule := range cfg.Rules {
		actionKnown := !rule.Action.IsUnknown()

		groupKnown := !rule.SecurityGroup.IsUnknown()
		if !actionKnown || !groupKnown {
			continue
		}

		actionSet := !rule.Action.IsNull()
		groupSet := !rule.SecurityGroup.IsNull()
		rulePath := path.Root("rules").AtListIndex(i)

		if actionSet == groupSet {
			resp.Diagnostics.AddAttributeError(
				rulePath,
				"Invalid attribute combination",
				"Exactly one of `action` or `security_group` must be set for each rule.",
			)

			continue
		}

		if groupSet && rule.hasPacketAttributes() {
			resp.Diagnostics.AddAttributeError(
				rulePath.AtName("security_group"),
				"Invalid attribute combination",
				"Packet match attributes cannot be used with security_group rules; only comment and enabled are supported.",
			)
		}
	}
}

// Create creates the managed VNet firewall ruleset without overwriting unowned rules.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet := plan.VNet.ValueString()

	vnetClient := r.client.SDNVnets(vnet)
	if _, err := vnetClient.GetVnet(ctx); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create SDN VNet Firewall Rules %q", vnet), err.Error())
		return
	}

	rulesAPI := vnetClient.Firewall()

	existing, err := rulesAPI.ListRules(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create SDN VNet Firewall Rules %q", vnet), err.Error())
		return
	}

	if len(existing) > 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Create SDN VNet Firewall Rules %q", vnet),
			fmt.Sprintf(
				"The VNet already has %d firewall rule(s). Import the ruleset with "+
					"`terraform import proxmox_sdn_vnet_firewall_rules.<name> %s` before managing it.",
				len(existing),
				vnet,
			),
		)

		return
	}

	if err := createRules(ctx, rulesAPI, plan.Rules); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create SDN VNet Firewall Rules %q", vnet), err.Error())
		return
	}

	readModel, err := readRules(ctx, vnet, rulesAPI)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read SDN VNet Firewall Rules %q After Creation", vnet), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// Read refreshes the ordered VNet firewall rules from Proxmox.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet := state.VNet.ValueString()

	vnetClient := r.client.SDNVnets(vnet)
	if _, err := vnetClient.GetVnet(ctx); err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read SDN VNet Firewall Rules %q", vnet), err.Error())

		return
	}

	readModel, err := readRules(ctx, vnet, vnetClient.Firewall())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read SDN VNet Firewall Rules %q", vnet), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// Update reconciles rule identities, positions, mutable fields, and removals.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model
	var state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet := plan.VNet.ValueString()

	rulesAPI := r.client.SDNVnets(vnet).Firewall()
	if err := reconcileRules(ctx, rulesAPI, state.Rules, plan.Rules); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Update SDN VNet Firewall Rules %q", vnet), err.Error())
		return
	}

	readModel, err := readRules(ctx, vnet, rulesAPI)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read SDN VNet Firewall Rules %q After Update", vnet), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

// Delete removes all firewall rules owned by the VNet ruleset resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnet := state.VNet.ValueString()

	vnetClient := r.client.SDNVnets(vnet)
	if _, err := vnetClient.GetVnet(ctx); err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete SDN VNet Firewall Rules %q", vnet), err.Error())

		return
	}

	if err := deleteAllRules(ctx, vnetClient.Firewall()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete SDN VNet Firewall Rules %q", vnet), err.Error())
	}
}

// ImportState imports the complete firewall ruleset using its VNet identifier.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vnetClient := r.client.SDNVnets(req.ID)
	if _, err := vnetClient.GetVnet(ctx); err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("SDN VNet Not Found", fmt.Sprintf("SDN VNet with ID %q was not found", req.ID))
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN VNet Firewall Rules %q", req.ID), err.Error())

		return
	}

	readModel, err := readRules(ctx, req.ID, vnetClient.Firewall())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN VNet Firewall Rules %q", req.ID), err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (m *ruleModel) hasPacketAttributes() bool {
	return attribute.IsDefined(m.Dest) ||
		attribute.IsDefined(m.DPort) ||
		attribute.IsDefined(m.ICMPType) ||
		attribute.IsDefined(m.Log) ||
		attribute.IsDefined(m.Macro) ||
		attribute.IsDefined(m.Proto) ||
		attribute.IsDefined(m.Source) ||
		attribute.IsDefined(m.SPort)
}

func createRules(ctx context.Context, rulesAPI proxmoxfirewall.Rule, rules []ruleModel) error {
	for i, v := range slices.Backward(rules) {
		if err := rulesAPI.CreateRule(ctx, v.toCreateAPI()); err != nil {
			return fmt.Errorf("error creating rule at desired position %d: %w", i, err)
		}
	}

	return nil
}

func readRules(ctx context.Context, vnet string, rulesAPI proxmoxfirewall.Rule) (*model, error) {
	rulePositions, err := rulesAPI.ListRules(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(rulePositions, func(i, j int) bool {
		return rulePositions[i].Pos < rulePositions[j].Pos
	})

	rules := make([]*proxmoxfirewall.RuleGetResponseData, 0, len(rulePositions))
	for _, position := range rulePositions {
		rule, err := rulesAPI.GetRule(ctx, position.Pos)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	result := &model{}
	result.fromAPI(vnet, rules)

	return result, nil
}

func reconcileRules(ctx context.Context, rulesAPI proxmoxfirewall.Rule, oldRules, newRules []ruleModel) error {
	oldBySignature := make(map[string][]ruleModel)

	for _, rule := range oldRules {
		signature := rule.signature()
		oldBySignature[signature] = append(oldBySignature[signature], rule)
	}

	consumed := make(map[string]int)
	matched := make(map[int]ruleModel)

	for i, rule := range newRules {
		signature := rule.signature()
		oldIndex := consumed[signature]
		consumed[signature]++

		if oldIndex < len(oldBySignature[signature]) {
			matched[i] = oldBySignature[signature][oldIndex]
			continue
		}

		if err := rulesAPI.CreateRule(ctx, rule.toCreateAPI()); err != nil {
			return fmt.Errorf("error creating rule for desired position %d: %w", i, err)
		}
	}

	for targetPosition, rule := range newRules {
		currentPositions, err := rulesAPI.ListRules(ctx)
		if err != nil {
			return err
		}

		sort.Slice(currentPositions, func(i, j int) bool {
			return currentPositions[i].Pos < currentPositions[j].Pos
		})

		currentPosition := -1

		for _, position := range currentPositions {
			if position.Pos < targetPosition {
				continue
			}

			currentRule, err := rulesAPI.GetRule(ctx, position.Pos)
			if err != nil {
				if errors.Is(err, proxmoxfirewall.ErrNoRuleAtPosition) {
					continue
				}

				return err
			}

			if apiRuleSignature(currentRule) == rule.signature() {
				currentPosition = position.Pos
				break
			}
		}

		if currentPosition == -1 {
			return fmt.Errorf("unable to find the rule for desired position %d during reconciliation", targetPosition)
		}

		if currentPosition != targetPosition {
			if err := rulesAPI.UpdateRule(ctx, currentPosition, &proxmoxfirewall.RuleUpdateRequestBody{
				MoveTo: &targetPosition,
			}); err != nil {
				return fmt.Errorf("error moving rule from position %d to %d: %w", currentPosition, targetPosition, err)
			}
		}
	}

	for position, previous := range matched {
		if err := rulesAPI.UpdateRule(ctx, position, newRules[position].toUpdateAPI(previous)); err != nil {
			return fmt.Errorf("error updating rule at position %d: %w", position, err)
		}
	}

	currentPositions, err := rulesAPI.ListRules(ctx)
	if err != nil {
		return err
	}

	sort.Slice(currentPositions, func(i, j int) bool {
		return currentPositions[i].Pos < currentPositions[j].Pos
	})

	for i := len(currentPositions) - 1; i >= len(newRules); i-- {
		position := currentPositions[i].Pos
		if err := rulesAPI.DeleteRule(ctx, position); err != nil {
			return fmt.Errorf("error deleting rule at position %d: %w", position, err)
		}
	}

	return nil
}

func deleteAllRules(ctx context.Context, rulesAPI proxmoxfirewall.Rule) error {
	rules, err := rulesAPI.ListRules(ctx)
	if err != nil {
		return err
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Pos > rules[j].Pos
	})

	for _, rule := range rules {
		if err := rulesAPI.DeleteRule(ctx, rule.Pos); err != nil {
			return fmt.Errorf("error deleting rule at position %d: %w", rule.Pos, err)
		}
	}

	return nil
}
