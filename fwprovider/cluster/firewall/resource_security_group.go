/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

const sdkSecurityGroupTypeName = "proxmox_virtual_environment_cluster_firewall_security_group"

var (
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
	_ resource.ResourceWithModifyPlan  = &securityGroupResource{}
	_ resource.ResourceWithMoveState   = &securityGroupResource{}
)

type securityGroupResource struct {
	client clusterfirewall.API
}

// NewSecurityGroupResource creates the proxmox_cluster_firewall_security_group resource.
func NewSecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

func (r *securityGroupResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_cluster_firewall_security_group"
}

func (r *securityGroupResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

	r.client = cfg.Client.Cluster().Firewall()
}

func (r *securityGroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	optionalString := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			Description: description,
			Optional:    true,
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		}
	}

	resp.Schema = schema.Schema{
		Description: "Manages a cluster-level firewall security group and its rules. A security group is a " +
			"collection of rules that can be referenced from cluster, node, VM and container rules.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The security group name.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Security group name. Must be lowercase, because PVE lowercases security group names. " +
					"Changing it renames the group in place, and PVE updates the rules that reference it.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 18),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-_]+$`),
						"must start with a letter and contain only letters, digits, '-' and '_'",
					),
					lowercaseValidator{},
				},
			},
			"comment": optionalString("Security group comment."),
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "Firewall rule, in evaluation order (multiple blocks supported).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"pos": schema.Int64Attribute{
							Description: "Position of the rule in the security group.",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "Rule action (`ACCEPT`, `DROP`, `REJECT`).",
							Required:    true,
							Validators:  []validator.String{stringvalidator.OneOf("ACCEPT", "DROP", "REJECT")},
						},
						"type": schema.StringAttribute{
							Description: "Rule type (`in`, `out`, `forward`).",
							Required:    true,
							Validators:  []validator.String{stringvalidator.OneOf("in", "out", "forward")},
						},
						"comment": optionalString("Rule comment."),
						"dest": optionalString("Restrict packet destination address. This can refer to a single IP " +
							"address, an IP set ('+ipsetname') or an IP alias definition. You can also specify an address " +
							"range like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
							"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists."),
						"dport": optionalString("Restrict TCP/UDP destination port. You can use service names or simple " +
							"numbers (0-65535), as defined in '/etc/services'. Port ranges can be specified with " +
							"'\\d+:\\d+', for example '80:85', and you can use comma separated list to match several " +
							"ports or ranges."),
						"enabled": schema.BoolAttribute{
							Description: "Enable this rule. Defaults to `true`.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
						},
						"log": schema.StringAttribute{
							Description: "Log level for this rule (`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, " +
								"`info`, `debug`, `nolog`).",
							Optional: true,
							Validators: []validator.String{stringvalidator.OneOf(
								"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog",
							)},
						},
						"macro": optionalString("Use predefined standard macro."),
						"proto": optionalString("Restrict packet protocol. You can use protocol names or simple " +
							"numbers (0-255), as defined in '/etc/protocols'."),
						"source": optionalString("Restrict packet source address. This can refer to a single IP " +
							"address, an IP set ('+ipsetname') or an IP alias definition. You can also specify an address " +
							"range like '20.34.101.207-201.3.9.99', or a list of IP addresses and networks (entries are " +
							"separated by comma). Please do not mix IPv4 and IPv6 addresses inside such lists."),
						"sport": optionalString("Restrict TCP/UDP source port. You can use service names or simple " +
							"numbers (0-65535), as defined in '/etc/services'. Port ranges can be specified with " +
							"'\\d+:\\d+', for example '80:85', and you can use comma separated list to match several " +
							"ports or ranges."),
					},
				},
			},
		},
	}
}

// ModifyPlan keeps id in step with name, so a rename shows the new id in the plan instead of "known after apply".
func (r *securityGroupResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var name types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)

	if resp.Diagnostics.HasError() || name.IsUnknown() {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), name)...)
}

