package ucs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getPolicyResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	// getPolicy: Query the UCS Policy detail
	var (
		region           = acceptance.HCS_REGION_NAME
		getPolicyHttpUrl = "v1/permissions/rules"
		getPolicyProduct = "ucs"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl

	getPolicyResp, err := pagination.ListAllItems(
		getPolicyClient,
		"offset",
		getPolicyPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}

	getPolicyRespJson, err := json.Marshal(getPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}
	var getPolicyRespBody interface{}
	err = json.Unmarshal(getPolicyRespJson, &getPolicyRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Policy: %s", err)
	}

	jsonPath := fmt.Sprintf("items[?metadata.uid=='%s']|[0]", state.Primary.ID)
	getPolicyRespBody = utils.PathSearch(jsonPath, getPolicyRespBody, nil)
	if getPolicyRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getPolicyRespBody, nil
}

func TestAccPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcs_ucs_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUcs(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "admin"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "custom"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "details.0.operations.0", "*"),
					resource.TestCheckResourceAttr(rName, "details.0.resources.0", "*"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config:            testPolicy_basic_import(name),
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testPolicy_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "readonly"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "details.0.operations.0", "*"),
					resource.TestCheckResourceAttr(rName, "details.0.resources.0", "*"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPolicy_basic_update3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "develop"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "details.0.operations.0", "*"),
					resource.TestCheckResourceAttr(rName, "details.0.resources.0", "*"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
		},
	})
}

func testPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = [%[2]s]
  type         = "admin"
  description  = "created by terraform"
}
`, name, acceptance.HCS_IAM_USER1_ID)
}

func testPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "custom"
  description  = "created by terraform update"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}

func testPolicy_basic_import(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test" {
	name         = "%[1]s"
	iam_user_ids = [%[2]s, %[3]s]
	type         = "custom"
	description  = "created by terraform update"
	details {
	  operations = ["*"]
	  resources  = ["*"]
	}
  }
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}

func testPolicy_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "readonly"
  description  = "created by terraform update"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}

func testPolicy_basic_update3(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test" {
  name         = "%[1]s"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "develop"
  description  = "created by terraform update"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}
