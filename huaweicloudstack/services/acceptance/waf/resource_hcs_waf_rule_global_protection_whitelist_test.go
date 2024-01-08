package waf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getRuleGlobalProtectionWhitelistResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}"
		getProduct = "waf"
	)
	getClient, err := cfg.NewServiceClient(getProduct, acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getPath := getClient.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.Attributes["policy_id"])
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)

	queryParam := ""
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		queryParam = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	getPath += queryParam

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getResp, err := getClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RuleGlobalProtectionWhitelist: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccRuleGlobalProtectionWhitelist_basic(t *testing.T) {
	var obj interface{}

	rName := "hcs_waf_rule_global_protection_whitelist.test"
	randName := acceptance.RandomAccResourceName()

	rc := hwacceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleGlobalProtectionWhitelistResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRuleGlobalProtectionWhitelist_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "hcs_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "domains.0", "hcs_waf_dedicated_domain.domain_1", "domain"),
					resource.TestCheckResourceAttr(rName, "ignore_waf_protection", "xss;webshell"),
					resource.TestCheckResourceAttr(rName, "advanced_field", "params"),
					resource.TestCheckResourceAttr(rName, "advanced_content", "test_content"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "params"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "test content"),
					resource.TestCheckResourceAttr(rName, "conditions.0.subfield", "test_subfield"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "ip"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "192.168.0.1"),
					resource.TestCheckResourceAttr(rName, "conditions.2.field", "ip"),
					resource.TestCheckResourceAttr(rName, "conditions.2.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.2.content", "192.168.0.2"),
					resource.TestCheckResourceAttr(rName, "conditions.2.subfield", "x-forwarded-for"),
				),
			},
			{
				Config: testRuleGlobalProtectionWhitelist_basic_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domains.#", "0"),
					resource.TestCheckResourceAttr(rName, "ignore_waf_protection", "030004;030006;030007"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "header"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "not_contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "test header content"),
					resource.TestCheckResourceAttr(rName, "conditions.0.subfield", "custom_subfield"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "prefix"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "https://example.com"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func TestAccRuleGlobalProtectionWhitelist_withEpsID(t *testing.T) {
	var obj interface{}

	rName := "hcs_waf_rule_global_protection_whitelist.test"
	randName := acceptance.RandomAccResourceName()

	rc := hwacceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleGlobalProtectionWhitelistResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRuleGlobalProtectionWhitelist_basicWithEpsID(randName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "hcs_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "domains.0", "hcs_waf_dedicated_domain.domain_1", "domain"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "ignore_waf_protection", "xss;webshell"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "params"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "test content"),
					resource.TestCheckResourceAttr(rName, "conditions.0.subfield", "test_subfield"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func testRuleGlobalProtectionWhitelist_basic(randName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_rule_global_protection_whitelist" "test" {
  policy_id             = hcs_waf_policy.policy_1.id
  domains               = [hcs_waf_dedicated_domain.domain_1.domain]
  ignore_waf_protection = "xss;webshell"
  advanced_field        = "params"
  advanced_content      = "test_content"
  description           = "test description"
  status                = 0

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }

  conditions {
    field   = "ip"
    logic   = "equal"
    content = "192.168.0.1"
  }

  conditions {
    field    = "ip"
    logic    = "equal"
    content  = "192.168.0.2"
    subfield = "x-forwarded-for"
  }
}
`, testAccWafDedicatedDomainV1_policy(randName))
}

func testRuleGlobalProtectionWhitelist_basic_update(randName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_rule_global_protection_whitelist" "test" {
  policy_id             = hcs_waf_policy.policy_1.id
  domains               = []
  ignore_waf_protection = "030004;030006;030007"
  status                = 1

  conditions {
    field    = "header"
    logic    = "not_contain"
    content  = "test header content"
    subfield = "custom_subfield"
  }

  conditions {
    field   = "url"
    logic   = "prefix"
    content = "https://example.com"
  }
}
`, testAccWafDedicatedDomainV1_policy(randName))
}

func testRuleGlobalProtectionWhitelist_basicWithEpsID(randName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_rule_global_protection_whitelist" "test" {
  policy_id             = hcs_waf_policy.policy_1.id
  domains               = [hcs_waf_dedicated_domain.domain_1.domain]
  enterprise_project_id = "%s"
  ignore_waf_protection = "xss;webshell"

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }
}
`, testAccWafDedicatedDomainV1_policy_withEpsID(randName, epsID), epsID)
}
