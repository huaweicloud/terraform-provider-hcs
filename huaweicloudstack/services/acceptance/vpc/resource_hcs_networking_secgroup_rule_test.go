/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2023-2023. All rights reserved.
 */

package vpc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getSecRuleResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getDNSRecordset: Query DNS recordset
	getSecRuleClient, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return nil, fmt.Errorf("error creating network Client: %s", err)
	}
	getSecRuleHttpUrl := fmt.Sprintf("v1/%s/security-group-rules/{rules_security_groups_id}", acceptance.HCS_PROJECT_ID)

	getSecRulePath := getSecRuleClient.Endpoint + getSecRuleHttpUrl
	getSecRulePath = strings.ReplaceAll(getSecRulePath, "{rules_security_groups_id}", state.Primary.ID)

	getSecRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := getSecRuleClient.Request("GET", getSecRulePath, &getSecRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS recordset: %s", err)
	}
	return utils.FlattenResponse(getDNSRecordsetResp)
}

func TestAccNetworkingSecGroupRule_basic(t *testing.T) {
	var resourceRuleName string = "hcs_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceRuleName,
		&obj,
		getSecRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "description", "This is a basic acc test"),
					resource.TestCheckResourceAttr(resourceRuleName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_oldPorts(t *testing.T) {
	var resourceRuleName string = "hcs_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceRuleName,
		&obj,
		getSecRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_oldPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_remoteGroup(t *testing.T) {
	var resourceRuleName string = "hcs_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceRuleName,
		&obj,
		getSecRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_remoteGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceRuleName, "remote_group_id"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_lowerCaseCIDR(t *testing.T) {
	var resourceRuleName string = "hcs_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceRuleName,
		&obj,
		getSecRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_lowerCaseCIDR(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "2001:558:fc00::/39"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_noPorts(t *testing.T) {
	var resourceRuleName string = "hcs_networking_secgroup_rule.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceRuleName,
		&obj,
		getSecRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_noPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
		},
	})
}

func testAccNetworkingSecGroupRule_base(rName string) string {
	return fmt.Sprintf(`
resource "hcs_networking_secgroup" "secgroup_test" {
 name        = "%s-secgroup"
 description = "terraform security group rule acceptance test"
}
`, rName)
}

func testAccNetworkingSecGroupRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup_rule" "secgroup_rule_test" {
 direction         = "ingress"
 description       = "This is a basic acc test"
 ethertype         = "IPv4"
 ports             = 80
 protocol          = "tcp"
 remote_ip_prefix  = "0.0.0.0/0"
 security_group_id = hcs_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_oldPorts(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup_rule" "secgroup_rule_test" {
 direction         = "ingress"
 ethertype         = "IPv4"
 port_range_min    = 80
 port_range_max    = 80
 protocol          = "tcp"
 remote_ip_prefix  = "0.0.0.0/0"
 security_group_id = hcs_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_remoteGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup_rule" "secgroup_rule_test" {
 direction         = "ingress"
 ethertype         = "IPv4"
 ports             = 80
 protocol          = "tcp"
 remote_group_id   = hcs_networking_secgroup.secgroup_test.id
 security_group_id = hcs_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_lowerCaseCIDR(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup_rule" "secgroup_rule_test" {
 direction         = "ingress"
 ethertype         = "IPv6"
 ports             = 80
 protocol          = "tcp"
 remote_ip_prefix  = "2001:558:FC00::/39"
 security_group_id = hcs_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_noPorts(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup_rule" "test" {
 security_group_id = hcs_networking_secgroup.secgroup_test.id
 direction         = "ingress"
 ethertype         = "IPv4"
 protocol          = "icmp"
 remote_ip_prefix  = "0.0.0.0/0"
}
`, testAccNetworkingSecGroupRule_base(rName))
}
