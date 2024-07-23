package dms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDmsMaintainWindowDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcs_dms_maintainwindow.maintainwindow1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "seq", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "begin", "22"),
				),
			},
		},
	})
}

var testAccDmsMaintainWindowDataSource_basic = `
data "hcs_dms_maintainwindow" "maintainwindow1" {
  seq = 1
}
`
