package applier

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func getTestClient() (*Client, error) {
	apiToken := os.Getenv("PVE_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("PVE_TOKEN env variable not set")
	}

	conURL := os.Getenv("PVE_URL")
	if conURL == "" {
		return nil, fmt.Errorf("PVE_URL env variable not set")
	}

	conn, err := api.NewConnection(conURL, true, "")
	if err != nil {
		return nil, err
	}

	creds := api.Credentials{
		TokenCredentials: &api.TokenCredentials{
			APIToken: apiToken,
		},
	}

	client, err := api.NewClient(creds, conn)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}

func TestApplyConfig(t *testing.T) {
	t.Parallel()

	client, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := client.ApplyConfig(ctx); err != nil {
		t.Fatalf("ApplyConfig failed: %v", err)
	}

	t.Logf("SDN configuration applied successfully")
}
