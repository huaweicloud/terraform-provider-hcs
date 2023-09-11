package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/listeners"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getELBListenerResourceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return listeners.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Listener_basic(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_elb_listener.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getELBListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3ListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "false"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "true"),
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

func testAccElbV3ListenerConfig_basic(rName string) string {
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

resource "hcs_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = hcs_elb_loadbalancer.test.id
  advanced_forwarding_enabled = false

  idle_timeout = 62
  request_timeout = 63
  response_timeout = 64

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, rName)
}

func testAccElbV3ListenerConfig_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name              = "%s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.hcs_vpc_subnet.test.ipv4_subnet_id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "hcs_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = hcs_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true

  idle_timeout = 62
  request_timeout = 63
  response_timeout = 64

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, rNameUpdate, rNameUpdate)
}
