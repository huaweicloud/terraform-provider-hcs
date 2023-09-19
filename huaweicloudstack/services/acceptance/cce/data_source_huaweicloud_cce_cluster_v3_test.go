package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func TestAccCCEClusterV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find cluster data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("cluster data source ID not set ")
		}

		return nil
	}
}

func testAccCCEClusterV3DataSource_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
}
`, common.TestVpc(rName), rName)
}

func testAccCCEClusterV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_cce_cluster" "test" {
  depends_on = [hcs_cce_cluster.test]

  name = hcs_cce_cluster.test.name
}
`, testAccCCEClusterV3DataSource_base(rName))
}
