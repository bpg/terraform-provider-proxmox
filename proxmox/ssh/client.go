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
	"golang.org/x/net/proxy"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	// TrySudo is a shell function that tries to execute a command with sudo if the user has sudo permissions.
	TrySudo = `try_sudo(){ if [ $(sudo -n pvesm apiinfo 2>&1 | grep "APIVER" | wc -l) -gt 0 ]; then sudo $1; else $1; fi }`
)

// NewErrUserHasNoPermission creates a new error indicating that the SSH user does not have required permissions.
func NewErrUserHasNoPermission(username string) error {
	return fmt.Errorf("the SSH user '%s' does not have required permissions. "+
		"Make sure 'sudo' is installed and the user is configured in sudoers file. "+
		"Refer to the documentation for more details", username)
}

// Client is an interface for performing SSH requests against the Proxmox Nodes.
type Client interface {
	// Username returns the SSH username.
	Username() string

	// ExecuteNodeCommands executes a command on a node.
	ExecuteNodeCommands(ctx context.Context, nodeName string, commands []string) ([]byte, error)

	// NodeUpload uploads a file to a node.
	NodeUpload(ctx context.Context, nodeName string,
		remoteFileDir string, fileUploadRequest *api.FileUploadRequest) error
}

type client struct {
	username       string
	password       string
	agent          bool
	agentSocket    string
	privateKey     string
	socks5Server   string
	socks5Username string
	socks5Password string
	nodeResolver   NodeResolver
}

// NewClient creates a new SSH client.
func NewClient(
	username string, password string,
	agent bool, agentSocket string,
	privateKey string,
	socks5Server string, socks5Username string, socks5Password string,
	nodeResolver NodeResolver,
) (Client, error) {
	if agent && runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "freebsd" {
		return nil, errors.New(
			"the ssh agent flag is only supported on POSIX systems, please set it to 'false'" +
				" or remove it from your provider configuration",
		)
	}

	if (socks5Username != "" || socks5Password != "") && socks5Server == "" {
		return nil, errors.New("socks5 server is required when socks5 username or password is set")
	}

	if nodeResolver == nil {
		return nil, errors.New("node resolver is required")
	}

	return &client{
		username:       username,
		password:       password,
		agent:          agent,
		agentSocket:    agentSocket,
		privateKey:     privateKey,
		socks5Server:   socks5Server,
		socks5Username: socks5Username,
		socks5Password: socks5Password,
		nodeResolver:   nodeResolver,
	}, nil
}

func (c *client) Username() string {
	return c.username
}

// ExecuteNodeCommands executes commands on a given node.
func (c *client) ExecuteNodeCommands(ctx context.Context, nodeName string, commands []string) ([]byte, error) {
	node, err := c.nodeResolver.Resolve(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to find node endpoint: %w", err)
	}

	tflog.Debug(ctx, "executing commands on the node using SSH", map[string]interface{}{
		"node_address": node.Address,
		"node_port":    node.Port,
		"commands":     commands,
	})

	closeOrLogError := utils.CloseOrLogError(ctx)

	sshClient, err := c.openNodeShell(ctx, node)
	if err != nil {
		return nil, err
	}

	defer closeOrLogError(sshClient)

	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	defer closeOrLogError(sshSession)

	script := strings.Join(commands, "; ")

	output, err := sshSession.CombinedOutput(
		fmt.Sprintf(
			// explicitly use bash to support shell features like pipes and var assignment
			"/bin/bash -c '%s'",
			// shell script escaping for single quotes
			strings.ReplaceAll(script, `'`, `'"'"'`),
		),
	)
	if err != nil {
		return nil, errors.New(string(output))
	}

	return output, nil
}

