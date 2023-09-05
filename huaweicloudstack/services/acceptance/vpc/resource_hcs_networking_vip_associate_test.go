package vpc

import (
	"fmt"
	"log"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("hcs_networking_vip_associate.vip_associate_1",
						"port_ids.0", "hcs_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttrPair("hcs_networking_vip_associate.vip_associate_1",
						"vip_id", "hcs_networking_vip.vip_1", "id"),
				),
			},
			{
				ResourceName:      "hcs_networking_vip_associate.vip_associate_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNetworkingV2VIPAssociateImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	hcsConfig := acceptance.TestAccProvider.Meta().(*config.HcsConfig)
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

func testAccNetworkingV2VIPAssociateImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		vip, ok := s.RootModule().Resources["hcs_networking_vip.vip_1"]
		if !ok {
			return "", fmt.Errorf("vip not found: %s", vip)
		}
		instance, ok := s.RootModule().Resources["hcs_compute_instance.test"]
		if !ok {
			return "", fmt.Errorf("port not found: %s", instance)
		}
		if vip.Primary.ID == "" || instance.Primary.Attributes["network.0.port"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", vip.Primary.ID,
				instance.Primary.Attributes["network.0.port"])
		}
		return fmt.Sprintf("%s/%s", vip.Primary.ID, instance.Primary.Attributes["network.0.port"]), nil
	}
}

const testAccCompute_data = `
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_compute_instance" "test" {
  name                = "%s"
  image_id            = data.hcs_images_image.test.id
  stop_before_destroy = true

  network {
    uuid              = data.hcs_vpc_subnet.test.id
    source_dest_check = false
  }
}
`, testAccCompute_data, rName)
}

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_networking_port" "port" {
  port_id = hcs_compute_instance.test.network[0].port
}

resource "hcs_networking_vip" "vip_1" {
  network_id = data.hcs_vpc_subnet.test.id
}

resource "hcs_networking_vip_associate" "vip_associate_1" {
  vip_id   = hcs_networking_vip.vip_1.id
  port_ids = [hcs_compute_instance.test.network[0].port]
}
`, testAccComputeInstance_basic(rName))
}
