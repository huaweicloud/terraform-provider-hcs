package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/nat/v2/snats"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getPublicSnatRuleResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return snats.Get(client, state.Primary.ID)
}

func TestAccPublicSnatRule_basic(t *testing.T) {
	var (
		obj snats.Rule

		rName = "hcs_nat_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "hcs_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "hcs_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccPublicSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "hcs_nat_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPublicSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpc_eip" "test" {
  count = 2

  publicip {
    type = "%[3]s"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
  }
}

resource "hcs_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = hcs_vpc.test.id
  subnet_id             = hcs_vpc_subnet.test.id
}
`, common.TestBaseComputeResources(name), name, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccPublicSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_snat_rule" "test" {
  nat_gateway_id = hcs_nat_gateway.test.id
  subnet_id      = hcs_vpc_subnet.test.id
  floating_ip_id = hcs_vpc_eip.test[0].id
  description    = "Created by acc test"
}
`, testAccPublicSnatRule_base(name))
}

func testAccPublicSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_snat_rule" "test" {
  nat_gateway_id = hcs_nat_gateway.test.id
  subnet_id      = hcs_vpc_subnet.test.id
  floating_ip_id = join(",", hcs_vpc_eip.test[*].id)
}
`, testAccPublicSnatRule_base(name))
}
