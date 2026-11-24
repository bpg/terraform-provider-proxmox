//go:build acceptance || all

//testacc:tier=light
//testacc:resource=firewall

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	clusterfirewall "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const sgResource = "proxmox_cluster_firewall_security_group.test"

func randomGroupName() string {
	return "tf" + strings.ToLower(gofakeit.LetterN(10))
}

func TestAccResourceClusterFirewallSecurityGroup(t *testing.T) {
	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps func(t *testing.T, name string) []resource.TestStep
	}{
		{"create, reorder, update rules and import", func(_ *testing.T, name string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_cluster_firewall_security_group" "test" {
							name    = %q
							comment = "Managed by Terraform"

							rule {
								type    = "in"
								action  = "ACCEPT"
								macro   = "SSH"
								log     = "info"
								comment = "Allow SSH"
							}

							rule {
								type   = "in"
								action = "ACCEPT"
								proto  = "tcp"
								dport  = "80"
								dest   = "192.168.1.5"
							}

							rule {
								type    = "in"
								action  = "DROP"
								enabled = false
							}
						}`, name)),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{
							"id":             name,
							"name":           name,
							"comment":        "Managed by Terraform",
							"rule.#":         "3",
							"rule.0.pos":     "0",
							"rule.0.macro":   "SSH",
							"rule.0.log":     "info",
							"rule.0.enabled": "true",
							"rule.1.pos":     "1",
							"rule.1.dport":   "80",
							"rule.1.dest":    "192.168.1.5",
							"rule.2.pos":     "2",
							"rule.2.action":  "DROP",
							"rule.2.enabled": "false",
						}),
						test.NoResourceAttributesSet(sgResource, []string{"rule.1.comment", "rule.1.log"}),
						checkGroupRules(te, name, []string{
							"in ACCEPT macro=SSH log=info comment=Allow SSH enable=1",
							"in ACCEPT dest=192.168.1.5 dport=80 proto=tcp enable=1",
							"in DROP enable=0",
						}),
					),
				},
				{
					// Moves DROP to the top, adds HTTPS, clears the SSH log and comments the HTTP rule.
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_cluster_firewall_security_group" "test" {
							name    = %q
							comment = "Updated by Terraform"

							rule {
								type    = "in"
								action  = "DROP"
								enabled = false
							}

							rule {
								type    = "in"
								action  = "ACCEPT"
								macro   = "SSH"
								comment = "Allow SSH"
							}

							rule {
								type   = "in"
								action = "ACCEPT"
								proto  = "tcp"
								dport  = "443"
							}

							rule {
								type    = "in"
								action  = "ACCEPT"
								proto   = "tcp"
								dport   = "80"
								dest    = "192.168.1.5"
								comment = "Allow HTTP"
							}
						}`, name)),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(sgResource, plancheck.ResourceActionUpdate),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{
							"comment":        "Updated by Terraform",
							"rule.#":         "4",
							"rule.0.action":  "DROP",
							"rule.1.macro":   "SSH",
							"rule.2.dport":   "443",
							"rule.3.comment": "Allow HTTP",
							"rule.3.pos":     "3",
						}),
						test.NoResourceAttributesSet(sgResource, []string{"rule.1.log"}),
						checkGroupRules(te, name, []string{
							"in DROP enable=0",
							"in ACCEPT macro=SSH comment=Allow SSH enable=1",
							"in ACCEPT dport=443 proto=tcp enable=1",
							"in ACCEPT dest=192.168.1.5 dport=80 proto=tcp comment=Allow HTTP enable=1",
						}),
					),
				},
				{
					ResourceName:      sgResource,
					ImportState:       true,
					ImportStateId:     name,
					ImportStateVerify: true,
				},
			}
		}},
		{"rename in place keeps the rules", func(_ *testing.T, name string) []resource.TestStep {
			renamed := name + "x"

			cfg := func(n string) string {
				return te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_cluster_firewall_security_group" "test" {
						name = %q

						rule {
							type   = "out"
							action = "ACCEPT"
							proto  = "udp"
							dport  = "53"
						}
					}`, n))
			}

			return []resource.TestStep{
				{Config: cfg(name)},
				{
					Config: cfg(renamed),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(sgResource, plancheck.ResourceActionUpdate),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{
							"id":     renamed,
							"name":   renamed,
							"rule.#": "1",
						}),
						checkGroupAbsent(te, name),
						checkGroupRules(te, renamed, []string{"out ACCEPT dport=53 proto=udp enable=1"}),
					),
				},
			}
		}},
		{"minimal group, then add and clear optional fields", func(_ *testing.T, name string) []resource.TestStep {
			minimal := te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_cluster_firewall_security_group" "test" {
					name = %q
				}`, name))

			return []resource.TestStep{
				{
					Config: minimal,
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{"name": name, "rule.#": "0"}),
						test.NoResourceAttributesSet(sgResource, []string{"comment"}),
						checkGroupRules(te, name, nil),
					),
				},
				{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_cluster_firewall_security_group" "test" {
							name    = %q
							comment = "temporary"

							rule {
								type   = "forward"
								action = "ACCEPT"
								source = "10.0.0.0/8"
								proto  = "tcp"
								sport  = "1024:65535"
							}
						}`, name)),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{
							"comment":       "temporary",
							"rule.#":        "1",
							"rule.0.type":   "forward",
							"rule.0.source": "10.0.0.0/8",
							"rule.0.sport":  "1024:65535",
						}),
						checkGroupRules(te, name, []string{"forward ACCEPT source=10.0.0.0/8 sport=1024:65535 proto=tcp enable=1"}),
					),
				},
				{
					Config: minimal,
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{"rule.#": "0"}),
						test.NoResourceAttributesSet(sgResource, []string{"comment"}),
						checkGroupComment(te, name, ""),
						checkGroupRules(te, name, nil),
					),
				},
			}
		}},
		{"rules added outside Terraform are removed", func(t *testing.T, name string) []resource.TestStep {
			cfg := te.RenderConfig(fmt.Sprintf(`
				resource "proxmox_cluster_firewall_security_group" "test" {
					name = %q

					rule {
						type   = "in"
						action = "ACCEPT"
						macro  = "Ping"
					}
				}`, name))

			return []resource.TestStep{
				{Config: cfg},
				{
					PreConfig: func() {
						createRuleOutOfBand(t, te, name)
					},
					Config: cfg,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(sgResource, plancheck.ResourceActionUpdate),
						},
					},
					Check: checkGroupRules(te, name, []string{"in ACCEPT macro=Ping enable=1"}),
				},
			}
		}},
		{"macro keeps the configured casing", func(_ *testing.T, name string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_cluster_firewall_security_group" "test" {
							name = %q

							rule {
								type   = "in"
								action = "ACCEPT"
								macro  = "ssh"
							}
						}`, name)),
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{"rule.0.macro": "ssh"}),
						checkGroupRules(te, name, []string{"in ACCEPT macro=SSH enable=1"}),
					),
				},
			}
		}},
		{"existing group is adopted and reconciled", func(t *testing.T, name string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						fw := te.ClusterClient().Firewall()
						err := fw.CreateGroup(context.Background(), &clusterfirewall.GroupCreateRequestBody{
							Group:   name,
							Comment: new("pre-existing"),
						})
						if err != nil {
							t.Fatalf("creating pre-existing group: %v", err)
						}

						createRuleOutOfBand(t, te, name)
					},
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_cluster_firewall_security_group" "test" {
							name = %q

							rule {
								type   = "in"
								action = "ACCEPT"
								macro  = "HTTPS"
							}
						}`, name)),
					Check: resource.ComposeTestCheckFunc(
						test.NoResourceAttributesSet(sgResource, []string{"comment"}),
						checkGroupComment(te, name, ""),
						checkGroupRules(te, name, []string{"in ACCEPT macro=HTTPS enable=1"}),
					),
				},
			}
		}},
		{"uppercase and invalid names are rejected", func(_ *testing.T, _ string) []resource.TestStep {
			cfg := func(n string) string {
				return te.RenderConfig(fmt.Sprintf(`
					resource "proxmox_cluster_firewall_security_group" "test" {
						name = %q
					}`, n))
			}

			return []resource.TestStep{
				{Config: cfg("WebServers"), PlanOnly: true, ExpectError: regexp.MustCompile(`must be lowercase`)},
				{Config: cfg("1web"), PlanOnly: true, ExpectError: regexp.MustCompile(`must start with a letter`)},
				{Config: cfg("a-name-way-too-long-for-pve"), PlanOnly: true, ExpectError: regexp.MustCompile(`length`)},
			}
		}},
		{"moved from the SDK resource without changes", func(_ *testing.T, name string) []resource.TestStep {
			rules := `
				rule {
					type    = "in"
					action  = "ACCEPT"
					proto   = "tcp"
					dport   = "22"
					comment = "Allow SSH"
				}

				rule {
					type   = "in"
					action = "DROP"
				}`

			return []resource.TestStep{
				{
					Config: te.RenderConfig(fmt.Sprintf(`
						resource "proxmox_virtual_environment_cluster_firewall_security_group" "test" {
							name = %q
							%s
						}`, name, rules)),
				},
				{
					Config: te.RenderConfig(fmt.Sprintf(`
						moved {
							from = proxmox_virtual_environment_cluster_firewall_security_group.test
							to   = proxmox_cluster_firewall_security_group.test
						}

						resource "proxmox_cluster_firewall_security_group" "test" {
							name = %q
							%s
						}`, name, rules)),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(sgResource, plancheck.ResourceActionNoop),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						test.ResourceAttributes(sgResource, map[string]string{
							"id":             name,
							"rule.#":         "2",
							"rule.0.comment": "Allow SSH",
							"rule.1.action":  "DROP",
						}),
						test.NoResourceAttributesSet(sgResource, []string{"comment", "rule.1.comment", "rule.0.dest"}),
					),
				},
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := randomGroupName()

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				CheckDestroy:             checkGroupAbsentAfterDestroy(te, name),
				Steps:                    tt.steps(t, name),
			})
		})
	}
}

// describeRule renders a rule as PVE stores it, listing only the fields that are set.
func describeRule(r *firewall.RuleGetResponseData) string {
	parts := []string{r.Type, r.Action}

	for _, f := range []struct {
		key   string
		value *string
	}{
		{"macro", r.Macro},
		{"source", r.Source},
		{"sport", r.SPort},
		{"dest", r.Dest},
		{"dport", r.DPort},
		{"proto", r.Proto},
		{"log", r.Log},
		{"comment", r.Comment},
	} {
		if f.value != nil && *f.value != "" {
			parts = append(parts, f.key+"="+*f.value)
		}
	}

	enable := 0
	if r.Enable != nil && bool(*r.Enable) {
		enable = 1
	}

	return strings.Join(append(parts, fmt.Sprintf("enable=%d", enable)), " ")
}

func checkGroupRules(te *test.Environment, group string, want []string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		ctx := context.Background()
		api := te.ClusterClient().Firewall().SecurityGroup(group)

		list, err := api.ListRules(ctx)
		if err != nil {
			return fmt.Errorf("listing rules of %q: %w", group, err)
		}

		slices.SortFunc(list, func(a, b *firewall.RuleListResponseData) int { return a.Pos - b.Pos })

		got := make([]string, 0, len(list))

		for _, item := range list {
			r, err := api.GetRule(ctx, item.Pos)
			if err != nil {
				return fmt.Errorf("reading rule %d of %q: %w", item.Pos, group, err)
			}

			got = append(got, describeRule(r))
		}

		if !slices.Equal(got, want) {
			return fmt.Errorf("rules of %q:\n got: %q\nwant: %q", group, got, want)
		}

		return nil
	}
}

func findGroup(te *test.Environment, group string) (*clusterfirewall.GroupListResponseData, error) {
	groups, err := te.ClusterClient().Firewall().ListGroups(context.Background())
	if err != nil {
		return nil, fmt.Errorf("listing security groups: %w", err)
	}

	for _, g := range groups {
		if g.Group == group {
			return g, nil
		}
	}

	return nil, nil
}

func checkGroupComment(te *test.Environment, group, want string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		g, err := findGroup(te, group)
		if err != nil {
			return err
		}

		if g == nil {
			return fmt.Errorf("security group %q not found", group)
		}

		if got := ptrValue(g.Comment); got != want {
			return fmt.Errorf("comment of %q: got %q, want %q", group, got, want)
		}

		return nil
	}
}

func checkGroupAbsent(te *test.Environment, group string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		g, err := findGroup(te, group)
		if err != nil {
			return err
		}

		if g != nil {
			return fmt.Errorf("security group %q still exists", group)
		}

		return nil
	}
}

// checkGroupAbsentAfterDestroy also covers the renamed group used by the rename scenario.
func checkGroupAbsentAfterDestroy(te *test.Environment, group string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(checkGroupAbsent(te, group), checkGroupAbsent(te, group+"x"))
}

func createRuleOutOfBand(t *testing.T, te *test.Environment, group string) {
	t.Helper()

	enable := proxmoxtypes.CustomBool(true)

	err := te.ClusterClient().Firewall().SecurityGroup(group).CreateRule(context.Background(), &firewall.RuleCreateRequestBody{
		BaseRule: firewall.BaseRule{Enable: &enable, Proto: new("tcp"), DPort: new("8080")},
		Action:   "ACCEPT",
		Type:     "in",
	})
	if err != nil {
		t.Fatalf("creating out-of-band rule in %q: %v", group, err)
	}
}

func ptrValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
