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

func getaclRuleResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getDNSRecordset: Query DNS recordset
	getaclRuleClient, err := cfg.NewServiceClient("networkv2", region)
	if err != nil {
		return nil, fmt.Errorf("error creating network Client: %s", err)
	}
	getaclRuleHttpUrl := "v2.0/fwaas/firewall_rules/{firewall_rule_id}"

	getaclRulePath := getaclRuleClient.Endpoint + getaclRuleHttpUrl
	getaclRulePath = strings.ReplaceAll(getaclRulePath, "{firewall_rule_id}", state.Primary.ID)

	getaclRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := getaclRuleClient.Request("GET", getaclRulePath, &getaclRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS recordset: %s", err)
	}
	return utils.FlattenResponse(getDNSRecordsetResp)
}

func TestAccNetworkACLRule_basic(t *testing.T) {
	resourceKey := "hcs_network_acl_rule.rule_1"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceKey,
		&obj,
		getaclRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_basic_1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "deny"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "deny"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceKey, "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(resourceKey, "source_port", "444"),
					resource.TestCheckResourceAttr(resourceKey, "destination_port", "555"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "allow"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(resourceKey, "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(resourceKey, "source_port", "666"),
					resource.TestCheckResourceAttr(resourceKey, "destination_port", "777"),
				),
			},
			{
				ResourceName:      resourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkACLRule_anyProtocol(t *testing.T) {
	resourceKey := "hcs_network_acl_rule.rule_any"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var obj interface{}
	rc := acceptance.InitResourceCheck(
		resourceKey,
		&obj,
		getaclRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_anyProtocol(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "any"),
					resource.TestCheckResourceAttr(resourceKey, "action", "allow"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "192.168.199.0/24"),
				),
			},
		},
	})
}

func testAccNetworkACLRule_basic_1(rName string) string {
	return fmt.Sprintf(`
resource "hcs_network_acl_rule" "rule_1" {
  name = "%s"
  protocol = "udp"
  action = "deny"
}
`, rName)
}

func testAccNetworkACLRule_basic_2(rName string) string {
	return fmt.Sprintf(`
resource "hcs_network_acl_rule" "rule_1" {
  name = "%s"
  description = "Terraform accept test"
  protocol = "udp"
  action = "deny"
  source_ip_address = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_ports = ["444"]
  destination_ports = ["555"]
  enabled = true
}
`, rName)
}

func testAccNetworkACLRule_basic_3(rName string) string {
	return fmt.Sprintf(`
resource "hcs_network_acl_rule" "rule_1" {
  name = "%s"
  description = "Terraform accept test updated"
  protocol = "tcp"
  action = "allow"
  source_ip_address = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_ports = ["666"]
  destination_ports = ["777"]
  enabled = false
}
`, rName)
}

func testAccNetworkACLRule_anyProtocol(rName string) string {
	return fmt.Sprintf(`
resource "hcs_network_acl_rule" "rule_any" {
  name = "%s"
  description = "Allow any protocol"
  protocol = "any"
  action = "allow"
  source_ip_address = "192.168.199.0/24"
  enabled = true
}
`, rName)
}
