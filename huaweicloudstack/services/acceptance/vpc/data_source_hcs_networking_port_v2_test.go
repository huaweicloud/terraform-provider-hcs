package vpc

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resourceName := "data.hcs_networking_port.gw_port"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_fixed_ips.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic() string {
	return `
data "hcs_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "hcs_networking_port" "gw_port" {
  network_id = data.hcs_vpc_subnet.mynet.id
  fixed_ip   = data.hcs_vpc_subnet.mynet.gateway_ip
}
`
}
