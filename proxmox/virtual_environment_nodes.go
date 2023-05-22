/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/skeema/knownhosts"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ExecuteNodeCommands executes commands on a given node.
func (c *VirtualEnvironmentClient) ExecuteNodeCommands(
	ctx context.Context,
	nodeName string,
	commands []string,
) error {
	closeOrLogError := CloseOrLogError(ctx)

	sshClient, err := c.OpenNodeShell(ctx, nodeName)
	if err != nil {
		return err
	}

	defer closeOrLogError(sshClient)

	sshSession, err := sshClient.NewSession()
	if err != nil {
		return err
	}

	defer closeOrLogError(sshSession)

	script := strings.Join(commands, " && \\\n")
	output, err := sshSession.CombinedOutput(
		fmt.Sprintf(
			"/bin/bash -c '%s'",
			strings.ReplaceAll(script, "'", "'\"'\"'"),
		),
	)
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}

// GetNodeIP retrieves the IP address of a node.
func (c *VirtualEnvironmentClient) GetNodeIP(
	ctx context.Context,
	nodeName string,
) (*string, error) {
	networkDevices, err := c.ListNodeNetworkDevices(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	nodeAddress := ""

	for _, d := range networkDevices {
		if d.Address != nil {
			nodeAddress = *d.Address
			break
		}
	}

	if nodeAddress == "" {
		return nil, fmt.Errorf("failed to determine the IP address of node \"%s\"", nodeName)
	}

	nodeAddressParts := strings.Split(nodeAddress, "/")

	return &nodeAddressParts[0], nil
}

// GetNodeTime retrieves the time information for a node.
func (c *VirtualEnvironmentClient) GetNodeTime(
	ctx context.Context,
	nodeName string,
) (*VirtualEnvironmentNodeGetTimeResponseData, error) {
	resBody := &VirtualEnvironmentNodeGetTimeResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/time", url.PathEscape(nodeName)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ListNodeNetworkDevices retrieves a list of network devices for a specific nodes.
func (c *VirtualEnvironmentClient) ListNodeNetworkDevices(
	ctx context.Context,
	nodeName string,
) ([]*VirtualEnvironmentNodeNetworkDeviceListResponseData, error) {
	resBody := &VirtualEnvironmentNodeNetworkDeviceListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/network", url.PathEscape(nodeName)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Priority < resBody.Data[j].Priority
	})

	return resBody.Data, nil
}

// ListNodes retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListNodes(
	ctx context.Context,
) ([]*VirtualEnvironmentNodeListResponseData, error) {
	resBody := &VirtualEnvironmentNodeListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "nodes", nil, resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// OpenNodeShell establishes a new SSH connection to a node.
func (c *VirtualEnvironmentClient) OpenNodeShell(
	ctx context.Context,
	nodeName string,
) (*ssh.Client, error) {
	nodeAddress, err := c.GetNodeIP(ctx, nodeName)
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the home directory: %w", err)
	}

	sshHost := fmt.Sprintf("%s:22", *nodeAddress)

	sshPath := path.Join(homeDir, ".ssh")
	if _, err = os.Stat(sshPath); os.IsNotExist(err) {
		e := os.Mkdir(sshPath, 0o700)
		if e != nil {
			return nil, fmt.Errorf("failed to create %s: %w", sshPath, e)
		}
	}

	khPath := path.Join(sshPath, "known_hosts")
	if _, err = os.Stat(khPath); os.IsNotExist(err) {
		e := os.WriteFile(khPath, []byte{}, 0o600)
		if e != nil {
			return nil, fmt.Errorf("failed to create %s: %w", khPath, e)
		}
	}

	kh, err := knownhosts.New(khPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", khPath, err)
	}

	// Create a custom permissive hostkey callback which still errors on hosts
	// with changed keys, but allows unknown hosts and adds them to known_hosts
	cb := ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		kherr := kh(hostname, remote, key)
		if knownhosts.IsHostKeyChanged(kherr) {
			return fmt.Errorf("REMOTE HOST IDENTIFICATION HAS CHANGED for host %s! This may indicate a MitM attack", hostname)
		}

		if knownhosts.IsHostUnknown(kherr) {
			f, ferr := os.OpenFile(khPath, os.O_APPEND|os.O_WRONLY, 0o600)
			if ferr == nil {
				defer CloseOrLogError(ctx)(f)
				ferr = knownhosts.WriteKnownHost(f, hostname, remote, key)
			}
			if ferr == nil {
				tflog.Info(ctx, fmt.Sprintf("Added host %s to known_hosts", hostname))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to add host %s to known_hosts", hostname), map[string]interface{}{
					"error": kherr,
				})
			}
			return nil
		}
		return kherr
	})

	sshConfig := &ssh.ClientConfig{
		User:              c.SSHUsername,
		Auth:              []ssh.AuthMethod{ssh.Password(c.SSHPassword)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	tflog.Info(ctx, fmt.Sprintf("Agent is set to %t", c.SSHAgent))

	var sshClient *ssh.Client
	if c.SSHAgent {
		sshClient, err = c.CreateSSHClientAgent(ctx, cb, kh, sshHost)
		if err != nil {
			tflog.Error(ctx, "Failed ssh connection through agent, "+
				"falling back to password authentication",
				map[string]interface{}{
					"error": err,
				})
		} else {
			return sshClient, nil
		}
	}

	sshClient, err = ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", sshHost, err)
	}

	tflog.Debug(ctx, "SSH connection established", map[string]interface{}{
		"host": sshHost,
		"user": c.SSHUsername,
	})

	return sshClient, nil
}

// CreateSSHClientAgent establishes an ssh connection through the agent authentication mechanism.
func (c *VirtualEnvironmentClient) CreateSSHClientAgent(
	ctx context.Context,
	cb ssh.HostKeyCallback,
	kh knownhosts.HostKeyCallback,
	sshHost string,
) (*ssh.Client, error) {
	if c.SSHAgentSocket == "" {
		return nil, errors.New("failed connecting to SSH agent socket: the socket file is not defined, " +
			"authentication will fall back to password")
	}

	conn, err := net.Dial("unix", c.SSHAgentSocket)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to SSH auth socket '%s': %w", c.SSHAgentSocket, err)
	}

	ag := agent.NewClient(conn)

	sshConfig := &ssh.ClientConfig{
		User:              c.SSHUsername,
		Auth:              []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers), ssh.Password(c.SSHPassword)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	sshClient, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", sshHost, err)
	}

	tflog.Debug(ctx, "SSH connection established", map[string]interface{}{
		"host": sshHost,
		"user": c.SSHUsername,
	})

	return sshClient, nil
}

// UpdateNodeTime updates the time on a node.
func (c *VirtualEnvironmentClient) UpdateNodeTime(
	ctx context.Context,
	nodeName string,
	d *VirtualEnvironmentNodeUpdateTimeRequestBody,
) error {
	return c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("nodes/%s/time", url.PathEscape(nodeName)), d, nil)
}
