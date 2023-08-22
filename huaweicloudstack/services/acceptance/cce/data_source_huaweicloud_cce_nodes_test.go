package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccNodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcs_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccNodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_cce_nodes" "test" {
  cluster_id = hcs_cce_cluster.test.id
  name       = hcs_cce_node.test.name

  depends_on = [hcs_cce_node.test]
}
`, testAccNode_basic(rName))
}
