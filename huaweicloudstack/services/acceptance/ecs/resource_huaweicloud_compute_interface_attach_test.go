package ecs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/compute/v2/extensions/attachinterfaces"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccComputeInterfaceAttach_Basic(t *testing.T) {
	var ai attachinterfaces.Interface
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_interface_attach.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterfaceAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInterfaceAttachExists(resourceName, &ai),
					testAccCheckComputeInterfaceAttachIP(&ai, "192.168.0.199"),
					resource.TestCheckResourceAttr(resourceName, "source_dest_check", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"hcs_networking_secgroup.test", "id"),
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

func computeInterfaceAttachParseID(id string) (instanceID, portID string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		err = fmt.Errorf("unable to parse the resource ID, must be <instance_id>/<port_id> format")
		return
	}

	instanceID = idParts[0]
	portID = idParts[1]
	return
}

func testAccCheckComputeInterfaceAttachDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.HcsConfig)
	computeClient, err := cfg.ComputeV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_interface_attach" {
			continue
		}

		instanceId, portId, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err == nil {
			return fmt.Errorf("interface attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeInterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.HcsConfig)
		computeClient, err := cfg.ComputeV2Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		instanceId, portId, err := computeInterfaceAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err != nil {
			return err
		}
		if found.PortID != portId {
			return fmt.Errorf("interface attachment not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeInterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("requested ip (%s) does not exist on port", ip)
	}
}

func testAccComputeInterfaceAttach_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  vpc_id     = hcs_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(hcs_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(hcs_vpc.test.cidr, 4, 0), 1)
}

resource "hcs_networking_secgroup" "test" {
  name = "%[1]s"
}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_ims_images" "test" {
  name       = "ecs_mini_image"
}

resource "hcs_ecs_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}

resource "hcs_ecs_compute_interface_attach" "test" {
  instance_id        = hcs_ecs_compute_instance.test.id
  network_id         = hcs_vpc_subnet.test.id
  fixed_ip           = cidrhost(cidrsubnet(hcs_vpc.test.cidr, 4, 0), 199)
  security_group_ids = [hcs_networking_secgroup.test.id]
}
`, rName)
}
