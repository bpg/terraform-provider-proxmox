/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

// func Test_HAGroupDatasource(t *testing.T) {
// 	t.Parallel()
//
// 	accProviders := testAccMuxProviders(context.Background(), t)
//
// 	datasourceName := "data.proxmox_virtual_environment_hagroups.all"
//
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: accProviders,
// 		Steps: []resource.TestStep{
// 			// Read testing
// 			{
// 				Config: `data "proxmox_virtual_environment_hagroups" "all" {}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(datasourceName, "release", "8.0"),
// 					resource.TestCheckResourceAttrSet(datasourceName, "id"),
// 				),
// 			},
// 		},
// 	})
// }
