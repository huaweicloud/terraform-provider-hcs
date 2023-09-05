package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/routes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getRouteResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return routes.Get(c, state.Primary.ID).Extract()
}

// TestAccVpcRoute_basic: This function is *deprecated* as the resource ID format
// has changed, please run TestAccVpcRTBRoute_basic
func TestAccVpcRoute_basic(t *testing.T) {
	var route routes.Route

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_route.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getRouteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRoute_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "nexthop",
						"${hcs_vpc_peering_connection.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "destination",
						"${hcs_vpc.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${hcs_vpc.test1.id}"),
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

func testAccRoute_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/20"
}

resource "hcs_vpc" "test2" {
  name = "%s_2"
  cidr = "172.16.128.0/20"
}

resource "hcs_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = hcs_vpc.test1.id
  peer_vpc_id = hcs_vpc.test2.id
}

resource "hcs_vpc_route" "test" {
  type        = "peering"
  nexthop     = hcs_vpc_peering_connection.test.id
  destination = hcs_vpc.test2.cidr
  vpc_id      = hcs_vpc.test1.id
}
`, rName, rName, rName)
}
