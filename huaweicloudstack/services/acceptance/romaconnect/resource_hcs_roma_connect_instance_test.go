package romaconnect

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/romaconnect/v2/instances"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccRomaConnectInstance_basic(t *testing.T) {
	var instance instances.RomaInstance
	name := fmt.Sprintf("tf-roma-%s", acctest.RandString(5))
	resourceName := "hcs_roma_connect_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRomaConnectInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRomaConnectInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRomaConnectInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "00300-30101-0--0"),
					resource.TestCheckResourceAttr(resourceName, "flavor_type", "basic"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "hcs_vpc.test.id"),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "hcs_vpc_subnet.test.id"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", "hcs_networking_secgroup.test.id"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_all", "true"),
					resource.TestCheckResourceAttr(resourceName, "cpu_arch", "x86_64"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.trace_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.vpc_client_plain", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.retention_policy", "produce_reject"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.engine_version", "2.7"),
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

func TestAccRomaConnectInstance_arm(t *testing.T) {
	var instance instances.RomaInstance
	name := fmt.Sprintf("tf-roma-%s", acctest.RandString(5))
	resourceName := "hcs_roma_connect_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRomaConnectInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRomaConnectInstance_arm(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRomaConnectInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "00300-30102-0--0"),
					resource.TestCheckResourceAttr(resourceName, "flavor_type", "professional"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "hcs_vpc.test.id"),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "hcs_vpc_subnet.test.id"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", "hcs_networking_secgroup.test.id"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "cpu_arch", "aarch64"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.trace_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.vpc_client_plain", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.retention_policy", "time_base"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.engine_version", "2.3.0"),
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

func TestAccRomaConnectInstance_rocketmq(t *testing.T) {
	var instance instances.RomaInstance
	name := fmt.Sprintf("tf-roma-%s", acctest.RandString(5))
	resourceName := "hcs_roma_connect_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRomaConnectInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRomaConnectInstance_rocketmq(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRomaConnectInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "00300-30103-0--0"),
					resource.TestCheckResourceAttr(resourceName, "flavor_type", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "hcs_vpc.test.id"),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "hcs_vpc_subnet.test.id"),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", "hcs_networking_secgroup.test.id"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_all", "false"),
					resource.TestCheckResourceAttr(resourceName, "cpu_arch", "aarch64"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.rocketmq_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.mqs.enable_acl", "true"),
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

func testAccRomaConnectInstanceBase(name string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = hcs_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "hcs_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}
`, name)
}

func testAccRomaConnectInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_roma_connect_instance" "test" {
  name              = "%[2]s"
  description       = "terraform created"
  product_id        = "00300-30101-0--0"
  available_zones   = ["az0.dc0"]
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  ipv6_enable       = false
  enable_all        = true
  cpu_architecture  = "x86"

  enterprise_project_id = "0"
  maintain_begin        = "22:00:00"
  maintain_end          = "02:00:00"

  mqs {
    connector_enable = false
    enable_publicip  = false
    engine_version   = "2.3.0"
    retention_policy = "produce_reject"
    ssl_enable       = true
    vpc_client_plain = false
    trace_enable     = false
  }
}
`, testAccRomaConnectInstanceBase(name), name)
}

func testAccRomaConnectInstance_arm(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_roma_connect_instance" "test" {
  name              = "%[2]s"
  description       = "terraform created"
  product_id        = "00300-30102-0--0"
  available_zones   = ["az0.dc0"]
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  ipv6_enable       = false
  enable_all        = false
  cpu_architecture  = "arm"

  enterprise_project_id = "0"
  maintain_begin        = "20:00:00"
  maintain_end          = "04:00:00"

  mqs {
    connector_enable = true
    enable_publicip  = true
    engine_version   = "2.7"
    retention_policy = "time_base"
    ssl_enable       = false
    vpc_client_plain = true
    trace_enable     = true
  }
}
`, testAccRomaConnectInstanceBase(name), name)
}

func testAccRomaConnectInstance_rocketmq(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_roma_connect_instance" "test" {
  name              = "%[2]s"
  description       = "terraform created"
  product_id        = "00300-30103-0--0"
  available_zones   = ["az0.dc0"]
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  ipv6_enable       = false
  enable_all        = false
  cpu_architecture  = "arm"

  enterprise_project_id = "0"
  maintain_begin        = "20:00:00"
  maintain_end          = "04:00:00"

  mqs {
    rocketmq_enable = true
    ssl_enable      = true
    enable_acl      = true
  }
}
`, testAccRomaConnectInstanceBase(name), name)
}

func testAccCheckRomaConnectInstanceDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	romaConnectClient, err := cfg.RomaConnectV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ROMA Connect client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_roma_connect_instance" {
			continue
		}

		_, err := instances.Get(romaConnectClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ROMA Connect instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckRomaConnectInstanceExists(n string, c *instances.RomaInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		romaConnectClient, err := cfg.RomaConnectV2Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ROMA Connect client: %s", err)
		}

		found, err := instances.Get(romaConnectClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ROMA Connect instance not found")
		}

		*c = *found
		return nil
	}
}
