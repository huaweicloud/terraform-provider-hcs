/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2023-2023. All rights reserved.
 */

package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingSecGroupV3DataSource_byID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.hcs_networking_secgroup.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupV3DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.hcs_networking_secgroup.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func testAccNetworkingSecGroupV3DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "hcs_networking_secgroup" "test" {
  name                 = "%[1]s"
  description          = "My neutron security group"
  delete_default_rules = true
}

resource "hcs_networking_secgroup_rule" "ports" {
  security_group_id = hcs_networking_secgroup.test.id

  direction        = "ingress"
  action           = "allow"
  ethertype        = "IPv4"
  ports            = "80-100,8080"
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  priority         = 5
}

resource "hcs_networking_secgroup_rule" "port_range" {
  security_group_id = hcs_networking_secgroup.test.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 101
  port_range_max   = 200
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "hcs_networking_secgroup_rule" "remote_group" {
  security_group_id = hcs_networking_secgroup.test.id
  remote_group_id = hcs_networking_secgroup.test.id
  direction       = "ingress"
  action          = "allow"
  ethertype       = "IPv4"
}

resource "hcs_networking_secgroup_rule" "remote_address_group" {
  security_group_id = hcs_networking_secgroup.test.id

  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  ports                   = "8088"
  protocol                = "tcp"
  remote_group_id         = hcs_networking_secgroup.test.id
}  
`, rName)
}

func testAccNetworkingSecGroupV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_secgroup" "test" {
  name = hcs_networking_secgroup.test.name
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}

func testAccNetworkingSecGroupV3DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_secgroup" "test" {
  secgroup_id = hcs_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}
