package vdc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	RoleSDK "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"testing"
)

func TestAccVdcRoleResourceCreate(t *testing.T) {
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"

	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName),
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
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"

	rName := acceptance.RandomAccResourceName()

	rNameUpdate := rName + "_update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "XA"),
				),
			},
			{
				Config: testAccResourceVdcRoleUpdate(rNameUpdate),
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
	var role RoleSDK.VdcRoleModel
	resourceName := "hcs_vdc_role.test"
	resourceName2 := "hcs_vdc_role.test1"
	rName := acceptance.RandomAccResourceName()
	rName2 := rName + "1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVdcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVdcRoleCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "XA"),
				),
			},
			{
				Config: testAccResourceVdcRoleDelete(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVdcRoleExists(resourceName2, &role),
				),
			},
		},
	})
}

func testAccResourceVdcRoleCreate(rName string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_role" "test" {
		  domain_id = "82a610a5e8ac49bcbce7ca01e0092f4a"
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
		rName)
}

func testAccResourceVdcRoleUpdate(rName string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_role" "test" {
		  domain_id = "82a610a5e8ac49bcbce7ca01e0092f4a"
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
		rName)
}
func testAccResourceVdcRoleDelete(rName string) string {
	return fmt.Sprintf(`
		resource "hcs_vdc_role" "test1" {
		  domain_id = "82a610a5e8ac49bcbce7ca01e0092f4a"
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

func testAccCheckVdcDestroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	vdcClient, err := hcsConfig.VdcClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vdc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vdc_role" {
			continue
		}

		_, err := RoleSDK.Get(vdcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Vdc still exists")
		}
	}

	return nil
}
