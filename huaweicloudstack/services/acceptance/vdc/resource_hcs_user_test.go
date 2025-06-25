package vdc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/user"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceVdcUserV3_basic(t *testing.T) {
	var userModel user.VdcUserModel

	rName := acceptance.RandomAccResourceName()
	vdcId := "a18c2ce0-5379-4b34-8a12-eee47f5cfa89"
	pwd := ""
	displayName := "user_display_name" + acceptance.RandomAccResourceName()
	description := "user_description" + acceptance.RandomAccResourceName()
	var accessMode = [3]string{"default", "console", "programmatic"}
	createParams := map[string]string{
		"name":         rName,
		"vdcId":        vdcId,
		"password":     pwd,
		"display_name": displayName,
		"description":  description,
		"auth_type":    "LOCAL_AUTH",
		"access_mode":  accessMode[0],
	}

	upParams := map[string]string{
		"name":         rName,
		"vdcId":        vdcId,
		"password":     pwd + "Edit",
		"display_name": "user_display_name_edit" + acceptance.RandomAccResourceName(),
		"description":  "user_description_edit" + acceptance.RandomAccResourceName(),
		"auth_type":    "LOCAL_AUTH",
		"access_mode":  accessMode[1],
	}

	resourceName := "hcs_vdc_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcUserV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVdcUserV3Basic(createParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcUserV3Exists(resourceName, &userModel),
					resource.TestCheckResourceAttr(resourceName, "vdc_id", vdcId),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", createParams["display_name"]),
					resource.TestCheckResourceAttr(resourceName, "description", createParams["description"]),
					resource.TestCheckResourceAttr(resourceName, "auth_type", createParams["auth_type"]),
					resource.TestCheckResourceAttr(resourceName, "access_mode", createParams["access_mode"]),
				),
			},
			{
				Config: testAccVdcUserV3Update(upParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcUserV3Exists(resourceName, &userModel),
					resource.TestCheckResourceAttr(resourceName, "vdc_id", vdcId),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", upParams["display_name"]),
					resource.TestCheckResourceAttr(resourceName, "description", upParams["description"]),
					resource.TestCheckResourceAttr(resourceName, "auth_type", upParams["auth_type"]),
					resource.TestCheckResourceAttr(resourceName, "access_mode", upParams["access_mode"]),
				),
			},
		},
	})
}

func testAccVdcUserV3Update(params map[string]string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_user" "test" {
      vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}
	`, params["vdcId"], params["name"], params["password"], params["display_name"], params["description"], params["auth_type"], params["access_mode"])
}

func testAccCheckVdcUserV3Exists(n string, userModel *user.VdcUserModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		userClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc user client: %s", err)
		}

		found, err := user.Get(userClient, rs.Primary.ID).ToExtract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Vdc user not found")
		}

		*userModel = found

		return nil
	}
}

func testAccVdcUserV3Basic(params map[string]string) string {
	return fmt.Sprintf(`
	resource "hcs_vdc_user" "test" {
      vdc_id       = "%s"
	  name         = "%s"
	  password     = "%s"
	  display_name = "%s"
	  description  = "%s"
	  auth_type    = "%s"
	  access_mode  = "%s"
	}
	`, params["vdcId"], params["name"], params["password"], params["display_name"], params["description"], params["auth_type"], params["access_mode"])
}

func testAccCheckVdcUserV3Destroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	userClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vdc user client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vdc_user" {
			continue
		}

		_, err := user.Get(userClient, rs.Primary.ID).ToExtract()
		if err == nil {
			return fmtp.Errorf("Vdc user still exists")
		}
	}

	return nil
}
