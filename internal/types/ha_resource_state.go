/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/internal/validators"
)

// HAResourceState represents the requested state of a HA resource.
type HAResourceState int

// Ensure various interfaces are supported by the HA resource state type.
// NOTE: the global variable created here is only meant to be used in this block. There is, to my knowledge, no
// other way to enforce interface implementation at compile time unless the value is wrapped into a struct. Because
// of this, the linter is disabled.
var (
	//nolint:gochecknoglobals
	_haResourceStateValue HAResourceState
	_                     fmt.Stringer     = &_haResourceStateValue
	_                     json.Marshaler   = &_haResourceStateValue
	_                     json.Unmarshaler = &_haResourceStateValue
	_                     query.Encoder    = &_haResourceStateValue
)

const (
	// HAResourceStateStarted indicates that a HA resource should be started.
	HAResourceStateStarted HAResourceState = 0
	// HAResourceStateStopped indicates that a HA resource should be stopped, but that it should still be relocated
	// on node failure.
	HAResourceStateStopped HAResourceState = 1
	// HAResourceStateDisabled indicates that a HA resource should be stopped. No relocation should occur on node failure.
	HAResourceStateDisabled HAResourceState = 2
	// HAResourceStateIgnored indicates that a HA resource is not managed by the cluster resource manager. No relocation
	// or status change will occur.
	HAResourceStateIgnored HAResourceState = 3
)

// ParseHAResourceState converts the string representation of a HA resource state into the corresponding
// enum value. An error is returned if the input string does not match any known state. This function also
// parses the `enabled` value which is an alias for `started`.
func ParseHAResourceState(input string) (HAResourceState, error) {
	switch input {
	case "started":
		return HAResourceStateStarted, nil
	case "enabled":
		return HAResourceStateStarted, nil
	case "stopped":
		return HAResourceStateStopped, nil
	case "disabled":
		return HAResourceStateDisabled, nil
	case "ignored":
		return HAResourceStateIgnored, nil
	default:
		return HAResourceStateIgnored, fmt.Errorf("illegal HA resource state '%s'", input)
	}
}

// HAResourceStateValidator returns a new HA resource state validator.
func HAResourceStateValidator() validator.String {
	return validators.NewParseValidator(ParseHAResourceState, "value must be a valid HA resource state")
}

// String converts a HAResourceState value into a string.
func (s HAResourceState) String() string {
	switch s {
	case HAResourceStateStarted:
		return "started"
	case HAResourceStateStopped:
		return "stopped"
	case HAResourceStateDisabled:
		return "disabled"
	case HAResourceStateIgnored:
		return "ignored"
	default:
		panic(fmt.Sprintf("unknown HA resource state value: %d", s))
	}
}

// MarshalJSON marshals a HA resource state into JSON value.
func (s HAResourceState) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(s.String())
	if err != nil {
		return nil, fmt.Errorf("cannot marshal HA resource state: %w", err)
	}

	return bytes, nil
}

// UnmarshalJSON unmarshals a Proxmox HA resource state.
func (s *HAResourceState) UnmarshalJSON(b []byte) error {
	var stateString string

	err := json.Unmarshal(b, &stateString)
	if err != nil {
		return fmt.Errorf("cannot unmarshal HA resource state: %w", err)
	}

	state, err := ParseHAResourceState(stateString)
	if err == nil {
		*s = state
	}

	return err
}

// EncodeValues encodes a HA resource state field into an URL-encoded set of values.
func (s HAResourceState) EncodeValues(key string, v *url.Values) error {
	v.Add(key, s.String())
	return nil
}

// ToValue converts a HA resource state into a Terraform value.
func (s HAResourceState) ToValue() types.String {
	return types.StringValue(s.String())
}
