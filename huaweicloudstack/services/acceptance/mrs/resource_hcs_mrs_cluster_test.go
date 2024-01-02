package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccMrsMapReduceCluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "hcs_mrs_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceCluster_basic(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom_separate(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "hcs_mrs_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_customSeparate(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom_fullsize(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "hcs_mrs_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_customFullsize(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func testAccCheckMRSV2ClusterDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.MrsV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_mrs_cluster" {
			continue
		}

		clusterGet, err := cluster.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("the MRS cluster (%s) is still exists", rs.Primary.ID)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV2ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no MRS cluster ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		mrsClient, err := config.MrsV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating MRS client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*clusterGet = *found
		return nil
	}
}

func testAccMrsMapReduceClusterConfig_base(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/20"
  vpc_id     = hcs_vpc.test.id
  gateway_ip = "192.168.0.1"
}
`, rName, rName)
}

func testAccMrsMapReduceCluster_basic(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "hcs_mrs_cluster" "test" {
  availability_zone  = data.hcs_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.2.1-LTS.1"
  safe_mode          = true
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = hcs_vpc_subnet.test.id
  vpc_id             = hcs_vpc.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_2"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, nodeNum1)
}

func testAccMrsMapReduceClusterConfig_customSeparate(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "hcs_mrs_cluster" "test" {
  availability_zone  = data.hcs_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.1.0"
  safe_mode          = true
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = hcs_vpc_subnet.test.id
  vpc_id             = hcs_vpc.test.id
  template_id        = "mgmt_control_separated_v4"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 5
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:3,4",
      "KerberosServer:3,4",
      "KerberosAdmin:3,4",
      "quorumpeer:3,4,5",
      "NameNode:4,5",
      "Zkfc:4,5",
      "JournalNode:3,4,5",
      "ResourceManager:4,5",
      "JobHistoryServer:5",
      "DBServer:3,5",
      "HttpFS:3,5",
      "TimelineServer:5",
      "RangerAdmin:3,4",
      "UserSync:4",
      "TagSync:4",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, nodeNum1)
}

func testAccMrsMapReduceClusterConfig_customFullsize(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "hcs_mrs_cluster" "test" {
  availability_zone  = data.hcs_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.2.1-LTS.1"
  safe_mode          = true
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = hcs_vpc_subnet.test.id
  vpc_id             = hcs_vpc.test.id
  template_id        = "mgmt_control_data_separated_v4.1"
  component_list     = ["Hadoop", "Ranger", "ZooKeeper"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 9
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:5,6",
      "KerberosServer:5,6",
      "KerberosAdmin:5,6",
      "quorumpeer:5,6,7,8,9",
      "NameNode:3,4",
      "Zkfc:3,4",
      "JournalNode:5,6,7",
      "ResourceManager:8,9",
      "JobHistoryServer:8,9",
      "DBServer:8,9",
      "HttpFS:8,9",
      "TimelineServer:5",
      "RangerAdmin:4,5",
      "UserSync:5",
      "TagSync:5",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, nodeNum1)
}

func TestAccMrsMapReduceCluster_bootstrap(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "hcs_mrs_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsBootstrapScript(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_bootstrap(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.name", "bootstrap_0"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.uri",
						acceptance.HCS_MAPREDUCE_BOOTSTRAP_SCRIPT),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.parameters", "a"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.before_component_start", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.execute_need_sudo_root", "true"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.fail_action", "continue"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.active_master", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.nodes.0", "master_node_default_group"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.nodes.1", "node_group_1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_bootstrap(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_mrs_cluster" "test" {
  availability_zone  = data.hcs_availability_zones.test.names[0]
  name               = "%[2]s"
  version            = "MRS 3.2.1-LTS.1"
  type               = "CUSTOM"
  safe_mode          = false
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  vpc_id             = hcs_vpc.test.id
  subnet_id          = hcs_vpc_subnet.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:2,3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  bootstrap_scripts {
    name                   = "bootstrap_0"
    uri                    = "%[4]s"
    parameters             = "a"
    before_component_start = false
    execute_need_sudo_root = true
    fail_action            = "continue"
    active_master          = false
    nodes = [
      "master_node_default_group",
      "node_group_1"
    ]
  }
}
`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HCS_MAPREDUCE_BOOTSTRAP_SCRIPT)
}
