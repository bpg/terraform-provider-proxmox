package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PBSStorageModel maps the Terraform schema for PBS storage.
type PBSStorageModel struct {
	StorageModelBase
	Server                   types.String `tfsdk:"server"`
	Datastore                types.String `tfsdk:"datastore"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	Namespace                types.String `tfsdk:"namespace"`
	Fingerprint              types.String `tfsdk:"fingerprint"`
	EncryptionKey            types.String `tfsdk:"encryption_key"`
	EncryptionKeyFingerprint types.String `tfsdk:"encryption_key_fingerprint"`
	GenerateEncryptionKey    types.Bool   `tfsdk:"generate_encryption_key"`
	GeneratedEncryptionKey   types.String `tfsdk:"generated_encryption_key"`
}

// GetStorageType returns the storage type identifier.
func (m *PBSStorageModel) GetStorageType() types.String {
	return types.StringValue("pbs")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *PBSStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.PBSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.PBSStorageMutableFields.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Username = m.Username.ValueStringPointer()
	request.Password = m.Password.ValueStringPointer()
	request.Namespace = m.Namespace.ValueStringPointer()
	request.Server = m.Server.ValueStringPointer()
	request.Datastore = m.Datastore.ValueStringPointer()

	request.Fingerprint = m.Fingerprint.ValueStringPointer()

	if !m.GenerateEncryptionKey.IsNull() && m.GenerateEncryptionKey.ValueBool() {
		request.Encryption = types.StringValue("autogen").ValueStringPointer()
	} else if !m.EncryptionKey.IsNull() && m.EncryptionKey.ValueString() != "" {
		request.Encryption = m.EncryptionKey.ValueStringPointer()
	}

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *PBSStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.PBSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Fingerprint = m.Fingerprint.ValueStringPointer()

	if !m.GenerateEncryptionKey.IsNull() && m.GenerateEncryptionKey.ValueBool() {
		request.Encryption = types.StringValue("autogen").ValueStringPointer()
	} else if !m.EncryptionKey.IsNull() && m.EncryptionKey.ValueString() != "" {
		request.Encryption = m.EncryptionKey.ValueStringPointer()
	}

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
// Password is not returned by the API so we leave it as is in the state.
func (m *PBSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.Server != nil {
		m.Server = types.StringValue(*datastore.Server)
	}
	if datastore.Datastore != nil {
		m.Datastore = types.StringValue(*datastore.Datastore)
	}
	if datastore.Username != nil {
		m.Username = types.StringValue(*datastore.Username)
	}
	if datastore.Namespace != nil {
		m.Namespace = types.StringValue(*datastore.Namespace)
	}
	if datastore.Fingerprint != nil {
		m.Fingerprint = types.StringValue(*datastore.Fingerprint)
	}

	return nil
}
