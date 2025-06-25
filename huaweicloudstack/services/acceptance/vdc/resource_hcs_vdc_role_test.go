package vdc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/fileconfig"
	RoleSDK "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"testing"
)

func initConfig() *fileconfig.ConfigForAcceptance {
	var cfg, err = fileconfig.GetTestConfig()

	if err != nil {
		fmt.Printf("failed to get configuration.")
		return nil
	}
	return cfg
}

func TestAccVdcRoleResourceCreate(t *testing.T) {
	var testConfigEnv *fileconfig.ConfigForAcceptance
	testConfigEnv = initConfig()
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"

	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName, testConfigEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "XA"),
				),
			},
		},
	})
}

func TestAccVdcRoleResourceUpdate(t *testing.T) {

	var testConfigEnv *fileconfig.ConfigForAcceptance
	testConfigEnv = initConfig()
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"

	rName := acceptance.RandomAccResourceName()

	rNameUpdate := rName + "_update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName, testConfigEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "XA"),
				),
			},
			{
				Config: testAccResourceVdcRoleUpdate(rNameUpdate, testConfigEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "type", "AX"),
				),
			},
		},
	})
}

func TestAccVdcRoleResourceDelete(t *testing.T) {

	var testConfigEnv *fileconfig.ConfigForAcceptance
	testConfigEnv = initConfig()
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName, testConfigEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
				),
			},
			{
				Config:  " ",
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcDestroy(resourceName),
				),
			},
		},
	})
}

func testAccResourceVdcRoleCreate(rName string, testConfigEnv *fileconfig.ConfigForAcceptance) string {

	return fmt.Sprintf(`
		resource "hcs_vdc_role" "test" {
		  domain_id = "%s"
		  name = "%s"
		  type = "XA"
		  policy = <<EOF
		{
		  "Depends": [],
		  "Statement": [
			{
			  "Action": [
				"ecs:cloudServers:start",
				"ecs:cloudServers:list"
			  ],
			  "Effect": "Allow"
			}
		  ],
		  "Version": "1.1"
		}
		EOF
		}`,
		testConfigEnv.NewDomainId,
		rName)
}

func testAccResourceVdcRoleUpdate(rName string, testConfigEnv *fileconfig.ConfigForAcceptance) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_role" "test" {
		  domain_id = "%s"
		  name = "%s"
		  type = "AX"
		  policy = <<EOF
		{
		  "Depends": [],
		  "Statement": [
			{
			  "Action": [
				"ecs:cloudServers:start",
				"ecs:cloudServers:list"
			  ],
			  "Effect": "Allow"
			}
		  ],
		  "Version": "1.1"
		}
		EOF
		}`,
		testConfigEnv.NewDomainId,
		rName)
}
func testAccCheckVdcRoleExists(n string, role *RoleSDK.VdcRoleModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		found, err := RoleSDK.Get(vdcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("vdc not found")
		}

		*role = *found

		return nil
	}
}

func testAccCheckVdcDestroy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())

		vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vdc client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceName {
				continue
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("resource ID is empty")
			}
			role, err := RoleSDK.Get(vdcClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmtp.Errorf("Vdc still exists")
			}

			if role.ID != rs.Primary.ID {
				return fmtp.Errorf("resource ID does not match: expected=%s, got=%s", rs.Primary.ID, role.ID)
			}
		}

		return nil
	}

}
