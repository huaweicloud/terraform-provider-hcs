package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
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
resource "hcs_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = "%[2]s"
  port_id     = "%[3]s"
  approval    = false

  port_mapping {
    service_port  = 111
    terminal_port = 222
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = "%[2]s"
  network_id  = "%[4]s"
  enable_dns  = false
}
`, rName, acceptance.HCS_VPC_ID, acceptance.HCS_ECS_PORT_ID, acceptance.HCS_NETWORK_ID)
}
