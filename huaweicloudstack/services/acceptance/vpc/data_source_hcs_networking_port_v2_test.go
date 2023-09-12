package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resourceName := "data.hcs_networking_port.gw_port"
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_fixed_ips.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic(rname string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "my_vpc" {
 name        = "%s"
 cidr        = "191.168.0.0/16"
}

resource "hcs_vpc_subnet" "my_net" {
 name              = "%s"
 cidr              = "191.168.0.0/24"
 gateway_ip        = "191.168.0.1"
 vpc_id            = hcs_vpc.my_vpc.id
}

data "hcs_vpc_subnet" "my_net" {
 id = hcs_vpc_subnet.my_net.id
}

data "hcs_networking_port" "gw_port" {
 network_id = data.hcs_vpc_subnet.my_net.id
 fixed_ip   = data.hcs_vpc_subnet.my_net.gateway_ip
}
`, rname+"_vpc", rname+"_sub")
}
