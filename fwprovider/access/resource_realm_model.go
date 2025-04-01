/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
)

type realmResourceModel struct {
	Realm types.String `tfsdk:"realm"`
	Type  types.String `tfsdk:"type"`

	Acrvalues           types.String `tfsdk:"acr-values"`
	Autocreate          bool         `tfsdk:"autocreate"`
	Base_dn             types.String `tfsdk:"base_dn"`
	Bind_dn             types.String `tfsdk:"bind_dn"`
	Capath              types.String `tfsdk:"capath"`
	Casesensitive       bool         `tfsdk:"case-sensitive"`
	Cert                types.String `tfsdk:"cert"`
	Certkey             types.String `tfsdk:"certkey"`
	Checkconnection     bool         `tfsdk:"check-connection"`
	Clientid            types.String `tfsdk:"client-id"`
	Clientkey           types.String `tfsdk:"client-key"`
	Comment             types.String `tfsdk:"comment"`
	Default             bool         `tfsdk:"default"`
	Domain              types.String `tfsdk:"domain"`
	Filter              types.String `tfsdk:"filter"`
	Group_classes       types.String `tfsdk:"group_classes"`
	Group_dn            types.String `tfsdk:"group_dn"`
	Group_filter        types.String `tfsdk:"group_filter"`
	Group_name_attr     types.String `tfsdk:"group_name_attr"`
	Issuerurl           types.String `tfsdk:"issuer-url"`
	Mode                types.String `tfsdk:"mode"`
	Password            types.String `tfsdk:"password"`
	Port                int          `tfsdk:"port"`
	Prompt              types.String `tfsdk:"prompty"`
	Scopes              types.String `tfsdk:"scopes"`
	Secure              bool         `tfsdk:"secure"`
	Server1             types.String `tfsdk:"server1"`
	Server2             types.String `tfsdk:"server2"`
	Sslversion          types.String `tfsdk:"sslversion"`
	Syncdefaultsoptions types.String `tfsdk:"sync-defaults-options"`
	Sync_attributes     types.String `tfsdk:"sync_attributes"`
	Tfa                 types.String `tfsdk:"tfa"`
	User_attr           types.String `tfsdk:"user_attr"`
	User_classes        types.String `tfsdk:"user_classes"`
	Usernameclaim       types.String `tfsdk:"username-claim"`
	Verify              bool         `tfsdk:"verify"`
}

const promptFormat = "(?:none|login|consent|select_account|\\S+)"
const syncAttributesFormat = "\\w+=[^,]+(,\\s*\\w+=[^,]+)*"
const acrValuesFormat = "^[^\\x00-\\x1F\\x7F <>#\"]*$"
const domainFormat = "\\S+"

func parseValidateMode(id string) (*realmResourceModel, error) {
	/*
	 *	ldap | ldaps | ldap+starttls
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func parseValidateSslVersion(id string) (*realmResourceModel, error) {
	/*
	 *	tlsv1 | tlsv1_1 | tlsv1_2 | tlsv1_3
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func parseValidateSyncDefaultOptions(id string) (*realmResourceModel, error) {
	/*
	 *	[enable-new=<1|0>] [,full=<1|0>] [,purge=<1|0>] [,remove-vanished=([acl];[properties];[entry])|none] [,scope=<users|groups|both>]
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func parseValidatePort(id string) (*realmResourceModel, error) {
	/*
	 * (1 - 65535)
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func parseValidateTFA(id string) (*realmResourceModel, error) {
	/*
	 *	yubico | oath
	 * 	type=<TFATYPE> [,digits=<COUNT>] [,id=<ID>] [,key=<KEY>] [,step=<SECONDS>] [,url=<URL>]
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func parseValidateType(id string) (*realmResourceModel, error) {
	/*
	 *	ad | ldap | openid | pam | pve
	 *
	 */

	model := &realmResourceModel{}
	return model, nil
}

func (r *realmResourceModel) intoUpdateBody() *access.RealmUpdateRequestBody {
	body := &access.RealmUpdateRequestBody{
		// TODO
	}
	return body
}
