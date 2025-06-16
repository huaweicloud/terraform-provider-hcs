package vdc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"testing"
)

func TestAccVdcRoleDataSource_byExistName(t *testing.T) {
	existName := "readonly"

	dataSourceTestName := "data.hcs_vdc_role.test"
	dataSourceName := "hcs_vdc_role"

	dc := acceptance.InitDataSourceCheck(dataSourceTestName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcRoleByName(dataSourceName, existName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceTestName, "name", existName),
					resource.TestCheckResourceAttr(dataSourceTestName, "display_name", "Tenant Guest"),
					resource.TestCheckResourceAttr(dataSourceTestName, "type", "AA"),
				),
			},
		},
	})
}

func TestAccVdcRoleDataSource_byExistDisplayName(t *testing.T) {
	existDisplayName := "Tenant Guest"
	dataSourceTestName := "data.hcs_vdc_role.test1"
	dataSourceName := "hcs_vdc_role"

	dc := acceptance.InitDataSourceCheck(dataSourceTestName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcRoleByDisplayName(dataSourceName, existDisplayName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceTestName, "name", "readonly"),
					resource.TestCheckResourceAttr(dataSourceTestName, "display_name", existDisplayName),
					resource.TestCheckResourceAttr(dataSourceTestName, "type", "AA"),
				),
			},
		},
	})
}

func testAccDataSourceVdcRoleByName(rName, name string) string {

	return fmt.Sprintf(`
data "%s" "test" {
  name = "%s"
}`, rName, name)
}

func testAccDataSourceVdcRoleByDisplayName(rName, displayName string) string {
	return fmt.Sprintf(`
data "%s" "test1" {
 display_name = "%s"
}
`, rName, displayName)
}
