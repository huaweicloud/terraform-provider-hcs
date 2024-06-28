package swr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getSwrImageRetentionPolicyResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getSwrImageRetentionPolicy: Query SWR image retention policy
	var (
		getSwrImageRetentionPolicyHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}"
		getSwrImageRetentionPolicyProduct = "swr"
	)
	getSwrImageRetentionPolicyClient, err := cfg.NewServiceClient(getSwrImageRetentionPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <organization_name>/<repository_name>/<retention_id>")
	}
	organization := parts[0]
	repository := parts[1]
	retentionId := parts[2]

	getSwrImageRetentionPolicyPath := getSwrImageRetentionPolicyClient.Endpoint + getSwrImageRetentionPolicyHttpUrl
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{namespace}", organization)
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{repository}", repository)
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{retention_id}", retentionId)

	getSwrImageRetentionPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageRetentionPolicyResp, err := getSwrImageRetentionPolicyClient.Request("GET",
		getSwrImageRetentionPolicyPath, &getSwrImageRetentionPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR image retention policy: %s", err)
	}
	return utils.FlattenResponse(getSwrImageRetentionPolicyResp)
}

func TestAccSwrImageRetentionPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_swr_image_retention_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageRetentionPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageRetentionPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcs_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcs_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "type", "date_rule"),
					resource.TestCheckResourceAttr(rName, "number", "30"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.#", "3"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.pattern", "1.1"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.pattern", "1.2"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.2.kind", "regexp"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.2.pattern", "abc*"),
				),
			},
			{
				Config: testSwrImageRetentionPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcs_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcs_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "type", "date_rule"),
					resource.TestCheckResourceAttr(rName, "number", "25"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.pattern", "2.1"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.kind", "regexp"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.pattern", "xyz"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSwrImageRetentionPolicy_tag_rule(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_swr_image_retention_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageRetentionPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageRetentionPolicy_tag_rule(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcs_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcs_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "type", "tag_rule"),
					resource.TestCheckResourceAttr(rName, "number", "30"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.#", "3"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.pattern", "1.1"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.pattern", "1.2"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.2.kind", "regexp"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.2.pattern", "abc*"),
				),
			},
			{
				Config: testSwrImageRetentionPolicy_tag_rule_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"hcs_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"hcs_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "type", "tag_rule"),
					resource.TestCheckResourceAttr(rName, "number", "25"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.kind", "label"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.0.pattern", "2.1"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.kind", "regexp"),
					resource.TestCheckResourceAttr(rName, "tag_selectors.1.pattern", "xyz"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSwrImageRetentionPolicy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_swr_image_retention_policy" "test" {
  organization = hcs_swr_organization.test.name
  repository   = hcs_swr_repository.test.name
  type         = "date_rule"
  number       = 30

  tag_selectors {
    kind    = "label"
    pattern = "1.1"
  }
  tag_selectors {
    kind    = "label"
    pattern = "1.2"
  }
  tag_selectors {
    kind    = "regexp"
    pattern = "abc*"
  }
}
`, testAccSWRRepository_basic(name))
}

func testSwrImageRetentionPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_swr_image_retention_policy" "test" {
  organization = hcs_swr_organization.test.name
  repository   = hcs_swr_repository.test.name
  type         = "date_rule"
  number       = 25

  tag_selectors {
    kind    = "label"
    pattern = "2.1"
  }
  tag_selectors {
    kind    = "regexp"
    pattern = "xyz"
  }
}
`, testAccSWRRepository_basic(name))
}

func testSwrImageRetentionPolicy_tag_rule(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_swr_image_retention_policy" "test" {
  organization = hcs_swr_organization.test.name
  repository   = hcs_swr_repository.test.name
  type         = "tag_rule"
  number       = 30

  tag_selectors {
    kind    = "label"
    pattern = "1.1"
  }
  tag_selectors {
    kind    = "label"
    pattern = "1.2"
  }
  tag_selectors {
    kind    = "regexp"
    pattern = "abc*"
  }
}
`, testAccSWRRepository_basic(name))
}

func testSwrImageRetentionPolicy_tag_rule_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_swr_image_retention_policy" "test" {
  organization = hcs_swr_organization.test.name
  repository   = hcs_swr_repository.test.name
  type         = "tag_rule"
  number       = 25

  tag_selectors {
    kind    = "label"
    pattern = "2.1"
  }
  tag_selectors {
    kind    = "regexp"
    pattern = "xyz"
  }
}
`, testAccSWRRepository_basic(name))
}
