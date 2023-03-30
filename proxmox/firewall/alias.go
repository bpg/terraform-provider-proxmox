/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "context"

type Alias interface {
	CreateAlias(ctx context.Context, d *AliasCreateRequestBody) error
	DeleteAlias(ctx context.Context, name string) error
	GetAlias(ctx context.Context, name string) (*AliasGetResponseData, error)
	ListAliases(ctx context.Context) ([]*AliasGetResponseData, error)
	UpdateAlias(ctx context.Context, name string, d *AliasUpdateRequestBody) error
}

// AliasCreateRequestBody contains the data for an alias create request.
type AliasCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
	CIDR    string  `json:"cidr"              url:"cidr"`
}

// AliasGetResponseBody contains the body from an alias get response.
type AliasGetResponseBody struct {
	Data *AliasGetResponseData `json:"data,omitempty"`
}

// AliasGetResponseData contains the data from an alias get response.
type AliasGetResponseData struct {
	Comment   *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name      string  `json:"name"              url:"name"`
	CIDR      string  `json:"cidr"              url:"cidr"`
	Digest    *string `json:"digest"            url:"digest"`
	IPVersion int     `json:"ipversion"         url:"ipversion"`
}

// AliasListResponseBody contains the data from an alias get response.
type AliasListResponseBody struct {
	Data []*AliasGetResponseData `json:"data,omitempty"`
}

// AliasUpdateRequestBody contains the data for an alias update request.
type AliasUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ReName  string  `json:"rename"            url:"rename"`
	CIDR    string  `json:"cidr"              url:"cidr"`
}
