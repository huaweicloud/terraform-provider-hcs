package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
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
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"hcs_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
			},
		},
	})
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = true

  port_mapping {
    service_port  = 333
    terminal_port = 444
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = "%[2]s"
  network_id  = "%[4]s"
  enable_dns  = false
}

resource "hcs_vpcep_approval" "approval" {
  service_id = hcs_vpcep_service.test.id
  endpoints  = [hcs_vpcep_endpoint.test.id]
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID, acceptance.HCS_NETWORK_ID)
}

func testAccVPCEndpointApproval_Update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = true

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = "%[2]s"
  network_id  = "%[4]s"
  enable_dns  = false
}

resource "hcs_vpcep_approval" "approval" {
  service_id = hcs_vpcep_service.test.id
  endpoints  = []
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID, acceptance.HCS_NETWORK_ID)
}
