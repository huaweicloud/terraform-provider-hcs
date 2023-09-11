package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/monitors"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccElbV3Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_elb_monitor.monitor_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "interval", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/aa"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.aa.com"),
				),
			},
			{
				Config: testAccElbV3MonitorConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "10"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/bb"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.bb.com"),
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

func testAccCheckElbV3MonitorDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	elbClient, err := cfg.ElbV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_elb_monitor" {
			continue
		}

		_, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("monitor not found")
		}

		*monitor = *found

		return nil
	}
}

func testAccCheckElbV3MonitorConfig(rName string) string {
	return fmt.Sprintf(`
data "hcs_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "hcs_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.hcs_vpc_subnet.test.ipv4_subnet_id
}

resource "hcs_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = hcs_elb_loadbalancer.test.id
}
`, rName, rName)
}

func testAccElbV3MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_elb_monitor" "monitor_1" {
  protocol    = "HTTP"
  interval    = 20
  timeout     = 10
  max_retries = 5
  url_path    = "/aa"
  domain_name = "www.aa.com"
  pool_id     = hcs_elb_pool.test.id
}
`, testAccCheckElbV3MonitorConfig(rName))
}

func testAccElbV3MonitorConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_elb_monitor" "monitor_1" {
  protocol    = "HTTP"
  interval    = 30
  timeout     = 15
  max_retries = 10
  url_path    = "/bb"
  domain_name = "www.bb.com"
  port        = 8888
  pool_id     = hcs_elb_pool.test.id
}
`, testAccCheckElbV3MonitorConfig(rName))
}
