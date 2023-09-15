package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccComputeInstance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.port"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "network.0.source_dest_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "stop_before_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_eip_on_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "50"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccComputeInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "60"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination", "data_disks",
				},
			},
		},
	})
}

func TestAccComputeInstance_prePaid(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "delete_eip_on_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				Config: testAccComputeInstance_prePaidUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
		},
	})
}

func TestAccComputeInstance_spot(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_spot(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "spot"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination",
					"spot_maximum_price", "spot_duration", "spot_duration_count",
				},
			},
		},
	})
}

func TestAccComputeInstance_powerAction(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_powerAction(rName, "OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "ON"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "ON"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-REBOOT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-REBOOT"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccComputeInstance_powerAction(rName, "FORCE-OFF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "power_action", "FORCE-OFF"),
					resource.TestCheckResourceAttr(resourceName, "status", "SHUTOFF"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"delete_eip_on_termination",
					"power_action",
				},
			},
		},
	})
}

func TestAccComputeInstance_disk_encryption(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_disk_encryption(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccComputeInstance_withEPS(t *testing.T) {
	var instance cloudservers.CloudServer

	srcEPS := acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstance_withEPS(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccComputeInstance_withEPS(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy", "delete_eip_on_termination",
				},
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	computeClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_instance" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		computeClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

const testAccCompute_data = `
data "hcs_availability_zones" "test" {}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 1
  memory_size       = 1
}

data "hcs_vpc_subnets" "test" {
  name = "subnet_9879"
}

data "hcs_ims_images" "test" {
  name       = "ecs_mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "secgroup_2"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name                = "%s"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  stop_before_destroy = true
  availability_zone = data.hcs_availability_zones.test.names[0]
  agency_name         = "%s"

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 50

  data_disks {
    type = "business_type_01"
    size = "10"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeInstance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name                = "%s-update"
  description         = "terraform test update"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  stop_before_destroy = true
  agency_name         = "%s"

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 60

  data_disks {
    type = "business_type_01"
    size = "10"
  }

  tags = {
    foo = "bar2"
    key2 = "value2"
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeInstance_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }

  eip_type = "EIP"
  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }

  eip_type = "EIP"
  bandwidth {
    share_type  = "PER"
    size        = 5
    charge_mode = "bandwidth"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_spot(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  charging_mode      = "spot"
  spot_duration      = 2

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeInstance_powerAction(rName, powerAction string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  power_action       = "%s"

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }
}
`, testAccCompute_data, rName, powerAction)
}

func testAccComputeInstance_disk_encryption(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_kms_key" "test" {
  key_alias       = "%s"
  pending_days    = "7"
  key_description = "first test key"
  is_enabled      = true
}

resource "hcs_ecs_compute_instance" "test" {
  name                = "%s"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  stop_before_destroy = true

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 50

  data_disks {
    type = "business_type_01"
    size = "10"
    kms_key_id = hcs_kms_key.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeInstance_withEPS(rName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name                  = "%s"
  description           = "terraform test"
  image_id              = data.hcs_ims_images.test.images[0].id
  flavor_id             = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids    = [data.hcs_networking_secgroups.test.security_groups[0].id]
  enterprise_project_id = "%s"
  system_disk_type      = "business_type_01"
  system_disk_size      = 40

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName, epsID)
}
