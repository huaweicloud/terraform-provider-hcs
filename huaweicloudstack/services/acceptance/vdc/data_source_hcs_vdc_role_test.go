package vdc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"testing"
)

func TestAccVdcRoleDataSource_byNotExistDisplayName(t *testing.T) {
	notExistName := acceptance.RandomAccVdcRoleDisplayName()
	dataSourceTestName := "data.hcs_vdc_role.test"
	dataSourceName := "hcs_vdc_role"

	dc := acceptance.InitDataSourceCheck(dataSourceTestName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcRoleByDisplayName(dataSourceName, notExistName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceDestroy(),
				),
			},
		},
	})
}

func TestAccVdcRoleDataSource_byExistDisplayName(t *testing.T) {
	existDisplayName := acceptance.GetExistAccVdcRoleDisplayName()
	dataSourceTestName := "data.hcs_vdc_role.test"
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
					resource.TestCheckResourceAttr(dataSourceTestName, "display_name", existDisplayName),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "id", "${hcs_vdc_role.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "name", "${hcs_vdc_role.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "type", "${hcs_vdc_role.test.type}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "description", "${hcs_vdc_role.test.description}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "policy", "${hcs_vdc_role.test.policy}"),
				),
			},
		},
	})
}

func TestAccVdcRoleDataSource_byNotExistName(t *testing.T) {
	notExistName := acceptance.RandomAccVdcRoleName()
	dataSourceTestName := "data.hcs_vdc_role.test"
	dataSourceName := "hcs_vdc_role"

	dc := acceptance.InitDataSourceCheck(dataSourceTestName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcRoleByName(dataSourceName, notExistName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceDestroy(),
				),
			},
		},
	})
}
func TestAccVdcRoleDataSource_byExistName(t *testing.T) {
	existName := acceptance.GetExistAccVdcRoleName()
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
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "id", "${hcs_vdc_role.test.id}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "display_name", "${hcs_vdc_role.test.display_name}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "type", "${hcs_vdc_role.test.type}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "description", "${hcs_vdc_role.test.description}"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceTestName, "policy", "${hcs_vdc_role.test.policy}"),
				),
			},
		},
	})
}

func testAccDataSourceVdcRoleByName(rName, name string) string {
	return fmt.Sprintf(`
data "%s" "test" {
  name = "%s"
}
`, rName, name)
}

func testAccDataSourceVdcRoleByDisplayName(rName, displayName string) string {
	return fmt.Sprintf(`
data "%s" "test" {
  display_name = "%s"
}
`, rName, displayName)
}
