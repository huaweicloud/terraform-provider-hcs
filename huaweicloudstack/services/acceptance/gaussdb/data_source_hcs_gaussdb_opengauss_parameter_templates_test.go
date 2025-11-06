package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceConfiguration_basic(t *testing.T) {
	var (
		dataSourceName = "data.hcs_gaussdb_opengauss_parameter_template.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOpengaussConfiguration(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.#"),

					// Check the attributes of the first parameter if exists
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.need_restart"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.is_risk_parameter"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameters.0.risk_description"),
				),
			},
		},
	})
}

func testAccDataSourceConfiguration_basic() string {
	return fmt.Sprintf(`
data "hcs_gaussdb_opengauss_parameter_template" "test" {
  template_id = "%s"
}
`, acceptance.OPENGAUSS_PARAMETER_TEMPLATE_ID)
}
