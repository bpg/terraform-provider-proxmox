/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ExecuteNodeCommands executes commands on a given node.
func (c *VirtualEnvironmentClient) ExecuteNodeCommands(nodeName string, commands []string) error {
	// We must first retrieve the IP address of the node as we need to bypass the API and use SSH instead.
	networkDevices, err := c.ListNodeNetworkDevices(nodeName)

	if err != nil {
		return err
	}

	nodeAddress := ""

	for _, d := range networkDevices {
		if d.Address != nil {
			nodeAddress = *d.Address
			break
		}
	}

	if nodeAddress == "" {
		return fmt.Errorf("Failed to determine the IP address of node \"%s\"", nodeName)
	}

	// We can now go ahead and execute the commands using SSH.
	// Hopefully, the developers will add this feature to the REST API at some point.
	ur := strings.Split(c.Username, "@")

	sshConfig := &ssh.ClientConfig{
		User:            ur[0],
		Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", nodeAddress+":22", sshConfig)

	if err != nil {
		return err
	}

	defer sshClient.Close()

	sshSession, err := sshClient.NewSession()

	if err != nil {
		return err
	}

	defer sshSession.Close()

	_, err = sshSession.CombinedOutput(
		fmt.Sprintf(
			"/bin/bash -c '%s'",
			strings.ReplaceAll(strings.Join(commands, " && "), "'", "'\"'\"'"),
		),
	)

	if err != nil {
		return err
	}

	return nil
}

// ListNodeNetworkDevices retrieves a list of network devices for a specific nodes.
func (c *VirtualEnvironmentClient) ListNodeNetworkDevices(nodeName string) ([]*VirtualEnvironmentNodeNetworkDeviceListResponseData, error) {
	resBody := &VirtualEnvironmentNodeNetworkDeviceListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/network", url.PathEscape(nodeName)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Priority < resBody.Data[j].Priority
	})

	return resBody.Data, nil
}

// ListNodes retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListNodes() ([]*VirtualEnvironmentNodeListResponseData, error) {
	resBody := &VirtualEnvironmentNodeListResponseBody{}
	err := c.DoRequest(hmGET, "nodes", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}
