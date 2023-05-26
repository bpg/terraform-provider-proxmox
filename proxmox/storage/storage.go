/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// GetDatastore retrieves information about a datastore.
/*
Using undocumented API endpoints is not recommended, but sometimes it's the only way to get things done.
$ pvesh get /storage/local
┌─────────┬───────────────────────────────────────────┐
│ key     │ value                                     │
╞═════════╪═══════════════════════════════════════════╡
│ content │ images,vztmpl,iso,backup,snippets,rootdir │
├─────────┼───────────────────────────────────────────┤
│ digest  │ 5b65ede80f34631d6039e6922845cfa4abc956be  │
├─────────┼───────────────────────────────────────────┤
│ path    │ /var/lib/vz                               │
├─────────┼───────────────────────────────────────────┤
│ shared  │ 0                                         │
├─────────┼───────────────────────────────────────────┤
│ storage │ local                                     │
├─────────┼───────────────────────────────────────────┤
│ type    │ dir                                       │
└─────────┴───────────────────────────────────────────┘.
*/
func (c *Client) GetDatastore(
	ctx context.Context,
	datastoreID string,
) (*DatastoreGetResponseData, error) {
	resBody := &DatastoreGetResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("storage/%s", url.PathEscape(datastoreID)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving datastore %s: %w", datastoreID, err)
	}

	return resBody.Data, nil
}
