package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccVPCEndpointApproval_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_vpcep_approval.approval"

	rc := acceptance.InitResourceCheck(
		"hcs_vpcep_endpoint.test",
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApproval_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "id", "hcs_vpcep_service.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"hcs_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"hcs_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpcep_service" "test" {
  name        = "%[2]s"
  server_type = "VM"
  vpc_id      = hcs_vpc.test.id
  port_id     = hcs_ecs_compute_instance.test.network[0].port
  approval    = true

  port_mapping {
    service_port  = 333
    terminal_port = 444
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = hcs_vpc.test.id
  network_id  = hcs_vpc_subnet.test.id
  enable_dns  = false
}

resource "hcs_vpcep_approval" "approval" {
  service_id = hcs_vpcep_service.test.id
  endpoints  = [hcs_vpcep_endpoint.test.id]
}
`, testVPCService_base(rName), rName)
}

func testAccVPCEndpointApproval_Update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_vpcep_service" "test" {
  name        = "%[2]s"
  server_type = "VM"
  vpc_id      = hcs_vpc.test.id
  port_id     = hcs_ecs_compute_instance.test.network[0].port
  approval    = true

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = hcs_vpc.test.id
  network_id  = hcs_vpc_subnet.test.id
  enable_dns  = false
}

resource "hcs_vpcep_approval" "approval" {
  service_id = hcs_vpcep_service.test.id
  endpoints  = []
}
`, testVPCService_base(rName), rName)
}

func testVPCService_base(rName string) string {
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
