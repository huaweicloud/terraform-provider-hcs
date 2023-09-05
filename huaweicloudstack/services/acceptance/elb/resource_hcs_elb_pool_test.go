package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccElbV3Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_elb_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3PoolExists(resourceName, &pool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
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

func testAccCheckElbV3PoolDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	elbClient, err := cfg.ElbV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_elb_pool" {
			continue
		}

		_, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		elbClient, err := cfg.ElbV3Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ELB client: %s", err)
		}

		found, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*pool = *found

		return nil
	}
}

func testAccElbV3PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.hcs_vpc_subnet.test.ipv4_subnet_id
}

resource "hcs_elb_listener" "test" {
  name             = "%s"
  description      = "test description"
  protocol         = "HTTP"
  protocol_port    = 8080
  loadbalancer_id  = hcs_elb_loadbalancer.test.id
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "hcs_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcs_elb_listener.test.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}

func testAccElbV3PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.hcs_vpc_subnet.test.ipv4_subnet_id
}

resource "hcs_elb_listener" "test" {
  name             = "%s"
  description      = "test description"
  protocol         = "HTTP"
  protocol_port    = 8080
  loadbalancer_id  = hcs_elb_loadbalancer.test.id
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "hcs_elb_pool" "test" {
  name           = "%s"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  listener_id    = hcs_elb_listener.test.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rNameUpdate)
}
