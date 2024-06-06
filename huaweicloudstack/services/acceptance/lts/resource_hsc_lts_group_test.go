package lts

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

func getLtsGroupResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups"
	client, err := cfg.NewServiceClient("lts", acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing the log group: %s", err)
	}

	groupId := state.Primary.ID
	groupResult := utils.PathSearch(fmt.Sprintf("log_groups|[?log_group_id=='%s']|[0]", groupId), respBody, nil)
	if groupResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return groupResult, nil
}

func TestAccLtsGroup_basic(t *testing.T) {
	var (
		group        interface{}
		resourceName = "hcs_lts_group.test"
		rName        = acceptance.RandomAccResourceName()
		rc           = acceptance.InitResourceCheck(resourceName, &group, getLtsGroupResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsGroup_basic(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLtsGroup_basic(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
				),
			},
			{
				Config: testAccLtsGroup_tags(rName, 6),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "6"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testAccLtsGroup_basic(name string, ttl int) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d
}
`, name, ttl)
}

func testAccLtsGroup_tags(name string, ttl int) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, ttl)
}
