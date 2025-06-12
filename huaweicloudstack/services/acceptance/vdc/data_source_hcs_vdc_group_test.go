package vdc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestVdcGroupDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_vdc_group.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcGroupBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${hcs_vdc_group.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVdcGroupBase(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vdc_group" "test" {
  name = "%s"
}
`, rName)
}

func testAccDataSourceVdcGroupBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_vdc_group" "test" {
  name = hcs_vdc_group.test.name

  depends_on = [
	hcs_vdc_group.test
  ]
}
`, testAccDataSourceVdcGroupBase(rName))
}
