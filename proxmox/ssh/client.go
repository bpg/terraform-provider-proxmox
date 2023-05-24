/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/sftp"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helper"
)

type client struct {
	username    string
	password    string
	agent       bool
	agentSocket string
}

// NewClient creates a new SSH client.
func NewClient(username string, password string, agent bool, agentSocket string) (Client, error) {
	//goland:noinspection GoBoolExpressions
	if agent && runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "freebsd" {
		return nil, errors.New(
			"the ssh agent flag is only supported on POSIX systems, please set it to 'false'" +
				" or remove it from your provider configuration",
		)
	}

	return &client{
		username:    username,
		password:    password,
		agent:       agent,
		agentSocket: agentSocket,
	}, nil
}

// ExecuteNodeCommands executes commands on a given node.
func (c *client) ExecuteNodeCommands(ctx context.Context, nodeAddress string, commands []string) error {
	closeOrLogError := helper.CloseOrLogError(ctx)

	sshClient, err := c.openNodeShell(ctx, nodeAddress)
	if err != nil {
		return err
	}

	defer closeOrLogError(sshClient)

	sshSession, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
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

func (c *client) NodeUpload(
	ctx context.Context, nodeAddress string, remoteFileDir string,
	d *api.FileUploadRequest,
) error {
	// We need to upload all other files using SFTP due to API limitations.
	// Hopefully, this will not be required in future releases of Proxmox VE.
	tflog.Debug(ctx, "uploading file to datastore using SFTP", map[string]interface{}{
		"file_name":    d.FileName,
		"content_type": d.ContentType,
	})

	fileInfo, err := d.File.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()

	sshClient, err := c.openNodeShell(ctx, nodeAddress)
	if err != nil {
		return fmt.Errorf("failed to open SSH client: %w", err)
	}

	defer func(sshClient *ssh.Client) {
		e := sshClient.Close()
		if e != nil {
			tflog.Error(ctx, "failed to close SSH client", map[string]interface{}{
				"error": e,
			})
		}
	}(sshClient)

	if d.ContentType != "" {
		remoteFileDir = filepath.Join(remoteFileDir, d.ContentType)
	}

	remoteFilePath := strings.ReplaceAll(filepath.Join(remoteFileDir, d.FileName), `\`, `/`)

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	defer func(sftpClient *sftp.Client) {
		e := sftpClient.Close()
		if e != nil {
			tflog.Error(ctx, "failed to close SFTP client", map[string]interface{}{
				"error": e,
			})
		}
	}(sftpClient)

	err = sftpClient.MkdirAll(remoteFileDir)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", remoteFileDir, err)
	}

	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", remoteFilePath, err)
	}

	defer func(remoteFile *sftp.File) {
		e := remoteFile.Close()
		if e != nil {
			tflog.Error(ctx, "failed to close remote file", map[string]interface{}{
				"error": e,
			})
		}
	}(remoteFile)

	bytesUploaded, err := remoteFile.ReadFrom(d.File)
	if err != nil {
		return fmt.Errorf("failed to upload file %s: %w", remoteFilePath, err)
	}

	if bytesUploaded != fileSize {
		return fmt.Errorf("failed to upload file %s: uploaded %d bytes, expected %d bytes",
			remoteFilePath, bytesUploaded, fileSize)
	}

	tflog.Debug(ctx, "uploaded file to datastore", map[string]interface{}{
		"remote_file_path": remoteFilePath,
		"size":             bytesUploaded,
	})

	return nil
}

// openNodeShell establishes a new SSH connection to a node.
func (c *client) openNodeShell(ctx context.Context, nodeAddress string) (*ssh.Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the home directory: %w", err)
	}

	sshHost := fmt.Sprintf("%s:22", nodeAddress)

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
				defer helper.CloseOrLogError(ctx)(f)
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
		User:              c.username,
		Auth:              []ssh.AuthMethod{ssh.Password(c.password)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	tflog.Info(ctx, fmt.Sprintf("Agent is set to %t", c.agent))

	var sshClient *ssh.Client
	if c.agent {
		sshClient, err = c.createSSHClientAgent(ctx, cb, kh, sshHost)
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
		"user": c.username,
	})

	return sshClient, nil
}

// createSSHClientAgent establishes an ssh connection through the agent authentication mechanism.
func (c *client) createSSHClientAgent(
	ctx context.Context,
	cb ssh.HostKeyCallback,
	kh knownhosts.HostKeyCallback,
	sshHost string,
) (*ssh.Client, error) {
	if c.agentSocket == "" {
		return nil, errors.New("failed connecting to SSH agent socket: the socket file is not defined, " +
			"authentication will fall back to password")
	}

	conn, err := net.Dial("unix", c.agentSocket)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to SSH auth socket '%s': %w", c.agentSocket, err)
	}

	ag := agent.NewClient(conn)

	sshConfig := &ssh.ClientConfig{
		User:              c.username,
		Auth:              []ssh.AuthMethod{ssh.PublicKeysCallback(ag.Signers), ssh.Password(c.password)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	sshClient, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", sshHost, err)
	}

	tflog.Debug(ctx, "SSH connection established", map[string]interface{}{
		"host": sshHost,
		"user": c.username,
	})

	return sshClient, nil
}
