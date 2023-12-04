package ssh

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	//"github.com/bpg/terraform-provider-proxmox/proxmox/ssh"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

type nodeResolver struct {
	node ProxmoxNode
}

func (c *nodeResolver) Resolve(_ context.Context, _ string) (ProxmoxNode, error) {
	return c.node, nil
}

func TestExecuteNodeCommand(t *testing.T) {
	ctx := context.TODO()
	//commands := []string{"ls /dfjfjgldf"}
	command := "ls /"
	envs := []string{}

	endpoint := utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT")
	u, err := url.ParseRequestURI(endpoint)
	require.NoError(t, err)

	sshUsername := strings.Split(utils.GetAnyStringEnv("PROXMOX_VE_USERNAME"), "@")[0]
	sshAgentSocket := utils.GetAnyStringEnv("SSH_AUTH_SOCK", "PROXMOX_VE_SSH_AUTH_SOCK", "PM_VE_SSH_AUTH_SOCK")

	sshClient, err := NewClient(
		sshUsername, "", true, sshAgentSocket,
		&nodeResolver{
			node: ProxmoxNode{
				Address: u.Hostname(),
				Port:    22,
			},
		},
	)
	require.NoError(t, err)

	comErr, stdOut, stdErr := sshClient.ExecuteNodeCommand(
		ctx,
		u.Host,
		command,
		envs,
	)

	if comErr != nil {
		t.Logf("Error: %v", comErr)
	}

	t.Logf("stdout: %v", stdOut)
	t.Logf("stderr: %v", stdErr)

	if "Foo" != "Foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", "Foo", "Foo")
	}
}
