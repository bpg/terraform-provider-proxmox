/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Options interface {
	SetOptions(ctx context.Context, d *OptionsPutRequestBody) error
	GetOptions(ctx context.Context) (*OptionsGetResponseData, error)
}

type OptionsPutRequestBody struct {
	EBTables     *types.CustomBool   `json:"ebtables,omitempty"      url:"ebtables,omitempty,int"`
	Enable       *types.CustomBool   `json:"enable,omitempty"        url:"enable,omitempty,int"`
	LogRateLimit *CustomLogRateLimit `json:"log_ratelimit,omitempty" url:"log_ratelimit,omitempty"`
	PolicyIn     *string             `json:"policy_in,omitempty"     url:"policy_in,omitempty"`
	PolicyOut    *string             `json:"policy_out,omitempty"    url:"policy_out,omitempty"`
}

type CustomLogRateLimit struct {
	Enable types.CustomBool `json:"enable,omitempty" url:"enable,omitempty,int"`
	Burst  *int             `json:"burst,omitempty"  url:"burst,omitempty,int"`
	Rate   *string          `json:"rate,omitempty"   url:"rate,omitempty"`
}

type OptionsGetResponseBody struct {
	Data *OptionsGetResponseData `json:"data,omitempty"`
}

type OptionsGetResponseData struct {
	EBTables     *types.CustomBool   `json:"ebtables"      url:"ebtables, int"`
	Enable       *types.CustomBool   `json:"enable"        url:"enable,int"`
	LogRateLimit *CustomLogRateLimit `json:"log_ratelimit" url:"log_ratelimit"`
	PolicyIn     *string             `json:"policy_in"     url:"policy_in"`
	PolicyOut    *string             `json:"policy_out"    url:"policy_out"`
}

// EncodeValues converts a CustomWatchdogDevice struct to a URL vlaue.
func (r *CustomLogRateLimit) EncodeValues(key string, v *url.Values) error {
	var values []string

	if r.Enable {
		values = append(values, "enable=1")
	} else {
		values = append(values, "enable=0")
	}

	if r.Burst != nil {
		values = append(values, fmt.Sprintf("burst=%d", *r.Burst))
	}

	if r.Rate != nil {
		values = append(values, fmt.Sprintf("rate=%s", *r.Rate))
	}

	v.Add(key, strings.Join(values, ","))

	return nil
}

func (r *CustomLogRateLimit) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("error unmarshaling json: %w", err)
	}

	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")

	for _, p := range pairs {
		v := strings.Split(strings.TrimSpace(p), "=")

		if len(v) == 1 {
			r.Enable = v[0] == "1"
		} else if len(v) == 2 {
			switch v[0] {
			case "enable":
				r.Enable = v[1] == "1"
			case "burst":
				iv, err := strconv.Atoi(v[1])
				if err != nil {
					return fmt.Errorf("error converting burst to int: %w", err)
				}
				r.Burst = &iv
			case "rate":
				r.Rate = &v[1]
			}
		}
	}

	return nil
}
