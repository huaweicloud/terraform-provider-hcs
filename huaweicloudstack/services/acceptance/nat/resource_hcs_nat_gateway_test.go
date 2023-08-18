package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/nat/v2/gateways"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/nat"
)

func getPublicGatewayResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return gateways.Get(client, state.Primary.ID)
}

func TestAccPublicGateway_basic(t *testing.T) {
	var (
		obj gateways.Gateway

		rName         = "hcs_nat_gateway.test"
		name          = acceptance.RandomAccResourceNameWithDash()
		updateName    = acceptance.RandomAccResourceNameWithDash()
		relatedConfig = common.TestBaseNetwork(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicGateway_basic_step_1(name, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "spec", string(nat.PublicSpecTypeSmall)),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
				),
			},
			{
				Config: testAccPublicGateway_basic_step_2(updateName, relatedConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "spec", string(nat.PublicSpecTypeMedium)),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func testAccPublicGateway_basic_step_1(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  description           = "Created by acc test"
  vpc_id                = hcs_vpc.test.id
  subnet_id             = hcs_vpc_subnet.test.id
}
`, relatedConfig, name)
}

func testAccPublicGateway_basic_step_2(name, relatedConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = hcs_vpc.test.id
  subnet_id             = hcs_vpc_subnet.test.id
}
`, relatedConfig, name)
}
