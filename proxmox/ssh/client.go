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
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

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
	// The result is cached both in the provider (across shell invocations) and in _try_sudo_use_sudo (within a script).
	TrySudo = `_try_sudo_check(){ ` +
		`if [ -z "${_try_sudo_cached:-}" ]; then ` +
		`_try_sudo_use_sudo="${TRY_SUDO_USE_SUDO:-0}"; ` +
		`_try_sudo_cached=1; ` +
		`fi; ` +
		`}; ` +
		`try_sudo(){ ` +
		`_try_sudo_check; ` +
		`if [ "$_try_sudo_use_sudo" = "1" ]; then ` +
		`sudo "$@"; ` +
		`else ` +
		`"$@"; ` +
		`fi; ` +
		`}`
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

	// NodeStreamUpload uploads a file to a node by streaming its content over SSH.
	NodeStreamUpload(ctx context.Context, nodeName string,
		remoteFileDir string, fileUploadRequest *api.FileUploadRequest) error
}

type client struct {
	username        string
	password        string
	agent           bool
	agentSocket     string
	agentForwarding bool
	privateKey      string
	socks5Server    string
	socks5Username  string
	socks5Password  string
	nodeResolver    NodeResolver
	sudoCache       map[string]bool
	sudoCacheMu     sync.RWMutex
}

// NewClient creates a new SSH client.
func NewClient(
	username string, password string,
	agent bool, agentSocket string, agentForwarding bool,
	privateKey string,
	socks5Server string, socks5Username string, socks5Password string,
	nodeResolver NodeResolver,
) (Client, error) {
	if agent &&
		runtime.GOOS != "linux" &&
		runtime.GOOS != "darwin" &&
		runtime.GOOS != "freebsd" &&
		runtime.GOOS != "windows" {
		return nil, errors.New(
			"the ssh agent flag is only supported on POSIX and Windows systems, please set it to 'false'" +
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
		username:        username,
		password:        password,
		agent:           agent,
		agentSocket:     agentSocket,
		agentForwarding: agentForwarding,
		privateKey:      privateKey,
		socks5Server:    socks5Server,
		socks5Username:  socks5Username,
		socks5Password:  socks5Password,
		nodeResolver:    nodeResolver,
		sudoCache:       make(map[string]bool),
		sudoCacheMu:     sync.RWMutex{},
	}, nil
}

func (c *client) Username() string {
	return c.username
}

// getSudoAvailability gets the cached sudo availability or checks and caches it.
// Uses double-checked locking to prevent race conditions when multiple goroutines
// check the same uncached node simultaneously.
func (c *client) getSudoAvailability(ctx context.Context, nodeName string) (bool, error) {
	c.sudoCacheMu.RLock()
	cached, found := c.sudoCache[nodeName]
	c.sudoCacheMu.RUnlock()

	if found {
		return cached, nil
	}

	c.sudoCacheMu.Lock()
	defer c.sudoCacheMu.Unlock()

	// Re-check the cache after acquiring the write lock in case another goroutine
	// populated it while we were waiting.
	cached, found = c.sudoCache[nodeName]
	if found {
		return cached, nil
	}

	checkCmd := `if [ "$(id -u)" = "0" ]; then echo "0"; ` +
		`elif [ $(sudo -n /sbin/pvesm apiinfo 2>&1 | grep "APIVER" | wc -l) -gt 0 ] ` +
		`|| sudo -n /sbin/qm --help >/dev/null 2>&1; then echo "1"; else echo "0"; fi`

	node, err := c.nodeResolver.Resolve(ctx, nodeName)
	if err != nil {
		return false, fmt.Errorf("failed to find node endpoint: %w", err)
	}

	sshClient, err := c.openNodeShell(ctx, node)
	if err != nil {
		return false, err
	}

	defer func() {
		if e := sshClient.Close(); e != nil {
			tflog.Warn(ctx, "failed to close SSH client", map[string]any{"error": e})
		}
	}()

	output, err := c.executeCommands(ctx, sshClient, []string{checkCmd})
	if err != nil {
		return false, fmt.Errorf("failed to check sudo availability: %w", err)
	}

	shouldUseSudo := strings.TrimSpace(string(output)) == "1"
	c.sudoCache[nodeName] = shouldUseSudo

	return shouldUseSudo, nil
}

// ExecuteNodeCommands executes commands on a given node.
func (c *client) ExecuteNodeCommands(ctx context.Context, nodeName string, commands []string) ([]byte, error) {
	commandsStr := strings.Join(commands, "; ")
	needsSudoCheck := strings.Contains(commandsStr, "try_sudo")

	sudoValue := c.getSudoValue(ctx, nodeName, needsSudoCheck)

	node, err := c.nodeResolver.Resolve(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to find node endpoint: %w", err)
	}

	tflog.Debug(ctx, "executing commands on the node using SSH", map[string]any{
		"node_address": node.Address,
		"node_port":    node.Port,
		"commands":     commands,
	})

	sshClient, err := c.openNodeShell(ctx, node)
	if err != nil {
		return nil, err
	}

	defer func() {
		if e := sshClient.Close(); e != nil {
			tflog.Warn(ctx, "failed to close SSH client", map[string]any{"error": e})
		}
	}()

	execCommands := commands
	if sudoValue != "" {
		execCommands = append([]string{fmt.Sprintf("export TRY_SUDO_USE_SUDO=%s", sudoValue)}, commands...)
	}

	return c.executeCommands(ctx, sshClient, execCommands)
}

// getSudoValue returns the sudo value to use, or empty string if check failed.
func (c *client) getSudoValue(ctx context.Context, nodeName string, needsCheck bool) string {
	if !needsCheck {
		return ""
	}

	shouldUseSudo, err := c.getSudoAvailability(ctx, nodeName)
	if err != nil {
		tflog.Warn(ctx, "failed to check sudo availability, proceeding without sudo", map[string]any{
			"error": err,
		})

		return ""
	}

	if shouldUseSudo {
		return "1"
	}

	return "0"
}

func (c *client) openSession(ctx context.Context, sshClient *ssh.Client) (*ssh.Session, func(), error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	closer := func() {
		e := sshSession.Close()
		if e != nil && !errors.Is(e, io.EOF) {
			tflog.Warn(ctx, "failed to close SSH session", map[string]any{
				"error": e,
			})
		}
	}

	if c.agentForwarding {
		tflog.Debug(ctx, "Requesting SSH agent forwarding")

		err = agent.RequestAgentForwarding(sshSession)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to request agent forwarding: %w", err)
		}

		if err = agent.ForwardToRemote(sshClient, c.agentSocket); err != nil {
			return nil, nil, fmt.Errorf("failed to forward agent connection to remote: %w", err)
		}
	}

	return sshSession, closer, nil
}

