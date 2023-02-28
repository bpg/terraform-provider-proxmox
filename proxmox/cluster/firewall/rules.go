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
	"net/http"
	"net/url"
)

// CreateRule creates a firewall rule.
func (a *API) CreateRule(ctx context.Context, d *RuleCreateRequestBody) error {
	err := a.DoRequest(ctx, http.MethodPost, "cluster/firewall/rules", d, nil)
	return fmt.Errorf("error creating firewall rule: %w", err)
}

// CreateGroupRule creates a security group firewall rule.
func (a *API) CreateGroupRule(ctx context.Context, group string, d *RuleCreateRequestBody) error {
	err := a.DoRequest(ctx,
		http.MethodPost,
		fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)),
		d,
		nil,
	)
	return fmt.Errorf("error creating security group: %w", err)
}

// GetRule retrieves a firewall rule.
func (a *API) GetRule(ctx context.Context, pos int) (*RuleGetResponseData, error) {
	resBody := &RuleGetResponseBody{}
	err := a.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/rules/%d", pos),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rule %d: %w", pos, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetGroupRule retrieves a security group firewall rule.
func (a *API) GetGroupRule(ctx context.Context, group string, pos int) (*RuleGetResponseData, error) {
	resBody := &RuleGetResponseBody{}
	err := a.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rule %d: %w", pos, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListRules retrieves a list of firewall rules.
func (a *API) ListRules(ctx context.Context) ([]*RuleListResponseData, error) {
	resBody := &RuleListResponseBody{}
	err := a.DoRequest(ctx, http.MethodGet, "cluster/firewall/rules", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListGroupRules retrieves a list of security group firewall rules.
func (a *API) ListGroupRules(ctx context.Context, group string) ([]*RuleListResponseData, error) {
	resBody := &RuleListResponseBody{}
	err := a.DoRequest(ctx, http.MethodGet, fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateRule updates a firewall rule.
func (a *API) UpdateRule(ctx context.Context, pos int, d *RuleUpdateRequestBody) error {
	err := a.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("cluster/firewall/rules/%d", pos),
		d,
		nil,
	)
	return fmt.Errorf("error updating firewall rule %d: %w", pos, err)
}

// UpdateGroupRule updates a security group firewall rule.
func (a *API) UpdateGroupRule(ctx context.Context, group string, pos int, d *RuleUpdateRequestBody) error {
	err := a.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		d,
		nil,
	)
	return fmt.Errorf("error updating firewall rule %d: %w", pos, err)
}

// DeleteRule deletes a firewall rule.
func (a *API) DeleteRule(ctx context.Context, pos int) error {
	err := a.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("cluster/firewall/rules/%d", pos), nil, nil)
	return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
}

// DeleteGroupRule deletes a security group firewall rule.
func (a *API) DeleteGroupRule(ctx context.Context, group string, pos int) error {
	err := a.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		nil,
		nil,
	)
	return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
}
