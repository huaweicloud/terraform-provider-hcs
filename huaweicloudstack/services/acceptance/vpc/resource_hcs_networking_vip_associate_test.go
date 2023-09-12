package vpc

import (
	"fmt"
	"log"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("hcs_networking_vip_associate.vip_associate_1", "port_ids.0", "hcs_networking_port.vip_ass_test_port", "id"),
					resource.TestCheckResourceAttrPair("hcs_networking_vip_associate.vip_associate_1", "vip_id", "hcs_networking_vip.vip_test", "id"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	networkingClient, err := hcsConfig.NetworkingV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_networking_vip_associate" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	log.Printf("[DEBUG] Destroy NetworkingVIPAssociated success!")
	return nil
}

const (
	vpcCidr         = "176.16.128.0/20"
	subnetCidr      = "176.16.140.0/22"
	subnetGatewayIp = "176.16.140.1"
)

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`

resource "hcs_vpc" "vip_ass_test_vpc" {
  name = "%s"
  cidr = "%s"
}

resource "hcs_vpc_subnet" "vip_ass_test_subnet" {
  name       = "%s"
  cidr       = "%s"
  gateway_ip = "%s"
  vpc_id     = hcs_vpc.vip_ass_test_vpc.id
}

resource "hcs_networking_port" "vip_ass_test_port" {
  name           = "%s"
  network_id     = hcs_vpc_subnet.vip_ass_test_subnet.id
  admin_state_up = "true"
}

resource "hcs_networking_vip" "vip_test" {
  network_id = hcs_vpc_subnet.vip_ass_test_subnet.id
}


resource "hcs_networking_vip_associate" "vip_associate_1" {
  vip_id   = hcs_networking_vip.vip_test.id
  port_ids = [hcs_networking_port.vip_ass_test_port.id]
}


`, rName+"_vpc", vpcCidr, rName+"_sub", subnetCidr, subnetGatewayIp, rName+"_port")
}
