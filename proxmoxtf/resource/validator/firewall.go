/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	rateExpression  = regexp.MustCompile(`[1-9][0-9]*/(second|minute|hour|day)`)
	ifaceExpression = regexp.MustCompile(`net\d+`)
)

// FirewallRate returns a schema validation function for a firewall rate.
func FirewallRate() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringMatch(
		rateExpression,
		"Must be a valid rate expression, e.g. '1/second'",
	))
}

// FirewallIFace returns a schema validation function for a firewall iface.
func FirewallIFace() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringMatch(
		ifaceExpression,
		"Must be a valid VM/Container iface key, e.g. 'net0'",
	))
}

// FirewallPolicy returns a schema validation function for a firewall policy.
func FirewallPolicy() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(
		[]string{"ACCEPT", "REJECT", "DROP"},
		false,
	))
}

// FirewallLogLevel returns a schema validation function for a firewall log level.
func FirewallLogLevel() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(
		[]string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug", "nolog"},
		false,
	))
}
