package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDatasourceMrsCluster_basic(t *testing.T) {
	rName := "data.hcs_mrs_cluster.test"
	dc := hwacceptance.InitDataSourceCheck(rName)
	name := hwacceptance.RandomAccResourceName()
	pwd := hwacceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMrsCluster_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "clusters.0.cluster_name", rName),
					resource.TestCheckResourceAttr(rName, "clusters.0.availability_zone_id", "data.hcs_availability_zones.test.names[0]"),
					resource.TestCheckResourceAttr(rName, "clusters.0.version", "MRS 3.2.1-LTS.1"),
					resource.TestCheckResourceAttr(rName, "clusters.0.safe_mode", "true"),
					resource.TestCheckResourceAttr(rName, "clusters.0.vpc_id", "hcs_vpc_subnet.test.id"),
					resource.TestCheckResourceAttr(rName, "clusters.0.subnet_id", "hcs_vpc.test.id"),
					resource.TestCheckResourceAttr(rName, "clusters.0.template_id", "mgmt_control_combined_v4.1"),
					resource.TestCheckResourceAttr(rName, "clusters.0.version", "MRS"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.core_data_volume_count"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.core_data_volume_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.data_center"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.deployment_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.external_alternate_ip"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.fee"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_data_volume_count"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_data_volume_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_data_volume_type"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_node_num"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_node_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.master_node_spec_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.data_volume_count"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.data_volume_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.group_name"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.node_num"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.node_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.node_spec_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.root_volume_size"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.node_groups.0.root_volume_type"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.security_groups_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.slave_security_groups_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.oms_alternate_business_ip"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.oms_business_ip"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.oms_business_ip_port"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.private_ip_first"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.tenant_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.total_node_num"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vnc"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.create_at"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.update_at"),
				),
			},
		},
	})
}

func testAccDatasourceMrsCluster_basic(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_mrs_cluster" "test" {
 cluster_id = hcs_mrs_cluster.test.id

 depends_on = [
   hcs_mrs_cluster.test
 ]
}
`, testAccMrsMapReduceCluster_basic(name, pwd, 3), name)
}