func (c *client) NodeUpload(
	ctx context.Context,
	nodeName string,
	remoteFileDir string,
	d *api.FileUploadRequest,
) error {
	ip, err := c.nodeResolver.Resolve(ctx, nodeName)
	if err != nil {
		return fmt.Errorf("failed to find node endpoint: %w", err)
	}

	tflog.Debug(ctx, "uploading file to the node datastore using SFTP", map[string]interface{}{
		"node_address": ip,
		"remote_dir":   remoteFileDir,
		"file_name":    d.FileName,
		"content_type": d.ContentType,
	})

	fileInfo, err := d.File.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()

	sshClient, err := c.openNodeShell(ctx, ip)
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

	remoteFilePath := strings.ReplaceAll(filepath.Join(remoteFileDir, d.FileName), `\`, "/")

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
func (c *client) openNodeShell(ctx context.Context, node ProxmoxNode) (*ssh.Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the home directory: %w", err)
	}

	sshHost := fmt.Sprintf("%s:%d", node.Address, node.Port)

	sshPath := path.Join(homeDir, ".ssh")
	if _, err = os.Stat(sshPath); os.IsNotExist(err) {
		e := os.Mkdir(sshPath, 0o700)
		if e != nil && !os.IsExist(e) {
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
				defer utils.CloseOrLogError(ctx)(f)
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

	tflog.Info(ctx, fmt.Sprintf("agent is set to %t", c.agent))

	var sshClient *ssh.Client
	if c.agent {
		sshClient, err = c.createSSHClientAgent(ctx, cb, kh, sshHost)
		if err == nil {
			return sshClient, nil
		}

		tflog.Error(ctx, "Failed SSH connection through agent",
			map[string]interface{}{
				"error": err,
			})
	}

	if c.privateKey != "" {
		sshClient, err = c.createSSHClientWithPrivateKey(ctx, cb, kh, sshHost)
		if err == nil {
			return sshClient, nil
		}

		tflog.Error(ctx, "Failed SSH connection with private key",
			map[string]interface{}{
				"error": err,
			})

	}

	tflog.Info(ctx, "Falling back to password authentication for SSH connection")
	sshClient, err = c.createSSHClient(ctx, cb, kh, sshHost)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate user %q over SSH to %q. Please verify that ssh-agent is "+
			"correctly loaded with an authorized key via 'ssh-add -L' (NOTE: configurations in ~/.ssh/config are "+
			"not considered by the provider): %w", c.username, sshHost, err)
	}

	return sshClient, nil
}

func (c *client) createSSHClient(
	ctx context.Context,
	cb ssh.HostKeyCallback,
	kh knownhosts.HostKeyCallback,
	sshHost string,
) (*ssh.Client, error) {
	if c.password == "" {
		tflog.Error(ctx, "Using password authentication fallback for SSH connection, but the SSH password is empty")
	}

	sshConfig := &ssh.ClientConfig{
		User:              c.username,
		Auth:              []ssh.AuthMethod{ssh.Password(c.password)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	return c.connect(ctx, sshHost, sshConfig)
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

	return c.connect(ctx, sshHost, sshConfig)
}

func (c *client) createSSHClientWithPrivateKey(
	ctx context.Context,
	cb ssh.HostKeyCallback,
	kh knownhosts.HostKeyCallback,
	sshHost string,
) (*ssh.Client, error) {
	privateKey, err := ssh.ParsePrivateKey([]byte(c.privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:              c.username,
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
		HostKeyCallback:   cb,
		HostKeyAlgorithms: kh.HostKeyAlgorithms(sshHost),
	}

	return c.connect(ctx, sshHost, sshConfig)
}

func (c *client) connect(ctx context.Context, sshHost string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	if c.socks5Server != "" {
		sshClient, err := c.socks5SSHClient(sshHost, sshConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s via SOCKS5 proxy %s: %w", sshHost, c.socks5Server, err)
		}

		tflog.Debug(ctx, "SSH connection via SOCKS5 established", map[string]interface{}{
			"host":          sshHost,
			"socks5_server": c.socks5Server,
			"user":          c.username,
		})

		return sshClient, nil
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

func (c *client) socks5SSHClient(sshServerAddress string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", c.socks5Server, &proxy.Auth{
		User:     c.socks5Username,
		Password: c.socks5Password,
	}, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 proxy dialer: %w", err)
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s via SOCKS5 proxy %s: %w", sshServerAddress, c.socks5Server, err)
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client connection: %w", err)
	}

	return ssh.NewClient(sshConn, chans, reqs), nil
}
