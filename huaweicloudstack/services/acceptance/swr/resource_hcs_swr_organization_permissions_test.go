package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/swr/v2/namespaces"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getResourcePermissions(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	return namespaces.GetAccess(swrClient, state.Primary.ID).Extract()
}

func TestAccSwrOrganizationPermissions_basic(t *testing.T) {
	var permissions namespaces.Access
	organizationName := acceptance.RandomAccResourceName()
	userName := acceptance.RandomAccResourceName()
	resourceName := "hcs_swr_organization_permissions.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&permissions,
		getResourcePermissions,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckUcs(t)
			acceptance.TestAccPreCheckSwrOrgPermissions(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccswrOrganizationPermissions_basic(organizationName, userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${hcs_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName+"_1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Read"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName+"_2"),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.2.user_name", userName+"_3"),
					resource.TestCheckResourceAttr(resourceName, "users.2.permission", "Manage"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccswrOrganizationPermissions_update(organizationName, userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${hcs_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName+"_1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName+"_2"),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Read"),
					resource.TestCheckResourceAttr(resourceName, "users.2.user_name", userName+"_4"),
					resource.TestCheckResourceAttr(resourceName, "users.2.permission", "Manage"),
					resource.TestCheckResourceAttr(resourceName, "users.3.user_name", userName+"_5"),
					resource.TestCheckResourceAttr(resourceName, "users.3.permission", "Read"),
				),
			},
		},
	})
}

func testAccswrOrganizationPermissions_basic(organizationName, userName string) string {
	return fmt.Sprintf(`
resource "hcs_swr_organization" "test" {
  name = "%[1]s"
}

resource "hcs_swr_organization_permissions" "test" {
  organization = hcs_swr_organization.test.name

  users {
    user_name  = "%[2]s"
    user_id    = "%[3]s"
    permission = "Read"
  }
  users {
    user_name  = "%[4]s"
    user_id    = "%[5]s"
    permission = "Write"
  }
  users {
    user_id    = "%[6]s"
    permission = "Manage"
  }
}
`, organizationName, acceptance.HCS_IAM_USER1_name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_name,
		acceptance.HCS_IAM_USER2_ID, acceptance.HCS_IAM_USER3_ID)
}

func testAccswrOrganizationPermissions_update(organizationName, userName string) string {
	return fmt.Sprintf(`
resource "hcs_swr_organization" "test" {
  name = "%[1]s"
}

resource "hcs_swr_organization_permissions" "test" {
  organization = hcs_swr_organization.test.name

  users {
    user_name  = "%[2]s"
    user_id    = "%[3]s"
    permission = "Manage"
  }

  users {
    user_id    = "%[4]s"
    permission = "Manage"
  }
}
`, organizationName, acceptance.HCS_IAM_USER1_name, acceptance.HCS_IAM_USER1_ID, acceptance.HCS_IAM_USER2_ID)
}
