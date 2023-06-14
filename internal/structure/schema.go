/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package structure

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// MergeSchema merges the schema.Schema from src into dst. Safety
// against conflicts is enforced by panicking.
func MergeSchema(dst, src schema.Schema) {
	for a, v := range src.GetAttributes() {
		if _, ok := dst.Attributes[a]; ok {
			panic(fmt.Errorf("conflicting schema attribute: %s", a))
		}

		dst.Attributes[a] = v
	}

	for b, v := range src.GetBlocks() {
		if _, ok := dst.Blocks[b]; ok {
			panic(fmt.Errorf("conflicting schema block: %s", b))
		}

		dst.Blocks[b] = v
	}
}
