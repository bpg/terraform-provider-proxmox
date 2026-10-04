/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	proxmoxaccess "github.com/bpg/terraform-provider-proxmox/proxmox/access"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

type userTokenTestAPI struct {
	api.Client

	t      *testing.T
	method string
	calls  int
}

func (c *userTokenTestAPI) DoRequest(_ context.Context, method, path string, _, response any) error {
	require.Equal(c.t, c.method, method)
	require.Equal(c.t, "access/users/test@pve/token/test-token", path)
	c.calls++

	if method == http.MethodPost {
		body := response.(*proxmoxaccess.UserTokenCreateResponseBody)
		body.Data = &proxmoxaccess.UserTokenCreateResponseData{FullTokenID: "test@pve!test-token", Value: "test-value"}
	}

	return nil
}

func TestUserTokenStateDiagnostics(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		for _, invalidState := range []bool{false, true} {
			name := method + "/valid-state"
			if invalidState {
				name = method + "/invalid-state"
			}

			t.Run(name, func(t *testing.T) {
				ctx := t.Context()
				client := &userTokenTestAPI{t: t, method: method}
				r := &userTokenResource{client: proxmox.NewClient(client, nil, "")}
				var schemaResp resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

				model := userTokenModel{
					Comment: types.StringNull(), ExpirationDate: types.StringNull(),
					ID: types.StringValue("test@pve!test-token"), PrivSeparation: types.BoolValue(true),
					UserID: types.StringValue("test@pve"), TokenName: types.StringValue("test-token"),
					Value: types.StringValue("test@pve!test-token=test-value"),
				}
				plan := tfsdk.Plan{Schema: schemaResp.Schema}
				require.False(t, plan.Set(ctx, model).HasError())

				state := tfsdk.State{Schema: schemaResp.Schema}
				require.False(t, state.Set(ctx, model).HasError())
				output := state

				if invalidState {
					// Keep the request valid but make the response schema incompatible with the model.
					r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
					delete(schemaResp.Schema.Attributes, "value")
					output.Schema = schemaResp.Schema
				}

				if method == http.MethodPost {
					resp := resource.CreateResponse{State: output}
					r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
					require.Equal(t, invalidState, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

					if !invalidState {
						var actual userTokenModel
						require.False(t, resp.State.Get(ctx, &actual).HasError())
						require.Equal(t, model, actual)
					}
				} else {
					resp := resource.UpdateResponse{State: output}
					r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)
					require.Equal(t, invalidState, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

					if !invalidState {
						var actual userTokenModel
						require.False(t, resp.State.Get(ctx, &actual).HasError())
						require.Equal(t, model, actual)
					}
				}

				require.Equal(t, 1, client.calls)
			})
		}
	}
}