func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	err := r.client.CreateGroup(ctx, &clusterfirewall.GroupCreateRequestBody{
		Group:   name,
		Comment: attribute.StringPtrFromValue(plan.Comment),
	})
	if err != nil {
		if !errors.Is(err, api.ErrResourceAlreadyExists) {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create Security Group %q", name), err.Error())
			return
		}

		// Adopt the existing group, as the SDK resource did, and reconcile it to the plan below.
		err = r.updateGroup(ctx, name, name, plan.Comment)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Update Existing Security Group %q", name), err.Error())
			return
		}
	}

	// Track the group before touching rules, so a failed rule call leaves it tainted instead of orphaned.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)

	r.applyAndRead(ctx, name, &plan, &resp.State, &resp.Diagnostics)
}

func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.ID.ValueString()
	prior := state.ruleModels(ctx, &resp.Diagnostics)

	found := r.read(ctx, name, &state, prior, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state securityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	oldName := state.ID.ValueString()
	newName := plan.Name.ValueString()

	if oldName != newName || !plan.Comment.Equal(state.Comment) {
		err := r.updateGroup(ctx, oldName, newName, plan.Comment)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Update Security Group %q", oldName), err.Error())
			return
		}
	}

	r.applyAndRead(ctx, newName, &plan, &resp.State, &resp.Diagnostics)
}

func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.ID.ValueString()
	rules := r.client.SecurityGroup(name)

	// PVE refuses to delete a group that still has rules, including ones added outside Terraform.
	current, err := readRules(ctx, rules)
	if err != nil {
		if isGroupNotFound(err) {
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read Rules of Security Group %q", name), err.Error())

		return
	}

	err = applyRuleOps(ctx, rules, planRuleOps(current, nil))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete Rules of Security Group %q", name), err.Error())
		return
	}

	err = r.client.DeleteGroup(ctx, name)
	if err != nil && !isGroupNotFound(err) {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete Security Group %q", name), err.Error())
	}
}

func (r *securityGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *securityGroupResource) MoveState(_ context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: moveSecurityGroupFromSDK}}
}

// updateGroup renames and/or re-comments a group. PVE's update is a create with rename set to the current name; an
// empty comment clears it.
func (r *securityGroupResource) updateGroup(ctx context.Context, oldName, newName string, comment types.String) error {
	c := comment.ValueString()

	return r.client.UpdateGroup(ctx, &clusterfirewall.GroupUpdateRequestBody{
		Group:   newName,
		ReName:  &oldName,
		Comment: &c,
	})
}

// applyAndRead reconciles the group's rules with the plan and stores what PVE reports back.
func (r *securityGroupResource) applyAndRead(
	ctx context.Context,
	name string,
	plan *securityGroupModel,
	state stateSetter,
	diags *diag.Diagnostics,
) {
	planned := plan.ruleModels(ctx, diags)
	if diags.HasError() {
		return
	}

	desired := make([]rule, len(planned))
	for i, m := range planned {
		desired[i] = m.toRule()
	}

	rules := r.client.SecurityGroup(name)

	current, err := readRules(ctx, rules)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read Rules of Security Group %q", name), err.Error())
		return
	}

	err = applyRuleOps(ctx, rules, planRuleOps(current, desired))
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Update Rules of Security Group %q", name), err.Error())
		return
	}

	var result securityGroupModel

	found := r.read(ctx, name, &result, planned, diags)
	if diags.HasError() {
		return
	}

	if !found {
		diags.AddError(
			fmt.Sprintf("Unable to Read Security Group %q", name),
			"The security group was not found after it was created or updated.",
		)

		return
	}

	diags.Append(state.Set(ctx, &result)...)
}

type stateSetter interface {
	Set(ctx context.Context, val any) diag.Diagnostics
}

// read fills m from PVE and reports whether the group exists.
func (r *securityGroupResource) read(
	ctx context.Context,
	name string,
	m *securityGroupModel,
	prior []ruleModel,
	diags *diag.Diagnostics,
) bool {
	groups, err := r.client.ListGroups(ctx)
	if err != nil {
		diags.AddError("Unable to List Security Groups", err.Error())
		return false
	}

	idx := slices.IndexFunc(groups, func(g *clusterfirewall.GroupListResponseData) bool { return g.Group == name })
	if idx < 0 {
		return false
	}

	rules, err := readRules(ctx, r.client.SecurityGroup(name))
	if err != nil {
		if isGroupNotFound(err) {
			return false
		}

		diags.AddError(fmt.Sprintf("Unable to Read Rules of Security Group %q", name), err.Error())

		return false
	}

	m.fromAPI(ctx, groups[idx], rules, prior, diags)

	return true
}

