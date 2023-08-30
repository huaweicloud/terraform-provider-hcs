/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2023-2023. All rights reserved.
 */

package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getACLZoneResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}
	return firewall_groups.Get(c, state.Primary.ID).Extract()
}

func TestAccNetworkACL_basic(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "hcs_network_acl.fw_1"
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceKey,
		&obj,
		getACLZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "created by terraform test acc"),
					resource.TestCheckResourceAttr(resourceKey, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceKey, "inbound_policy_id"),
				),
			},
		},
	})
}

func TestAccNetworkACL_no_subnets(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "hcs_network_acl.fw_1"
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceKey,
		&obj,
		getACLZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "network acl without subents"),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
				),
			},
		},
	})
}

func TestAccNetworkACL_remove(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "hcs_network_acl.fw_1"
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceKey,
		&obj,
		getACLZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
				),
			},
		},
	})
}

func testAccNetworkACLRules(name string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "vpc_1" {
  name = "%s_vpc"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "subnet_1" {
  name = "%s_subnet_1"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = hcs_vpc.vpc_1.id
}

resource "hcs_vpc_subnet" "subnet_2" {
	name = "%s_subnet_2"
	cidr = "192.168.10.0/24"
	gateway_ip = "192.168.10.1"
	vpc_id = hcs_vpc.vpc_1.id
  }

resource "hcs_network_acl_rule" "rule_1" {
  name             = "%s-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
}

resource "hcs_network_acl_rule" "rule_2" {
  name             = "%s-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
}
`, name, name, name, name, name)
}

func testAccNetworkACL_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_network_acl" "fw_1" {
  name        = "%s"
  description = "created by terraform test acc"

  inbound_rules = [hcs_network_acl_rule.rule_1.id]
  subnets = [hcs_vpc_subnet.subnet_1.id]
}
`, testAccNetworkACLRules(name), name)
}

func testAccNetworkACL_no_subnets(name string) string {
	return fmt.Sprintf(`
resource "hcs_network_acl" "fw_1" {
  name        = "%s"
  description = "network acl without subents"
}
`, name)
}
