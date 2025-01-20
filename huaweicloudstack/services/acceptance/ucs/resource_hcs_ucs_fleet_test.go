package ucs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getFleetResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	// getFleet: Query the UCS Fleet detail
	var (
		region          = acceptance.HCS_REGION_NAME
		getFleetHttpUrl = "v1/clustergroups/{id}"
		getFleetProduct = "ucs"
	)
	getFleetClient, err := cfg.NewServiceClient(getFleetProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getFleetPath := getFleetClient.Endpoint + getFleetHttpUrl
	getFleetPath = strings.ReplaceAll(getFleetPath, "{id}", state.Primary.ID)

	getFleetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getFleetResp, err := getFleetClient.Request("GET", getFleetPath, &getFleetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Fleet: %s", err)
	}

	getFleetRespBody, err := utils.FlattenResponse(getFleetResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Fleet: %s", err)
	}

	return getFleetRespBody, nil
}

func TestAccFleet_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcs_ucs_fleet.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFleetResourceFunc,
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
				Config: testFleet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "permissions.0.namespaces.0", "default"),
					resource.TestCheckResourceAttrPair(rName, "permissions.0.policy_ids.0",
						"hcs_ucs_policy.test1", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testFleet_update_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "permissions.0.namespaces.0", "*"),
					resource.TestCheckResourceAttr(rName, "permissions.1.namespaces.0", "default"),
					resource.TestCheckResourceAttr(rName, "permissions.1.namespaces.1", "kube-system"),
					resource.TestCheckResourceAttrPair(rName, "permissions.0.policy_ids.0",
						"hcs_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.0",
						"hcs_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.1",
						"hcs_ucs_policy.test2", "id"),
				),
			},
			{
				Config: testFleet_update_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "permissions.#", "0"),
				),
			},
			{
				Config: testFleet_update_3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "permissions.0.namespaces.0", "*"),
					resource.TestCheckResourceAttr(rName, "permissions.1.namespaces.0", "kube-public"),
					resource.TestCheckResourceAttrPair(rName, "permissions.0.policy_ids.0",
						"hcs_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.0",
						"hcs_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.1",
						"hcs_ucs_policy.test2", "id"),
				),
			},
		},
	})
}

func testFleet_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test1" {
  name         = "%[1]s"
  iam_user_ids = [%[2]s]
  type         = "admin"
  description  = "created by terraform"
}

resource "hcs_ucs_fleet" "test" {
  name        = "%[1]s"
  description = "created by terraform"

  permissions {
    namespaces = ["default"]
    policy_ids = [hcs_ucs_policy.test1.id]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID)
}

func testFleet_update_1(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test1" {
  name         = "%[1]s-1"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "admin"
  description  = "created by terraform"
}

resource "hcs_ucs_policy" "test2" {
  name         = "%[1]s-2"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "custom"
  description  = "created by terraform"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}

resource "hcs_ucs_fleet" "test" {
  name        = "%[1]s"
  description = "created by terraform update"

  permissions {
    namespaces = ["*"]
    policy_ids = [hcs_ucs_policy.test1.id]
  }

  permissions {
    namespaces = ["default", "kube-system"]
    policy_ids = [hcs_ucs_policy.test1.id, hcs_ucs_policy.test2.id]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}

func testFleet_update_2(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_fleet" "test" {
  name = "%[1]s"
}
`, name)
}

func testFleet_update_3(name string) string {
	return fmt.Sprintf(`
resource "hcs_ucs_policy" "test1" {
  name         = "%[1]s-1"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "admin"
  description  = "created by terraform"
}

resource "hcs_ucs_policy" "test2" {
  name         = "%[1]s-2"
  iam_user_ids = [%[2]s, %[3]s]
  type         = "custom"
  description  = "created by terraform"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}

resource "hcs_ucs_fleet" "test" {
  name        = "%[1]s"
  description = "created by terraform update"

  permissions {
    namespaces = ["*"]
    policy_ids = [hcs_ucs_policy.test1.id]
  }

  permissions {
    namespaces = ["kube-public"]
    policy_ids = [hcs_ucs_policy.test1.id, hcs_ucs_policy.test2.id]
  }
}
`, name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}
