package vdc

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group_role_assignment"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"strconv"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVdcGroupRoleAssignmentV3_basic(t *testing.T) {

	groupId := ""
	roleId := ""
	epsId := ""
	projectId := ""

	resourceName := "hcs_vdc_group_role_assignment.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcGroupRoleAssignmentV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVdcGroupRoleAssignmentBasic(groupId, roleId, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcGroupRoleAssignmentV3Exists(resourceName, groupId),
				),
			},
			{
				Config: testAccVdcGroupRoleAssignmentV3Update(groupId, roleId, epsId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcGroupRoleAssignmentV3Exists(resourceName, groupId),
				),
			},
		},
	})
}

func testAccVdcGroupRoleAssignmentV3Update(groupId string, roleId string, epsId string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_group_role_assignment" "test" {
		group_id = "%s"
		role_assignment {
			 role_id = "%s"
			enterprise_project_id = "%s"
		}
	}
	`, groupId, roleId, epsId)
}

func testAccCheckVdcGroupRoleAssignmentV3Exists(n string, groupId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		var oldRoleIs []string
		if val, err := strconv.Atoi(rs.Primary.Attributes["role_assignment.#"]); err == nil && val == 1 {
			if u, ok := rs.Primary.Attributes["role_assignment.0.role_id"]; ok {
				oldRoleIs = append(oldRoleIs, u)
			} else {
				return fmtp.Errorf("No roles in result")
			}
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		client, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		allUsers, err := group_role_assignment.GetVdcGroupRoleAssignmentAllRoles(client, groupId)
		if err != nil {
			return err
		}

		for _, id := range oldRoleIs {
			filter := map[string]interface{}{
				"ID": id,
			}
			_, err1 := utils.FilterSliceWithZeroField(allUsers, filter)
			if err1 != nil {
				return fmtp.Errorf("Error to find role from all vdc group roles: %s", err1)
			}
		}

		return nil
	}
}

func testAccVdcGroupRoleAssignmentBasic(groupId string, roleId string, projectId string) string {
	return fmt.Sprintf(`
 
	resource "hcs_vdc_group_role_assignment" "test" {
		group_id = "%s"
		role_assignment {
			role_id = "%s"
			project_id = "%s"
		}
	}
	`, groupId, roleId, projectId)
}

func testAccCheckVdcGroupRoleAssignmentV3Destroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	client, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vdc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vdc_group_role_assignment" {
			continue
		}

		roles, err1 := group_role_assignment.GetVdcGroupRoleAssignmentAllRoles(client, rs.Primary.ID)
		if err1 != nil {
			return fmtp.Errorf("Error to retrive all group roles: %s", err1)
		}

		if len(roles) != 0 {
			return fmtp.Errorf("Failed to delete data. Data still exists in the group: %s", err1)
		}
	}

	return nil
}
