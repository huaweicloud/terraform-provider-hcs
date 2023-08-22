package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vpcep/v1/services"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccVPCEPService_Basic(t *testing.T) {
	var service services.Service

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_vpcep_service.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getVpcepServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPService_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
				),
			},
			{
				Config: testAccVPCEPService_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCEPService_Permission(t *testing.T) {
	var service services.Service

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_vpcep_service.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getVpcepServiceResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPService_Permission(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
				),
			},
			{
				Config: testAccVPCEPService_PermissionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
				),
			},
		},
	})
}

func getVpcepServiceResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return services.Get(vpcepClient, state.Primary.ID).Extract()
}

func testAccVPCEPService_Basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID)
}

func testAccVPCEPService_Update(rName string) string {
	return fmt.Sprintf(`

resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = true

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID)
}

func testAccVPCEPService_Permission(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = false
  permissions = ["iam:domain::6e9dfd51d1124e8d8498dce894923a0d", "iam:domain::5fc973eea581490997e82ea11a1d0102"]

  port_mapping {
    service_port  = 11
    terminal_port = 22
  } 
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID)
}

func testAccVPCEPService_PermissionUpdate(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = false
  permissions = ["iam:domain::1e6dfd51d1124e8d8498dce123456a0e"]

  port_mapping {
    service_port  = 11
    terminal_port = 22
  }
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID)
}
