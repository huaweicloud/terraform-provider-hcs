package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDatasourceTemplateDetail_basic(t *testing.T) {
	rName := "data.hcs_dcs_template_detail.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceTemplateDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "params.0.param_id"),
					resource.TestCheckResourceAttrSet(rName, "params.0.param_name"),
					resource.TestCheckResourceAttrSet(rName, "params.0.default_value"),
					resource.TestCheckResourceAttrSet(rName, "params.0.value_range"),
					resource.TestCheckResourceAttrSet(rName, "params.0.value_type"),
					resource.TestCheckResourceAttrSet(rName, "params.0.description"),
					resource.TestCheckResourceAttrSet(rName, "params.0.need_restart"),
					resource.TestCheckOutput("params_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceTemplateDetail_basic() string {
	return `
data "hcs_dcs_templates" "test" {
  type = "sys"
  name = "Default-Redis-4.0-single-generic-DRAM"
}

data "hcs_dcs_template_detail" "test" {
  template_id = data.hcs_dcs_templates.test.templates[0].template_id
  type        = "sys"
}

data "hcs_dcs_template_detail" "params_filter" {
  template_id = data.hcs_dcs_templates.test.templates[0].template_id
  type        = "sys"

  params {
    param_name = "timeout"
  }

  depends_on = [data.hcs_dcs_template_detail.test]
}

output "params_filter_is_useful" {
  value = length(data.hcs_dcs_template_detail.params_filter.params) > 0 && alltrue(
    [for v in data.hcs_dcs_template_detail.params_filter.params[*].param_name : v == "timeout"]
  )  
}
`
}
