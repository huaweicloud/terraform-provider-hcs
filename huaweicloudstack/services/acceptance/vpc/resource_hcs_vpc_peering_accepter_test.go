package vpc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccVpcPeeringAccepter_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcPeeringAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringAccepter_basic(randName),
				ExpectError: regexp.MustCompile(
					`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckVpcPeeringAccepterDestroy(_ *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

func testAccVpcPeeringAccepter_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test1" {
  name = "%s_1"
  cidr = "192.168.0.0/20"
}

resource "hcs_vpc" "test2" {
  name = "%s_2"
  cidr = "192.168.128.0/20"
}

resource "hcs_vpc_peering" "test" {
  name        = "%s"
  vpc_id      = hcs_vpc.test1.id
  peer_vpc_id = hcs_vpc.test2.id
}

resource "hcs_vpc_peering_accepter" "test" {
  vpc_peering_id = hcs_vpc_peering.test.id

  accept = true
}
`, rName, rName, rName)
}
