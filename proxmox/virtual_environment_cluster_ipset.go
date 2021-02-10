/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// CreateIPSet create an IPSet
func (c *VirtualEnvironmentClient) CreateIPSet(d *VirtualEnvironmentClusterIPSetCreateRequestBody) error  {
	return c.DoRequest(hmPOST, "cluster/firewall/ipset", d, nil)
}

// Add IP or Network to IPSet
func (c *VirtualEnvironmentClient) AddCIDRToIPSet(id string, d *VirtualEnvironmentClusterIPSetGetResponseData) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("cluster/firewall/ipset/%s/", url.PathEscape(id)), d, nil)
}

// UpdateIPSet updates an IPSet.
func (c *VirtualEnvironmentClient) UpdateIPSet(d *VirtualEnvironmentClusterIPSetUpdateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprint("cluster/firewall/ipset/"), d, nil)
}

// DeleteIPSet delete an IPSet
func (c *VirtualEnvironmentClient) DeleteIPSet(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)), nil, nil)
}

// DeleteIPSetContent remove IP or Network from IPSet.
func (c *VirtualEnvironmentClient) DeleteIPSetContent(id string, cidr string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("cluster/firewall/ipset/%s/%s", url.PathEscape(id), url.PathEscape(cidr)), nil, nil)
}

// GetListIPSetContent retrieve a list of IPSet content
func (c *VirtualEnvironmentClient) GetListIPSetContent(id string) ([]*VirtualEnvironmentClusterIPSetGetResponseData, error) {
	resBody := &VirtualEnvironmentClusterIPSetGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetListIPSets retrieves list of IPSets.
func (c *VirtualEnvironmentClient) GetListIPSets() (*VirtualEnvironmentClusterIPSetListResponseBody, error) {
	resBody := &VirtualEnvironmentClusterIPSetListResponseBody{}
	err := c.DoRequest(hmGET, "cluster/firewall/ipset", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody, nil
}