package codearts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getProjectResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getProject: Query the Project
	var (
		getProjectHttpUrl = "v4/projects/{id}"
		getProjectProduct = "projectman"
	)
	getProjectClient, err := cfg.NewServiceClient(getProjectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Project Client: %s", err)
	}

	getProjectPath := getProjectClient.Endpoint + getProjectHttpUrl
	getProjectPath = strings.ReplaceAll(getProjectPath, "{id}", state.Primary.ID)

	getProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProjectResp, err := getProjectClient.Request("GET", getProjectPath, &getProjectOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Project: %s", err)
	}
	return utils.FlattenResponse(getProjectResp)
}

func TestAccProject_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_codearts_project.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCodeartsReqTempalteId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProject_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "scrum"),
				),
			},
			{
				Config: testProject_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "xboard"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
				),
			},
			{
				Config: testProject_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "basic"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description2"),
				),
			},
			{
				Config: testProject_basic_update3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "phoenix"),
				),
			},
			{
				Config: testProject_basic_update4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ipd"),
					resource.TestCheckResourceAttr(rName, "template_id", acceptance.HCS_CODEARTSREQ_TEMPLATE_ID),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
		},
	})
}

func testProject_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name = "%s"
  type = "scrum"
}
`, name)
}

func testProject_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name        = "%s"
  type        = "xboard"
  description = "demo_description1"
}
`, name)
}

func testProject_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name        = "%s"
  type        = "basic"
  description = "demo_description2"
}
`, name)
}

func testProject_basic_update3(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name = "%s"
  type = "phoenix"
}
`, name)
}

func testProject_basic_update4(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name = "%s"
  type = "ipd"

  template_id = "%s"
}
`, name, acceptance.HCS_CODEARTSREQ_TEMPLATE_ID)
}
