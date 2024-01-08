package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v2/jobs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	mrsRes "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
)

func TestAccMrsMapReduceJob_basic(t *testing.T) {
	var job jobs.Job
	resourceName := "hcs_mapreduce_job.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	pwd := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceJobConfig_basic(rName, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2JobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", mrsRes.JobHiveSQL),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMRSClusterSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckMRSV2JobDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.MrsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_mapreduce_job" {
			continue
		}

		_, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("the MRS cluster (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMRSV2JobExists(n string, job *jobs.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no MRS cluster ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.MrsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating MRS client: %s ", err)
		}

		found, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*job = *found
		return nil
	}
}

func testAccMRSClusterSubResourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
	}
}

func testAccMrsMapReduceJobConfig_base(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "hcs_mrs_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "ANALYSIS"
  version            = "MRS 1.9.2"
  safe_mode          = false
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = hcs_vpc_subnet.test.id
  vpc_id             = hcs_vpc.test.id
  component_list     = ["Hadoop", "Spark", "Hive", "Tez"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SSD"
    root_volume_size  = 300
    data_volume_type  = "SSD"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SSD"
    root_volume_size  = 300
    data_volume_type  = "SSD"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}

func testAccMrsMapReduceJobConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "hcs_mrs_job" "test" {
  cluster_id = hcs_mrs_cluster.test.id
  name       = "%s"
  type       = "HiveSql"
  sql        = "CREATE TABLE tname2 (name VARCHAR(50) NOT NULL);"
}`, testAccMrsMapReduceJobConfig_base(rName, pwd), rName)
}
