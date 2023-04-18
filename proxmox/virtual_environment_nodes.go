/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

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
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/skeema/knownhosts"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetNodeTaskStatus retrieves the status of a node task.
func (c *VirtualEnvironmentClient) GetNodeTaskStatus(
	ctx context.Context,
	nodeName string,
	upid string,
) (*VirtualEnvironmentNodeGetTaskStatusResponseData, error) {
	resBody := &VirtualEnvironmentNodeGetTaskStatusResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/tasks/%s/status", url.PathEscape(nodeName), url.PathEscape(upid)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
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
		return nil, errors.New("the server did not include a data object in the response")
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
		return nil, errors.New("the server did not include a data object in the response")
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

	ur := strings.Split(c.Username, "@")

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
		err := kh(hostname, remote, key)
		if knownhosts.IsHostKeyChanged(err) {
			return fmt.Errorf("REMOTE HOST IDENTIFICATION HAS CHANGED for host %s! This may indicate a MitM attack", hostname)
		}

		if knownhosts.IsHostUnknown(err) {
			f, ferr := os.OpenFile(khPath, os.O_APPEND|os.O_WRONLY, 0o600)
			if ferr == nil {
				defer CloseOrLogError(ctx)(f)
				ferr = knownhosts.WriteKnownHost(f, hostname, remote, key)
			}
			if ferr == nil {
				tflog.Info(ctx, fmt.Sprintf("Added host %s to known_hosts", hostname))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to add host %s to known_hosts", hostname), map[string]interface{}{
					"error": err,
				})
			}
			return nil
		}
		return err
	})

	sshConfig := &ssh.ClientConfig{
		User:              ur[0],
		Auth:              []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	if c.Agent {
		
		sshAuthSock := os.Getenv("SSH_AUTH_SOCK")

		if sshAuthSock == "" {
			return nil, fmt.Errorf("failed connecting to SSH_AUTH_SOCK: environment variable is empty")
		}

		conn, err := net.Dial("unix", sshAuthSock)

		if err != nil {
			return nil, fmt.Errorf("failed connecting to SSH_AUTH_SOCK: %v", err)
		}

		ag := agent.NewClient(conn)

		sshConfig = &ssh.ClientConfig{
			User:              ur[0],
			Auth:              []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers)},
			HostKeyCallback:   cb,
			HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
		}
	}

	sshClient, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", sshHost, err)
	}

	tflog.Debug(ctx, "SSH connection established", map[string]interface{}{
		"host": sshHost,
		"user": ur[0],
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

// WaitForNodeTask waits for a specific node task to complete.
func (c *VirtualEnvironmentClient) WaitForNodeTask(
	ctx context.Context,
	nodeName string,
	upid string,
	timeout int,
	delay int,
) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			status, err := c.GetNodeTaskStatus(ctx, nodeName, upid)
			if err != nil {
				return err
			}

			if status.Status != "running" {
				if status.ExitCode != "OK" {
					return fmt.Errorf(
						"task \"%s\" on node \"%s\" failed to complete with error: %s",
						upid,
						nodeName,
						status.ExitCode,
					)
				}
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf(
		"timeout while waiting for task \"%s\" on node \"%s\" to complete",
		upid,
		nodeName,
	)
}
