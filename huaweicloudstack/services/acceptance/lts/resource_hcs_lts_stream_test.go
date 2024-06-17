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

func getLtsStreamResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups/{log_group_id}/streams"
	client, err := cfg.NewServiceClient("lts", acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_group_id}", state.Primary.Attributes["group_id"])

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
		return nil, fmt.Errorf("error parsing the log stream: %s", err)
	}

	streamId := state.Primary.ID
	streamResult := utils.PathSearch(fmt.Sprintf("log_streams|[?log_stream_id=='%s']|[0]", streamId), respBody, nil)
	if streamResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return streamResult, nil
}

func TestAccLtsStream_basic(t *testing.T) {
	var (
		stream       interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "hcs_lts_stream.test"
		rc           = acceptance.InitResourceCheck(resourceName, &stream, getLtsStreamResourceFunc)
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
				Config: testAccLtsStream_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "hcs_lts_group.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStreamImportState(resourceName),
			},
		},
	})
}

func TestAccLtsStream_ttl(t *testing.T) {
	var (
		stream       interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "hcs_lts_stream.test"
		rc           = acceptance.InitResourceCheck(resourceName, &stream, getLtsStreamResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsStream_ttl(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "5"),
					resource.TestCheckResourceAttr(resourceName, "filter_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "hcs_lts_group.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl_in_days"},
				ImportStateIdFunc:       testLtsStreamImportState(resourceName),
			},
		},
	})
}

func testLtsStreamImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		streamID := rs.Primary.ID
		groupID := rs.Primary.Attributes["group_id"]

		return fmt.Sprintf("%s/%s", groupID, streamID), nil
	}
}

func testAccLtsStream_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 5
}

resource "hcs_lts_stream" "test" {
  group_id              = hcs_lts_group.test.id
  stream_name           = "%[1]s"
  enterprise_project_id = "%[2]s"
}
`, rName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccLtsStream_ttl(rName string) string {
	return fmt.Sprintf(`
resource "hcs_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 5
}

resource "hcs_lts_stream" "test" {
  group_id    = hcs_lts_group.test.id
  stream_name = "%[1]s"
  ttl_in_days = 5
}
`, rName)
}
