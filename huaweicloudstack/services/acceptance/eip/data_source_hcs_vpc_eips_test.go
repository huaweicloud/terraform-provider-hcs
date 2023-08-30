package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccVpcEipsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_vpc_eips.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEips_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.status", "UNBOUND"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_size", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttrPair(dataSourceName, "eips.0.enterprise_project_id",
						"hcs_vpc_eip.test", "enterprise_project_id"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEips_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpc_eips" "test" {
  enterprise_project_id = hcs_vpc_eip.test.enterprise_project_id
}
`, testAccVpcEip_basic(rName))
}
