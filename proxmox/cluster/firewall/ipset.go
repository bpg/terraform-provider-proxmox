/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

// CreateIPSet create an IPSet
func (a *API) CreateIPSet(ctx context.Context, d *IPSetCreateRequestBody) error {
	err := a.DoRequest(ctx, http.MethodPost, "cluster/firewall/ipset", d, nil)
	if err != nil {
		return fmt.Errorf("error creating IPSet: %w", err)
	}
	return nil
}

// AddCIDRToIPSet adds IP or Network to IPSet
func (a *API) AddCIDRToIPSet(ctx context.Context, id string, d IPSetGetResponseData) error {
	err := a.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("cluster/firewall/ipset/%s/", url.PathEscape(id)),
		&d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error adding CIDR to IPSet: %w", err)
	}
	return nil
}

// UpdateIPSet updates an IPSet.
func (a *API) UpdateIPSet(ctx context.Context, d *IPSetUpdateRequestBody) error {
	err := a.DoRequest(ctx, http.MethodPost, "cluster/firewall/ipset/", d, nil)
	if err != nil {
		return fmt.Errorf("error updating IPSet: %w", err)
	}
	return nil
}

// DeleteIPSet delete an IPSet
func (a *API) DeleteIPSet(ctx context.Context, id string) error {
	err := a.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting IPSet %s: %w", id, err)
	}
	return nil
}

// DeleteIPSetContent remove IP or Network from IPSet.
func (a *API) DeleteIPSetContent(ctx context.Context, id string, cidr string) error {
	err := a.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/ipset/%s/%s", url.PathEscape(id), url.PathEscape(cidr)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting IPSet content %s: %w", id, err)
	}
	return nil
}

// GetIPSetContent retrieve a list of IPSet content
func (a *API) GetIPSetContent(ctx context.Context, id string) ([]*IPSetGetResponseData, error) {
	resBody := &IPSetGetResponseBody{}
	err := a.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting IPSet content: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListIPSets retrieves list of IPSets.
func (a *API) ListIPSets(ctx context.Context) ([]*IPSetListResponseData, error) {
	resBody := &IPSetListResponseBody{}
	err := a.DoRequest(ctx, http.MethodGet, "cluster/firewall/ipset", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error getting IPSet list: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}
