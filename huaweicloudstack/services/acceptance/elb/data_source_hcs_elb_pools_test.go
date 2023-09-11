package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDatasourcePools_basic(t *testing.T) {
	rName := "data.hcs_elb_pools.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"hcs_elb_pool.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"hcs_elb_pool.test", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"hcs_elb_pool.test", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"hcs_elb_pool.test", "lb_method"),
				),
			},
		},
	})
}

func testAccDatasourcePools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_elb_pools" "test" {
  name = "%s"

  depends_on = [
    hcs_elb_pool.test
  ]
}
`, testAccElbV3PoolConfig_basic(name), name)
}
