package vpc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/routetables"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getVpcRTBRouteResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	vpcClient, err := conf.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating VPC client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("the format of resource ID %s is invalid", state.Primary.ID)
	}

	routeTableID := parts[0]
	destination := parts[1]
	routeTable, err := routetables.Get(vpcClient, routeTableID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving VPC route table %s: %s", routeTableID, err)
	}

	var route *routetables.Route
	for _, item := range routeTable.Routes {
		if item.DestinationCIDR == destination {
			route = &item
			break
		}
	}
	if route == nil {
		return nil, fmt.Errorf("can not find the vpc route %s with %s", routeTableID, destination)
	}

	return route, nil
}

func TestAccVpcRTBRoute_basic(t *testing.T) {
	var route routetables.Route
	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_route_table_route.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getVpcRTBRouteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRTBRoute_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
					resource.TestCheckResourceAttr(resourceName, "description", "peering route"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_name"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "nexthop",
						"${hcs_vpc.test2.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "destination",
						"${hcs_vpc.test2.cidr}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${hcs_vpc.test1.id}"),
				),
			},
			{
				Config: testAccVpcRTBRoute_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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

func TestAccVpcRTBRoute_vip(t *testing.T) {
	var route routetables.Route
	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_route_table_route.vip"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getVpcRTBRouteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRTBRoute_vip(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "vip"),
					resource.TestCheckResourceAttr(resourceName, "description", "vip route"),
					resource.TestCheckResourceAttr(resourceName, "route_table_name", randName),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "nexthop",
						"${hcs_networking_vip.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${hcs_vpc.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "route_table_id",
						"${hcs_vpc_route_table.test.id}"),
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

func testAccVpcRTBRoute_base(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/16"
}

resource "hcs_vpc" "test2" {
  name = "%s_2"
  cidr = "192.168.0.0/20"
}

resource "hcs_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = hcs_vpc.test1.id
  peer_vpc_id = hcs_vpc.test2.id
}
`, rName, rName, rName)
}

func testAccVpcRTBRoute_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpc_route_table_route" "test" {
  vpc_id      = hcs_vpc.test1.id
  destination = hcs_vpc.test2.cidr
  type        = "peering"
  nexthop     = hcs_vpc.test2.id
  description = "peering route"
}
`, testAccVpcRTBRoute_base(rName))
}

func testAccVpcRTBRoute_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpc_route_table_route" "test" {
  vpc_id      = hcs_vpc.test1.id
  destination = hcs_vpc.test2.cidr
  type        = "peering"
  nexthop     = hcs_vpc.test2.id
  description = ""
}
`, testAccVpcRTBRoute_base(rName))
}

func testAccVpcRTBRoute_vip(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "172.16.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  vpc_id     = hcs_vpc.test.id
  name       = "%s"
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

resource "hcs_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = hcs_vpc.test.id
  description = "a custom route table"
  subnets     = [hcs_vpc_subnet.test.ipv4_subnet_id]
}

resource "hcs_networking_vip" "test" {
  network_id = hcs_vpc_subnet.test.id
}

resource "hcs_vpc_route_table_route" "vip" {
  vpc_id         = hcs_vpc.test.id
  route_table_id = hcs_vpc_route_table.test.id
  destination    = "10.10.10.0/24"
  type           = "vip"
  nexthop        = hcs_networking_vip.test.id
  description    = "vip route"
}
`, rName, rName, rName)
}
