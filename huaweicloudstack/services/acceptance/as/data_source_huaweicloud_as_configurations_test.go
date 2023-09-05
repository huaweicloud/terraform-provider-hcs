package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceASConfiguration_basic(t *testing.T) {
	dataSourceName := "data.hcs_as_configurations.configurations"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASConfiguration_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "configurations.0.scaling_configuration_name", name),
				),
			},
		},
	})
}

func testAccDataSourceASConfiguration_conf(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_as_configurations" "configurations" {
  name     = hcs_as_configuration.acc_as_config.scaling_configuration_name
  image_id = hcs_as_configuration.acc_as_config.instance_config.0.image
  depends_on = [hcs_as_configuration.acc_as_config]
}
`, testAccASConfiguration_basic(name))
}
