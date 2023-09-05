/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2023-2023. All rights reserved.
 */

package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccNetworkingSecGroupsDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_networking_secgroups.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.0.description",
						"[Acc Test] The security group created by Terraform."),
					resource.TestCheckResourceAttrPair(dataSourceName, "security_groups.0.id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.updated_at"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupsDataSource_description(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_networking_secgroups.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupsDataSource_description(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.0.description",
						"[Acc Test] The security group created by Terraform."),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.updated_at"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupsDataSource_id(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_networking_secgroups.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupsDataSource_id(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.0.description",
						"[Acc Test] The security group created by Terraform."),
					resource.TestCheckResourceAttrPair(dataSourceName, "security_groups.0.id",
						"hcs_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.0.updated_at"),
				),
			},
		},
	})
}

func testAccNetworkingSecGroupsDataSource_base(rName string) string {
	return fmt.Sprintf(`
resource "hcs_networking_secgroup" "test" {
  name        = "%s"
  description = "[Acc Test] The security group created by Terraform."
}
`, rName)
}

func testAccNetworkingSecGroupsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_secgroups" "test" {
  name = hcs_networking_secgroup.test.name
}
`, testAccNetworkingSecGroupsDataSource_base(rName))
}

func testAccNetworkingSecGroupsDataSource_description(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_secgroups" "test" {
  description = hcs_networking_secgroup.test.description
}
`, testAccNetworkingSecGroupsDataSource_base(rName))
}

func testAccNetworkingSecGroupsDataSource_id(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_secgroups" "test" {
  id = hcs_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupsDataSource_base(rName))
}
