package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/autoscaling/v1/configurations"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccASConfiguration_basic(t *testing.T) {
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_configuration.acc_as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.metadata.some_key", "some_value"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
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

func TestAccASConfiguration_instance(t *testing.T) {
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_as_configuration.acc_as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_instance(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.instance_id",
						"hcs_ecs_compute_instance.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
				},
			},
		},
	})
}

func testAccCheckASConfigurationDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	asClient, err := cfg.AutoscalingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_as_configuration" {
			continue
		}

		_, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS configuration still exists")
		}
	}

	return nil
}

func testAccCheckASConfigurationExists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
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

		found, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Configuration not found")
		}

		configuration = &found
		return nil
	}
}

func testAccASConfiguration_base(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

data "hcs_ims_images" "test" {
  name       = "cirros-arm"
}
data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-c7bb"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_keypair" "acc_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}
`, rName)
}

func testAccASConfiguration_basic(rName string) string {
	return fmt.Sprintf(`
%s
resource "hcs_as_configuration" "acc_as_config" {
  scaling_configuration_name = "%s"

  instance_config {
    flavor             = data.hcs_ecs_compute_flavors.test.ids[0]
    image              = data.hcs_ims_images.test.images[0].id
    key_name           = hcs_ecs_compute_keypair.acc_key.id
    security_group_ids = [data.hcs_networking_secgroups.test.id]

    disk {
      size        = 40
      volume_type = "business_type_01"
      disk_type   = "SYS"
    }
    metadata = {
      some_key = "some_value"
    }
    user_data = <<EOT
#!/bin/sh
echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
EOT
  }
}
`, testAccASConfiguration_base(rName), rName)
}

func testAccASConfiguration_instance(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }
}

resource "hcs_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
    instance_id = hcs_ecs_compute_instance.test.id
    key_name    = hcs_ecs_compute_keypair.acc_key.id
    user_data   = "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkISBUaGUgdGltZSBpcyBub3cgJChkYXRlIC1SKSEiIHwgdGVlIC9yb290L291dHB1dC50eHQK"
  }
}
`, testAccASConfiguration_base(rName), rName, rName)
}
