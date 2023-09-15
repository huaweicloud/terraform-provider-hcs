package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccASGroup_basic(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.source_dest_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "300"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zones.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_instances",
				},
			},
			{
				Config: testASGroup_basic_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSED"),
				),
			},
			{
				Config: testASGroup_basic_enable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "15"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func TestAccASGroup_withEpsId(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccASGroup_forceDelete(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_forceDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func TestAccASGroup_sourceDestCheck(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_sourceDestCheck(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "networks.0.source_dest_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func testAccCheckASGroupDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	asClient, err := cfg.AutoscalingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_as_group" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	return nil
}

func testAccCheckASGroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		asClient, err := cfg.AutoscalingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating autoscaling client: %s", err)
		}

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Group not found")
		}

		group = &found
		return nil
	}
}

func testASGroup_Base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_ecs_compute_keypair" "acc_key" {
  name       = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "hcs_elb_loadbalancer" "loadbalancer_1" {
  name          = "%[2]s"
  ipv4_subnet_id = hcs_vpc_subnet.test.ipv4_subnet_id
}

resource "hcs_elb_listener" "listener_1" {
  name            = "%[2]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = hcs_elb_loadbalancer.loadbalancer_1.id
}

resource "hcs_elb_pool" "pool_1" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcs_elb_listener.listener_1.id
}

resource "hcs_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[2]s"
  instance_config {
	image    = data.hcs_ims_images.test.images[0].id
	flavor   = data.hcs_ecs_compute_flavors.test.ids[0]
    key_name = hcs_ecs_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "business_type_01"
      disk_type   = "SYS"
    }
  }
}`, common.TestBaseComputeResources(rName), rName)
}

func testASGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  max_instance_number      = 5

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
  lbaas_listeners {
	listener_id   = hcs_elb_listener.listener_1.id
    pool_id       = hcs_elb_pool.pool_1.id
    protocol_port = hcs_elb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_basic_disable(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  max_instance_number      = 5
  enable                   = false

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
  lbaas_listeners {
	listener_id   = hcs_elb_listener.listener_1.id
    pool_id       = hcs_elb_pool.pool_1.id
    protocol_port = hcs_elb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_basic_enable(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  max_instance_number      = 5
  enable                   = true

  cool_down_time                     = 600
  health_periodic_audit_time         = 15

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
  lbaas_listeners {
	listener_id   = hcs_elb_listener.listener_1.id
    pool_id       = hcs_elb_pool.pool_1.id
    protocol_port = hcs_elb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  enterprise_project_id    = "%s"

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
  lbaas_listeners {
	listener_id   = hcs_elb_listener.listener_1.id
    pool_id       = hcs_elb_pool.pool_1.id
    protocol_port = hcs_elb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName, acceptance.HCS_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testASGroup_forceDelete(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  min_instance_number      = 2
  max_instance_number      = 5
  force_delete             = true
  vpc_id                   = hcs_vpc.test.id

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_sourceDestCheck(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id

  networks {
    id                = hcs_vpc_subnet.test.id
    source_dest_check = false
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
}
`, testASGroup_Base(rName), rName)
}