func readRules(ctx context.Context, api firewall.Rule) ([]rule, error) {
	list, err := api.ListRules(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(list, func(a, b *firewall.RuleListResponseData) int { return a.Pos - b.Pos })

	rules := make([]rule, 0, len(list))

	for _, item := range list {
		data, err := api.GetRule(ctx, item.Pos)
		if errors.Is(err, firewall.ErrNoRuleAtPosition) {
			continue
		}

		if err != nil {
			return nil, err
		}

		rules = append(rules, ruleFromAPI(data))
	}

	return rules, nil
}

func applyRuleOps(ctx context.Context, api firewall.Rule, ops []ruleOp) error {
	for _, op := range ops {
		var err error

		switch op.kind {
		case ruleOpCreate:
			err = api.CreateRule(ctx, op.rule.createBody())
		case ruleOpMove:
			err = api.UpdateRule(ctx, op.pos, &firewall.RuleUpdateRequestBody{MoveTo: &op.moveTo})
		case ruleOpUpdate:
			err = api.UpdateRule(ctx, op.pos, op.rule.updateBody(op.current))
		case ruleOpDelete:
			err = api.DeleteRule(ctx, op.pos)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// isGroupNotFound matches PVE's "no such security group '<name>'", which the API client doesn't classify as
// ErrResourceDoesNotExist.
func isGroupNotFound(err error) bool {
	return errors.Is(err, api.ErrResourceDoesNotExist) || strings.Contains(err.Error(), "no such security group")
}

// sdkSecurityGroupState is the state JSON of proxmox_virtual_environment_cluster_firewall_security_group. Its rule
// security_group and iface fields, and the node_name/vm_id/container_id selectors, are dropped: PVE rejects group and
// iface in security group rules, and the selectors never applied to a cluster security group.
type sdkSecurityGroupState struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Rules   []struct {
		Action  string `json:"action"`
		Type    string `json:"type"`
		Comment string `json:"comment"`
		Dest    string `json:"dest"`
		DPort   string `json:"dport"`
		Enabled bool   `json:"enabled"`
		Log     string `json:"log"`
		Macro   string `json:"macro"`
		Proto   string `json:"proto"`
		Source  string `json:"source"`
		SPort   string `json:"sport"`
	} `json:"rule"`
}

// moveSecurityGroupFromSDK supports `moved` blocks from the SDK resource. The SDK stored unset strings as "", which
// become null here.
func moveSecurityGroupFromSDK(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != sdkSecurityGroupTypeName {
		return
	}

	if req.SourceRawState == nil {
		resp.Diagnostics.AddError("Unable to Move Security Group State", "The source state is empty.")
		return
	}

	var src sdkSecurityGroupState

	err := json.Unmarshal(req.SourceRawState.JSON, &src)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Move Security Group State", err.Error())
		return
	}

	rules := make([]rule, len(src.Rules))
	for i, sr := range src.Rules {
		rules[i] = rule(sr)
	}

	var m securityGroupModel

	m.fromAPI(ctx, &clusterfirewall.GroupListResponseData{Group: src.ID, Comment: &src.Comment}, rules, nil, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.TargetState.Set(ctx, &m)...)
}

// lowercaseValidator rejects names with uppercase letters: PVE lowercases security group names on creation, so a
// mixed-case name would make the follow-up rule calls fail with "no such security group".
type lowercaseValidator struct{}

func (v lowercaseValidator) Description(_ context.Context) string {
	return "value must be lowercase"
}

func (v lowercaseValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v lowercaseValidator) ValidateString(
	_ context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if lower := strings.ToLower(value); value != lower {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Security Group Name",
			fmt.Sprintf("%q must be lowercase: PVE lowercases security group names, use %q instead", value, lower),
		)
	}
}
