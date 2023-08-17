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
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
				),
			},
			{
				Config: testAccVPCEPService_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-"+rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
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
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
				),
			},
			{
				Config: testAccVPCEPService_PermissionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func testAccVPCEPService_Precondition(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpc" "myvpc" {
  name = "vpc-default"
}

resource "hcs_compute_instance" "ecs" {
  name               = "%s"
  image_id           = data.hcs_images_image.test.id
  flavor_id          = data.hcs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccVPCEPService_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.hcs_vpc.myvpc.id
  port_id     = hcs_compute_instance.ecs.network[0].port
  approval    = false
  description = "test description"

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

func testAccVPCEPService_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpcep_service" "test" {
  name        = "tf-%s"
  server_type = "VM"
  vpc_id      = data.hcs_vpc.myvpc.id
  port_id     = hcs_compute_instance.ecs.network[0].port
  approval    = true
  description = "test description update"

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc-update"
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

func testAccVPCEPService_Permission(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.hcs_vpc.myvpc.id
  port_id     = hcs_compute_instance.ecs.network[0].port
  approval    = false
  permissions = ["iam:domain::1234", "iam:domain::5678"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

func testAccVPCEPService_PermissionUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.hcs_vpc.myvpc.id
  port_id     = hcs_compute_instance.ecs.network[0].port
  approval    = false
  permissions = ["iam:domain::abcd"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

const testAccCompute_data = `
data "hcs_availability_zones" "test" {}

data "hcs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcs_images_image" "test" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

data "hcs_networking_secgroup" "test" {
  name = "default"
}
`
