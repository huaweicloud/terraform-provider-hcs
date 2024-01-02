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

func getAddressGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getWAFAddressGroup: Query WAF address group
	var (
		getWAFAddressGroupHttpUrl = "v1/{project_id}/waf/ip-group/{id}"
		getWAFAddressGroupProduct = "waf"
	)
	getWAFAddressGroupClient, err := cfg.NewServiceClient(getWAFAddressGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getWAFAddressGroupPath := getWAFAddressGroupClient.Endpoint + getWAFAddressGroupHttpUrl
	getWAFAddressGroupPath = strings.ReplaceAll(getWAFAddressGroupPath, "{project_id}",
		getWAFAddressGroupClient.ProjectID)
	getWAFAddressGroupPath = strings.ReplaceAll(getWAFAddressGroupPath, "{id}", state.Primary.ID)

	enterpriseProjectID := state.Primary.Attributes["enterprise_project_id"]
	if enterpriseProjectID != "" {
		getWAFAddressGroupPath += fmt.Sprintf("?enterprise_project_id=%s", enterpriseProjectID)
	}

	getWAFAddressGroupOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getWAFAddressGroupResp, err := getWAFAddressGroupClient.Request("GET", getWAFAddressGroupPath,
		&getWAFAddressGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving address group: %s", err)
	}
	return utils.FlattenResponse(getWAFAddressGroupResp)
}

func TestAccAddressGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_waf_address_group.test"

	rc := hwacceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "example_description"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(rName, "rules.#"),
				),
			},
			{
				Config: testAddressGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "example_description_update"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0", "192.168.1.0"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.1", "192.168.2.0/12"),
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

func TestAccAddressGroup_withEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_waf_address_group.test"

	rc := hwacceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroup_withEpsId(name, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "example_description"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(rName, "rules.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(rName),
			},
		},
	})
}

func testAddressGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_address_group" "test" {
  name         = "%s"
  description  = "example_description"
  ip_addresses = ["192.168.1.0/24"]

  depends_on   = [hcs_waf_dedicated_instance.instance_1]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAddressGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_address_group" "test" {
  name         = "%s_update"
  description  = "example_description_update"
  ip_addresses = ["192.168.1.0", "192.168.2.0/12"]

  depends_on   = [hcs_waf_dedicated_instance.instance_1]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAddressGroup_withEpsId(name, enterpriseProjectID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_waf_address_group" "test" {
  name                  = "%s"
  description           = "example_description"
  ip_addresses          = ["192.168.1.0/24"]
  enterprise_project_id = "%s"

  depends_on = [hcs_waf_dedicated_instance.instance_1]
}
`, testAccWafDedicatedInstance_epsId(name, enterpriseProjectID), name, enterpriseProjectID)
}

// testWAFResourceImportState use to return an id with format <id> or <id>/<enterprise_project_id>
func testWAFResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		epsID := rs.Primary.Attributes["enterprise_project_id"]
		if epsID == "" {
			return rs.Primary.ID, nil
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, epsID), nil
	}
}
