package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/nat/v2/dnats"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getPublicDnatRuleResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return dnats.Get(client, state.Primary.ID)
}

func TestAccPublicDnatRule_basic(t *testing.T) {
	var (
		obj dnats.Rule

		rName = "hcs_nat_dnat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicDnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "hcs_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "hcs_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "port_id", "hcs_ecs_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttr(rName, "protocol", "udp"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "internal_service_port", "80"),
					resource.TestCheckResourceAttr(rName, "external_service_port", "8080"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPublicDnatRule_basic2(t *testing.T) {
	var (
		obj dnats.Rule

		rName2 = "hcs_nat_dnat_rule.test2"
		name   = acceptance.RandomAccResourceNameWithDash()
	)

	rc2 := acceptance.InitResourceCheck(
		rName2,
		&obj,
		getPublicDnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc2.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicDnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "protocol", "tcp"),
					resource.TestCheckResourceAttr(rName2, "internal_service_port", "0"),
					resource.TestCheckResourceAttr(rName2, "external_service_port", "0"),
					resource.TestCheckResourceAttr(rName2, "internal_service_port_range", "23-823"),
					resource.TestCheckResourceAttr(rName2, "external_service_port_range", "8023-8823"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPublicDnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "%[3]s"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
  }
}

resource "hcs_ecs_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = "%[4]s"
  flavor_id          = "%[5]s"
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = "%[4]s"
    volume_type = "business_type_01"
    volume_size = 20
  }
}

resource "hcs_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = hcs_vpc.test.id
  subnet_id             = hcs_vpc_subnet.test.id
}
`, common.TestBaseComputeResources(name), name,
		acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME, acceptance.HCS_IMAGE_ID, acceptance.HCS_FLAVOR_ID)
}

func testAccPublicDnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_dnat_rule" "test" {
  nat_gateway_id        = hcs_nat_gateway.test.id
  floating_ip_id        = hcs_vpc_eip.test.id
  port_id               = hcs_ecs_compute_instance.test.network[0].port
  description           = "Created by acc test"
  protocol              = "udp"
  internal_service_port = 80
  external_service_port = 8080
}
`, testAccPublicDnatRule_base(name))
}

func testAccPublicDnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_dnat_rule" "test2" {
  nat_gateway_id              = hcs_nat_gateway.test.id
  floating_ip_id              = hcs_vpc_eip.test.id
  port_id                     = hcs_ecs_compute_instance.test.network[0].port
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
`, testAccPublicDnatRule_base(name))
}
