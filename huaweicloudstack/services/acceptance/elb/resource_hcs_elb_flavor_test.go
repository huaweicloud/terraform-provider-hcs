package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/elb/v3/flavors"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getELBFlavorResourceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return flavors.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Flavor_basic(t *testing.T) {
	var flavor flavors.Flavor
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "hcs_elb_flavor.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&flavor,
		getELBFlavorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3FlavorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "l7"),
					resource.TestCheckResourceAttr(resourceName, "info.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "info.0.flavor_type", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "info.0.value", "300"),
					resource.TestCheckResourceAttr(resourceName, "info.1.flavor_type", "connection"),
					resource.TestCheckResourceAttr(resourceName, "info.1.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "info.2.flavor_type", "cps"),
					resource.TestCheckResourceAttr(resourceName, "info.2.value", "100"),
					resource.TestCheckResourceAttr(resourceName, "info.3.flavor_type", "qps"),
					resource.TestCheckResourceAttr(resourceName, "info.3.value", "70"),
				),
			},
			{
				Config: testAccElbV3FlavorConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "info.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "info.0.flavor_type", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName, "info.0.value", "400"),
					resource.TestCheckResourceAttr(resourceName, "info.1.flavor_type", "connection"),
					resource.TestCheckResourceAttr(resourceName, "info.1.value", "300"),
					resource.TestCheckResourceAttr(resourceName, "info.2.flavor_type", "cps"),
					resource.TestCheckResourceAttr(resourceName, "info.2.value", "200"),
					resource.TestCheckResourceAttr(resourceName, "info.3.flavor_type", "qps"),
					resource.TestCheckResourceAttr(resourceName, "info.3.flavor_type", "qps"),
					resource.TestCheckResourceAttr(resourceName, "info.3.value", "80"),
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

func testAccElbV3FlavorConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_elb_flavor" "test" {
  name        = "%s"
  type        = "l7"
  info {
    flavor_type = "bandwidth"
    value       = 300
  }
  info {
    flavor_type = "connection"
    value       = 200
  }
  info {
    flavor_type = "cps"
    value       = 100
  }
  info {
    flavor_type = "qps"
    value       = 70
  }
}
`, rName)
}

func testAccElbV3FlavorConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_elb_flavor" "test" {
  name        = "%s"
  type        = "l7"
  info {
    flavor_type = "bandwidth"
    value       = 400
  }
  info {
    flavor_type = "connection"
    value       = 300
  }
  info {
    flavor_type = "cps"
    value       = 200
  }
  info {
    flavor_type = "qps"
    value       = 80
  }
}
`, rName)
}
