package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDatasourceMrsClusters_basic(t *testing.T) {
	rName := "data.hcs_mrs_clusters.name_filter"
	dc := hwacceptance.InitDataSourceCheck(rName)
	name := hwacceptance.RandomAccResourceName()
	pwd := hwacceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMrsClusters_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.name", "hcs_mrs_cluster.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.enterprise_project_id",
						"hcs_mrs_cluster.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(rName, "clusters.0.type", "1"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.version",
						"hcs_mrs_cluster.test", "version"),
					resource.TestCheckResourceAttr(rName, "clusters.0.safe_mode", "1"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vnc"),
					resource.TestCheckResourceAttr(rName, "clusters.0.component_list.0.component_name", "Storm"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceMrsClusters_basic(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_mrs_clusters" "status_filter" {
  status = "running"

  depends_on = [
    hcs_mrs_cluster.test
  ]
}

output "status_filter_is_useful" {
  value = length(data.hcs_mrs_clusters.status_filter.clusters) > 0 && alltrue(
    [for v in data.hcs_mrs_clusters.status_filter.clusters[*].status : v == "running"]
  )  
}

data "hcs_mrs_clusters" "name_filter" {
  name   = "%[2]s"
  status = "existing"

  depends_on = [
    hcs_mrs_cluster.test
  ]
}
output "name_filter_is_useful" {
  value = length(data.hcs_mrs_clusters.name_filter.clusters) > 0 && alltrue(
    [for v in data.hcs_mrs_clusters.name_filter.clusters[*].name : v == "%[2]s"]
  )  
}

data "hcs_mrs_clusters" "tags_filter" {
  tags = "foo*bar,key*value"

  depends_on = [
    hcs_mrs_cluster.test
  ]
}
output "tags_filter_is_useful" {
  value = length(data.hcs_mrs_clusters.tags_filter.clusters) > 0 && alltrue(
    [for v in data.hcs_mrs_clusters.tags_filter.clusters[*].tags.foo : v == "bar"]
  )  
}
`, testAccMrsMapReduceCluster_basic(name, pwd, 3), name)
}