func (c *client) executeCommands(ctx context.Context, sshClient *ssh.Client, commands []string) ([]byte, error) {
	sshSession, closer, err := c.openSession(ctx, sshClient)
	defer closer()

	if err != nil {
		return nil, fmt.Errorf("failed to open SSH session: %w", err)
	}

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

	tflog.Debug(ctx, "uploading file to the node datastore using SFTP", map[string]any{
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
			tflog.Warn(ctx, "failed to close SSH client", map[string]any{
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
			tflog.Warn(ctx, "failed to close SFTP client", map[string]any{
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
			tflog.Warn(ctx, "failed to close remote file", map[string]any{
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

	tflog.Debug(ctx, "uploaded file to datastore", map[string]any{
		"remote_file_path": remoteFilePath,
		"size":             bytesUploaded,
	})

	return nil
}

func (c *client) NodeStreamUpload(
	ctx context.Context,
	nodeName string,
	remoteFileDir string,
	d *api.FileUploadRequest,
) error {
	ip, err := c.nodeResolver.Resolve(ctx, nodeName)
	if err != nil {
		return fmt.Errorf("failed to find node endpoint: %w", err)
	}

	tflog.Debug(ctx, "uploading file to the node datastore via SSH input stream ", map[string]any{
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
			tflog.Warn(ctx, "failed to close SSH client", map[string]any{
				"error": e,
			})
		}
	}(sshClient)

	if d.ContentType != "" {
		remoteFileDir = filepath.Join(remoteFileDir, d.ContentType)
	}

	remoteFilePath := strings.ReplaceAll(filepath.Join(remoteFileDir, d.FileName), `\`, "/")

	err = c.uploadFile(ctx, sshClient, d, remoteFilePath, nodeName)
	if err != nil {
		return err
	}

	err = c.checkUploadedFile(ctx, sshClient, remoteFilePath, fileSize)
	if err != nil {
		return err
	}

	if d.Mode != "" {
		parsedFileMode, parseErr := strconv.ParseUint(d.Mode, 8, 12)
		if parseErr != nil {
			return fmt.Errorf("failed to parse file mode %q: %w", d.Mode, parseErr)
		}

		mode := uint32(parsedFileMode)
		if err = c.changeModeUploadedFile(ctx, sshClient, remoteFilePath, os.FileMode(mode)); err != nil {
			return err
		}
	}

	tflog.Debug(ctx, "uploaded file to datastore", map[string]any{
		"remote_file_path": remoteFilePath,
	})

	return nil
}

func (c *client) uploadFile(
	ctx context.Context,
	sshClient *ssh.Client,
	req *api.FileUploadRequest,
	remoteFilePath string,
	nodeName string,
) error {
	sshSession, closer, err := c.openSession(ctx, sshClient)
	defer closer()

	if err != nil {
		return fmt.Errorf("failed to open SSH session: %w", err)
	}

	sudoValue := c.getSudoValue(ctx, nodeName, nodeName != "")

	sudoEnv := ""
	if sudoValue != "" {
		sudoEnv = fmt.Sprintf("export TRY_SUDO_USE_SUDO=%s; ", sudoValue)
	}

	sshSession.Stdin = req.File

	cmd := fmt.Sprintf(`%s%s; try_sudo /usr/bin/tee %s`, sudoEnv, TrySudo, remoteFilePath)

	output, err := sshSession.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("error transferring file: %s", string(output))
	}

	return nil
}

func (c *client) checkUploadedFile(
	ctx context.Context,
	sshClient *ssh.Client,
	remoteFilePath string,
	fileSize int64,
) error {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	defer func(sftpClient *sftp.Client) {
		e := sftpClient.Close()
		if e != nil {
			tflog.Warn(ctx, "failed to close SFTP client", map[string]any{
				"error": e,
			})
		}
	}(sftpClient)

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file %s: %w", remoteFilePath, err)
	}

	remoteStat, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to read remote file %s: %w", remoteFilePath, err)
	}

	bytesUploaded := remoteStat.Size()
	if bytesUploaded != fileSize {
		return fmt.Errorf("failed to upload file %s: uploaded %d bytes, expected %d bytes",
			remoteFilePath, bytesUploaded, fileSize)
	}

	return nil
}

func (c *client) changeModeUploadedFile(
	ctx context.Context,
	sshClient *ssh.Client,
	remoteFilePath string,
	fileMode os.FileMode,
) error {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	defer func(sftpClient *sftp.Client) {
		e := sftpClient.Close()
		if e != nil {
			tflog.Warn(ctx, "failed to close SFTP client", map[string]any{
				"error": e,
			})
		}
	}(sftpClient)

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file %s: %w", remoteFilePath, err)
	}

	remoteStat, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to read remote file %s: %w", remoteFilePath, err)
	}

	if err = sftpClient.Chmod(remoteFilePath, fileMode); err != nil {
		return fmt.Errorf("failed to change file mode of remote file from %#o (%s) to %#o (%s): %w",
			remoteStat.Mode().Perm(), remoteStat.Mode(), fileMode.Perm(), fileMode, err)
	}

	tflog.Debug(ctx, "changed mode of uploaded file", map[string]any{
		"before": fmt.Sprintf("%#o (%s)", remoteStat.Mode().Perm(), remoteStat.Mode()),
		"after":  fmt.Sprintf("%#o (%s)", fileMode.Perm(), fileMode),
	})

	return nil
}

// openNodeShell establishes a new SSH connection to a node.
func (c *client) openNodeShell(ctx context.Context, node ProxmoxNode) (*ssh.Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine the home directory: %w", err)
	}

	var sshHost string
	if strings.Contains(node.Address, ":") {
		// IPv6
		sshHost = fmt.Sprintf("[%s]:%d", node.Address, node.Port)
	} else {
		// IPv4
		sshHost = fmt.Sprintf("%s:%d", node.Address, node.Port)
	}

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

	kh, err := knownhosts.NewDB(khPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", khPath, err)
	}

	// Create a custom permissive host key callback which still errors on hosts
	// with changed keys, but allows unknown hosts and adds them to known_hosts
	cb := ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		khErr := kh.HostKeyCallback()(hostname, remote, key)
		if knownhosts.IsHostKeyChanged(khErr) {
			return fmt.Errorf("REMOTE HOST IDENTIFICATION HAS CHANGED for host %s! This may indicate a MitM attack", hostname)
		}

		if knownhosts.IsHostUnknown(khErr) {
			f, fErr := os.OpenFile(khPath, os.O_APPEND|os.O_WRONLY, 0o600)
			if fErr == nil {
				defer utils.CloseOrLogError(ctx)(f)

				fErr = knownhosts.WriteKnownHost(f, hostname, remote, key)
			}

			if fErr == nil {
				tflog.Info(ctx, fmt.Sprintf("Added host %s to known_hosts", hostname))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to add host %s to known_hosts", hostname), map[string]any{
					"error": khErr,
				})
			}

			return nil
		}

		return khErr
	})

	tflog.Info(ctx, fmt.Sprintf("agent is set to %t", c.agent))

	var sshClient *ssh.Client
	if c.agent {
		sshClient, err = c.createSSHClientAgent(ctx, cb, kh, sshHost)
		if err == nil {
			return sshClient, nil
		}

		tflog.Error(ctx, "Failed SSH connection through agent",
			map[string]any{
				"error": err,
			})
	}

	if c.privateKey != "" {
		sshClient, err = c.createSSHClientWithPrivateKey(ctx, cb, kh, sshHost)
		if err == nil {
			return sshClient, nil
		}

		tflog.Error(ctx, "Failed SSH connection with private key",
			map[string]any{
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
	kh *knownhosts.HostKeyDB,
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
	kh *knownhosts.HostKeyDB,
	sshHost string,
) (*ssh.Client, error) {
	conn, err := dialSocket(ctx, c.agentSocket)
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
	kh *knownhosts.HostKeyDB,
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

		tflog.Debug(ctx, "SSH connection via SOCKS5 established", map[string]any{
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

	tflog.Debug(ctx, "SSH connection established", map[string]any{
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

	sshConn, ch, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client connection: %w", err)
	}

	return ssh.NewClient(sshConn, ch, reqs), nil
}
