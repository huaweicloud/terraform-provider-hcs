package vdc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVdcUserGroupV3_basic(t *testing.T) {
	var userGroupModel group.GroupModel

	vdcId := "a18c2ce0-5379-4b34-8a12-eee47f5cfa89"
	rName := "user_group_name_" + acceptance.RandomAccResourceName()
	description := "user_group_description_" + acceptance.RandomAccResourceName()

	createParams := map[string]string{
		"vdcId":       vdcId,
		"name":        rName,
		"description": description,
	}

	newName := rName + "-Edit"
	newDes := description + "-Edit"
	upParams := map[string]string{
		"vdcId":       vdcId,
		"name":        newName,
		"description": newDes,
	}

	resourceName := "hcs_vdc_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcUserGroupV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVdcUserGroupV3Basic(createParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcUserGroupV3Exists(resourceName, &userGroupModel),
					resource.TestCheckResourceAttr(resourceName, "vdc_id", vdcId),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccVdcUserGroupV3Update(upParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcUserGroupV3Exists(resourceName, &userGroupModel),
					resource.TestCheckResourceAttr(resourceName, "vdc_id", vdcId),
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "description", newDes),
				),
			},
		},
	})
}

func testAccVdcUserGroupV3Update(params map[string]string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_group" "test" {
      vdc_id       = "%s"
	  name         = "%s"
	  description  = "%s"
	}
	`, params["vdcId"], params["name"], params["description"])
}

func testAccCheckVdcUserGroupV3Exists(n string, userGroupModel *group.GroupModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		userGroupClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc user group client: %s", err)
		}

		found, err := group.Get(userGroupClient, rs.Primary.ID).ToExtract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Vdc user group not found")
		}

		*userGroupModel = found

		return nil
	}
}

func testAccVdcUserGroupV3Basic(params map[string]string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_group" "test" {
      vdc_id       = "%s"
	  name         = "%s"
	  description  = "%s"
	}
	`, params["vdcId"], params["name"], params["description"])
}

func testAccCheckVdcUserGroupV3Destroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	userGroupClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vdc user group client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vdc_group" {
			continue
		}

		_, err := group.Get(userGroupClient, rs.Primary.ID).ToExtract()
		if err == nil {
			return fmtp.Errorf("Vdc user group still exists")
		}
	}

	return nil
}
