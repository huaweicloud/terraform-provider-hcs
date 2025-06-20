package vdc

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group_membership"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"strconv"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVdcGroupMembershipV3_basic(t *testing.T) {

	groupId := ""
	rName := acceptance.RandomAccResourceName()
	rName1 := acceptance.RandomAccResourceName()
	rName2 := acceptance.RandomAccResourceName()
	rName3 := acceptance.RandomAccResourceName()
	vdcId := ""
	pd := ""
	displayName := "user_display_name" + acceptance.RandomAccResourceName()
	displayName1 := "user_display_name" + acceptance.RandomAccResourceName()
	displayName2 := "user_display_name" + acceptance.RandomAccResourceName()
	displayName3 := "user_display_name" + acceptance.RandomAccResourceName()
	description := "user_description" + acceptance.RandomAccResourceName()
	authType := "LOCAL_AUTH"
	accessMode := "default"

	resourceName := "hcs_vdc_group_membership.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcGroupMembershipV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVdcGroupMembershipBasic(rName, rName1, vdcId, pd, displayName, displayName1, description, authType, accessMode, groupId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcGroupMembershipV3Exists(resourceName, groupId),
				),
			},
			{
				Config: testAccVdcGroupMembershipV3Update(rName2, rName3, vdcId, pd, displayName2, displayName3, description, authType, accessMode, groupId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcGroupMembershipV3Exists(resourceName, groupId),
				),
			},
		},
	})
}

func testAccVdcGroupMembershipV3Update(rName string, rName1 string, vdcId string, pd string, displayName string, displayName1 string, description string, authType string, accessMode string, groupId string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_user" "test3" {
	  vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}

	resource "hcs_vdc_user" "test4" {
	  vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}

	resource "hcs_vdc_group_membership" "test" {
		group = "%s"
		users = [hcs_vdc_user.test3.id, hcs_vdc_user.test4.id]
	}
	`, vdcId, rName, pd, displayName, description, authType, accessMode, vdcId, rName1, pd, displayName1, description, authType, accessMode, groupId)
}

func testAccCheckVdcGroupMembershipV3Exists(n string, groupId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		var oldUserIds []string
		if val, err := strconv.Atoi(rs.Primary.Attributes["users.#"]); err == nil && val == 2 {
			if u, ok := rs.Primary.Attributes["users.0"]; ok {
				oldUserIds = append(oldUserIds, u)
			} else {
				return fmtp.Errorf("No users in result")
			}
			if u, ok := rs.Primary.Attributes["users.1"]; ok {
				oldUserIds = append(oldUserIds, u)
			} else {
				return fmtp.Errorf("No users in result")
			}
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		client, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		allUsers, err := group_membership.GetGroupMemberShipAllUser(client, groupId)
		if err != nil {
			return err
		}

		for _, id := range oldUserIds {
			filter := map[string]interface{}{
				"ID": id,
			}
			_, err1 := utils.FilterSliceWithZeroField(allUsers, filter)
			if err1 != nil {
				return fmtp.Errorf("error to find user from all group users: %s", err1)
			}
		}

		return nil
	}
}

func testAccVdcGroupMembershipBasic(rName string, rName1 string, vdcId string, pd string, displayName string, displayName1 string, description string, authType string, accessMode string, groupId string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_user" "test1" {
      vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}

	resource "hcs_vdc_user" "test2" {
      vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}

	resource "hcs_vdc_group_membership" "test" {
		group = "%s"
		users = [hcs_vdc_user.test1.id, hcs_vdc_user.test2.id]
	}
	`, vdcId, rName, pd, displayName, description, authType, accessMode, vdcId, rName1, pd, displayName1, description, authType, accessMode, groupId)
}

func testAccCheckVdcGroupMembershipV3Destroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	client, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vdc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vdc_group_membership" {
			continue
		}

		_, err1 := group_membership.GetGroupMemberShipAllUser(client, rs.Primary.ID)
		if err1 != nil {
			return fmtp.Errorf("error to retrive all group users: %s", err1)
		}
	}

	return nil
}
