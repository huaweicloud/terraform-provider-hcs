package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccAddonTemplateDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddonTemplateDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hcs_cce_addon_template.spark_operator_test", "spec"),
				),
			},
		},
	})
}

func testAccAddonTemplateDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_cce_addon_template" "spark_operator_test" {
  cluster_id = "3fa23c13-305a-11ee-84f9-0255ac100110"
  name       = "npd"
  version    = "1.18.10"
}

`)
}
