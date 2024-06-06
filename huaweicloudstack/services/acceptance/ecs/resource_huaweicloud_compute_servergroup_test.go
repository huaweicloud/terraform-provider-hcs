package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/servergroups"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_server_group.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccComputeServerGroup_scheduler(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_server_group.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_scheduler(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("hcs_ecs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_members(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_server_group.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_members(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("hcs_ecs_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_concurrency(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_concurrency(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("members_attached", "true"),
					resource.TestCheckOutput("volumes_attached", "true"),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	ecsClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_server_group" {
			continue
		}

		_, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("server group still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		ecsClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("server group not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeInstanceInServerGroup(instance *cloudservers.CloudServer, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("instance %s does not belong to server group %s", instance.ID, sg.ID)
	}
}

func testAccComputeServerGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_ecs_compute_server_group" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}
`, rName)
}

func testAccComputeServerGroup_scheduler(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_server_group" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}

resource "hcs_ecs_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.id
  flavor_id          = data.hcs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  scheduler_hints {
    group = hcs_ecs_compute_server_group.sg_1.id
  }
  network {
    uuid = data.hcs_vpc_subnets.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeServerGroup_members(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_server_group" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
  members  = [hcs_ecs_compute_instance.instance_1.id]
}

resource "hcs_ecs_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnets.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeServerGroup_concurrency(name string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_ims_images" "test" {
  flavor_id  = data.hcs_ecs_compute_flavors.test.ids[0]
  os         = "Ubuntu"
  visibility = "public"
}

resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.192.0/20"
}

resource "hcs_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = hcs_vpc.test.id
  cidr       = cidrsubnet(hcs_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(hcs_vpc.test.cidr, 4, 1), 1)
}

resource "hcs_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "hcs_ecs_compute_keypair" "test" {
  name = "%[1]s"
}

resource "hcs_ecs_compute_instance" "test" {
  count = 2

  name               = format("%[1]s_%%d", count.index)
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  image_id           = data.hcs_ims_images.test.images[0].id
  security_groups    = [hcs_networking_secgroup.test.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  key_pair = hcs_ecs_compute_keypair.test.name

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}

resource "hcs_ecs_compute_server_group" "test" {
  count = 2

  name     = format("%[1]s_%%d", count.index)
  policies = ["anti-affinity"]

  members = [
    hcs_ecs_compute_instance.test[count.index].id,
  ]
}

resource "hcs_evs_volume" "test" {
  count = 10

  name              = format("%[1]s-%%d", count.index)
  availability_zone = data.hcs_availability_zones.test.names[0]

  device_type = "SCSI"
  volume_type = "SAS"
  size        = 1000 * (count.index+1)
  multiattach = true
}

resource "hcs_compute_volume_attach" "attach_volumes_to_compute_test_1" {
  count = 10

  instance_id = hcs_ecs_compute_instance.test[0].id
  volume_id   = hcs_evs_volume.test[count.index].id
}

resource "hcs_compute_volume_attach" "attach_volumes_to_compute_test_2" {
  count = 10

  instance_id = hcs_ecs_compute_instance.test[1].id
  volume_id   = hcs_evs_volume.test[count.index].id
}

locals {
  attach_members_1 = hcs_ecs_compute_server_group.test[0].members
  attach_members_2 = hcs_ecs_compute_server_group.test[1].members

  attach_devices_1 = [for d in hcs_compute_volume_attach.attach_volumes_to_compute_test_1[*].device : d != ""]
  attach_devices_2 = [for d in hcs_compute_volume_attach.attach_volumes_to_compute_test_2[*].device : d != ""]
}

output "members_attached" {
  value = length(local.attach_members_1) == 1 && length(local.attach_members_2) == 1
}

output "volumes_attached" {
  value = length(local.attach_devices_1) == 10 && length(local.attach_devices_2) == 10
}
`, name)
}
