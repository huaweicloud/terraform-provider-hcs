package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/ports"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getNetworkVipResourceFunc(config *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating hcs VPC network v1 client: %s", err)
	}

	return ports.Get(client, state.Primary.ID)
}

func TestAccNetworkingVip_basic(t *testing.T) {
	var vip ports.Port
	resourceName := "hcs_networking_vip.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getNetworkVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
				),
			},
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName + "_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
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

func TestAccNetworkingVip_ipv6(t *testing.T) {
	var vip ports.Port
	resourceName := "hcs_networking_vip.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getNetworkVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
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

func testAccNetworkingV2VIPConfig_ipv4(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "193.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  vpc_id     = hcs_vpc.test.id
  name       = "%[1]s"
  cidr       = "193.168.1.0/24"
  gateway_ip = "193.168.1.1"
}

resource "hcs_networking_vip" "test" {
  name       = "%[1]s"
  network_id = hcs_vpc_subnet.test.id
}
`, rName)
}

func testAccNetworkingV2VIPConfig_ipv6(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "193.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  vpc_id      = hcs_vpc.test.id
  name        = "%[1]s"
  cidr        = "193.168.0.0/24"
  gateway_ip  = "193.168.0.1"
  ipv6_enable = true
}

resource "hcs_networking_vip" "test" {
  name       = "%[1]s"
  network_id = hcs_vpc_subnet.test.id
  ip_version = 6
}
`, rName)
}
