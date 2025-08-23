package storage

import "time"

// PBSStorageMutableFields defines the mutable attributes for 'pbs' type storage.
type PBSStorageMutableFields struct {
	DataStoreCommonMutableFields
	DataStoreWithBackups
	Fingerprint *string `json:"fingerprint,omitempty" url:"fingerprint,omitempty"`
	Encryption  *string `json:"encryption-key,omitempty" url:"encryption-key,omitempty"`
}

// PBSStorageImmutableFields defines the immutable attributes for 'pbs' type storage.
type PBSStorageImmutableFields struct {
	Username  *string `json:"username,omitempty" url:"username,omitempty"`
	Password  *string `json:"password,omitempty" url:"password,omitempty"`
	Namespace *string `json:"namespace,omitempty" url:"namespace,omitempty"`
	Server    *string `json:"server,omitempty" url:"server,omitempty"`
	Datastore *string `json:"datastore,omitempty" url:"datastore,omitempty"`
}

// PBSStorageCreateRequest defines the request body for creating a new PBS storage.
type PBSStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	PBSStorageMutableFields
	PBSStorageImmutableFields
}

// PBSStorageUpdateRequest defines the request body for updating an existing PBS storage.
type PBSStorageUpdateRequest struct {
	PBSStorageMutableFields
}

// EncryptionKey represents a Proxmox Backup Server encryption key object.
// Keys are stored as JSON and may include optional KDF (Key Derivation Function)
// parameters, creation/modification metadata, the key data itself, and its fingerprint.
//
// Example JSON:
//
//	{
//	  "kdf": { "Scrypt": { "n": 32768, "r": 8, "p": 1, "salt": "..." } },
//	  "created": "2025-08-18T15:04:05Z",
//	  "modified": "2025-08-18T15:04:05Z",
//	  "data": "base64-encoded-key",
//	  "fingerprint": "sha256:abcdef..."
//	}
type EncryptionKey struct {
	KDF         *KDF      `json:"kdf"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Data        string    `json:"data"`
	Fingerprint string    `json:"fingerprint"`
}

// KDF defines the Key Derivation Function configuration for an encryption key.
// Only one algorithm may be set at a time. If no KDF is used, this field is nil.
type KDF struct {
	Scrypt *ScryptParams `json:"Scrypt,omitempty"`
	PBKDF2 *PBKDF2Params `json:"PBKDF2,omitempty"`
}

// ScryptParams defines parameters for the scrypt key derivation function.
// The values control CPU/memory cost (N), block size (r), parallelization (p),
// and the random salt used to derive the key.
type ScryptParams struct {
	N    int    `json:"n"`
	R    int    `json:"r"`
	P    int    `json:"p"`
	Salt string `json:"salt"`
}

// PBKDF2Params defines parameters for the PBKDF2 key derivation function.
// It includes the iteration count (Iter) and the random salt value.
type PBKDF2Params struct {
	Iter int    `json:"iter"`
	Salt string `json:"salt"`
}
