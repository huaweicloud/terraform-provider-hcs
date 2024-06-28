package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceDcsFlavorsV2_basic(t *testing.T) {
	dataSourceName := "data.hcs_dcs_flavors.flavors"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsFlavorsV2_conf(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine", "redis"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.capacity", "2"),
				),
			},
		},
	})
}

func testAccDcsFlavorsV2_conf() string {
	return `
data "hcs_dcs_flavors" "flavors" {
  engine   = "Redis"
  capacity = 2
}
`
}
