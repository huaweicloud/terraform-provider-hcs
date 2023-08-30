package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccVPCEndpoint_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_vpcep_endpoint.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpoint_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getVpcepEndpointResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return endpoints.Get(vpcepClient, state.Primary.ID).Extract()
}

func testAccVPCEndpoint_Basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpcep_service" "test" {
  name        = "%[2]s"
  server_type = "VM"
  vpc_id      = hcs_vpc.test.id
  port_id     = hcs_ecs_compute_instance.test.network[0].port
  approval    = false

  port_mapping {
    service_port  = 111
    terminal_port = 222
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = hcs_vpc.test.id
  network_id  = hcs_vpc_subnet.test.id
  enable_dns  = false
}
`, testAccVPCService_base(rName), rName)
}

func testAccVPCService_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_availability_zones" "test" {}

resource "hcs_ecs_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = "%[3]s"
  flavor_id          = "%[4]s"
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = "%[3]s"
    volume_type = "business_type_01"
    volume_size = 20
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HCS_IMAGE_ID, acceptance.HCS_FLAVOR_ID)
}
