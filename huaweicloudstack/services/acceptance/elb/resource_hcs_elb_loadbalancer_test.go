package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/loadbalancers"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getELBResourceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	eipID := state.Primary.Attributes["ipv4_eip_id"]
	eipType := state.Primary.Attributes["iptype"]
	if eipType != "" && eipID != "" {
		eipClient, err := c.NetworkingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
		}

		if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
			return nil, err
		}
	}

	return loadbalancers.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccElbV3LoadBalancer_withEIP(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	eipName := acceptance.HCS_EIP_NAME
	resourceName := "hcs_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", eipName),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func testAccElbV3LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.hcs_vpc_subnet.test.ipv4_subnet_id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name              = "%s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.hcs_vpc_subnet.test.ipv4_subnet_id

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_withEIP(rName string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.hcs_vpc_subnet.test.ipv4_subnet_id

  iptype                = "%s"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5
}
`, rName, acceptance.HCS_EIP_NAME)
}
