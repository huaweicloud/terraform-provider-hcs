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

func getRepositoryResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getRepository: Query the resource detail of the codeartsrepo repository
	var (
		getRepositoryHttpUrl = "v2/repositories/{repository_uuid}"
		getRepositoryProduct = "codeartsrepo"
	)
	getRepositoryClient, err := cfg.NewServiceClient(getRepositoryProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating repository client: %s", err)
	}

	getRepositoryPath := getRepositoryClient.Endpoint + getRepositoryHttpUrl
	getRepositoryPath = strings.ReplaceAll(getRepositoryPath, "{repository_uuid}", state.Primary.ID)

	getRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRepositoryResp, err := getRepositoryClient.Request("GET", getRepositoryPath, &getRepositoryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving codehub repository: %s", err)
	}
	return utils.FlattenResponse(getRepositoryResp)
}

func TestAccRepository_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_codearts_repository.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRepositoryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRepository_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "hcs_codearts_project.test",
						"id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform acc test"),
					resource.TestCheckResourceAttr(rName, "gitignore_id", "Go"),
					resource.TestCheckResourceAttr(rName, "enable_readme", "0"),
					resource.TestCheckResourceAttr(rName, "visibility_level", "20"),
					resource.TestCheckResourceAttr(rName, "license_id", "2"),
					resource.TestCheckResourceAttr(rName, "import_members", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"description",
					"gitignore_id",
					"enable_readme",
					"license_id",
					"import_members",
				},
			},
		},
	})
}

func TestAccRepository_default(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_codearts_repository.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRepositoryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRepository_default(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "hcs_codearts_project.test",
						"id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enable_readme", "1"),
					resource.TestCheckResourceAttr(rName, "visibility_level", "0"),
					resource.TestCheckResourceAttr(rName, "license_id", "1"),
					resource.TestCheckResourceAttr(rName, "import_members", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
					"description",
					"gitignore_id",
					"enable_readme",
					"license_id",
					"import_members",
				},
			},
		},
	})
}

func testRepository_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name = "%[1]s"
  type = "scrum"
}

resource "hcs_codearts_repository" "test" {
  project_id = hcs_codearts_project.test.id

  name             = "%[1]s"
  description      = "Created by terraform acc test"
  gitignore_id     = "Go"
  enable_readme    = 0  
  visibility_level = 20 
  license_id       = 2  
  import_members   = 0  
}
`, name)
}

func testRepository_default(name string) string {
	return fmt.Sprintf(`
resource "hcs_codearts_project" "test" {
  name = "%[1]s"
  type = "scrum"
}

resource "hcs_codearts_repository" "test" {
  project_id = hcs_codearts_project.test.id

  name = "%[1]s"
}
`, name)
}
