package ucs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getClusterResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	// getCluster: Query the UCS Cluster detail
	var (
		region            = acceptance.HCS_REGION_NAME
		getClusterHttpUrl = "v1/clusters/{id}"
		getClusterProduct = "ucs"
	)
	getClusterClient, err := cfg.NewServiceClient(getClusterProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getClusterPath := getClusterClient.Endpoint + getClusterHttpUrl
	getClusterPath = strings.ReplaceAll(getClusterPath, "{id}", state.Primary.ID)

	getClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getClusterResp, err := getClusterClient.Request("GET", getClusterPath, &getClusterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Cluster: %s", err)
	}

	getClusterRespBody, err := utils.FlattenResponse(getClusterResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Cluster: %s", err)
	}

	return getClusterRespBody, nil
}

func TestAccCluster_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcs_ucs_cluster.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCluster_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "attachedcluster"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "privatek8s"),
					resource.TestCheckResourceAttr(rName, "manage_type", "discrete"),
					resource.TestCheckResourceAttr(rName, "cluster_name", name),
					resource.TestCheckResourceAttr(rName, "country", "CN"),
					resource.TestCheckResourceAttr(rName, "city", "110000"),
				),
			},
			{
				Config: testCluster_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "attachedcluster"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "privatek8s"),
					resource.TestCheckResourceAttr(rName, "manage_type", "grouped"),
					resource.TestCheckResourceAttr(rName, "cluster_name", name),
					resource.TestCheckResourceAttr(rName, "country", "US"),
					resource.TestCheckResourceAttr(rName, "city", "US"),
					resource.TestCheckResourceAttrPair(rName, "fleet_id", "hcs_ucs_fleet.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"annotations",
				},
			},
		},
	})
}

func TestAccCluster_cce(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcs_ucs_cluster.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCluster_cce(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "self"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "cce"),
					resource.TestCheckResourceAttr(rName, "manage_type", "grouped"),
					resource.TestCheckResourceAttrPair(rName, "fleet_id", "hcs_ucs_fleet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "hcs_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "cluster_region", acceptance.HCS_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_project_id", acceptance.HCS_PROJECT_ID),
				),
			},
			{
				Config: testCluster_cce_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "category", "self"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "cce"),
					resource.TestCheckResourceAttr(rName, "manage_type", "discrete"),
					resource.TestCheckResourceAttr(rName, "fleet_id", ""),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "hcs_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "cluster_region", acceptance.HCS_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_project_id", acceptance.HCS_PROJECT_ID),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCluster_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  cluster_version        = "v1.19.16-r1"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "hcs_ucs_fleet" "test" {
  name        = "%[2]s"
  description = "created by terraform"
}
`, common.TestVpc(rName), rName)
}

func testCluster_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ucs_cluster" "test" {
  category     = "attachedcluster"
  cluster_type = "privatek8s"
  cluster_name = "%s"
  country      = "CN"
  city         = "110000"

  annotations = {
    "kubeconfig" = hcs_cce_cluster.test.kube_config_raw
  }
}
`, testCluster_base(name), name)
}

func testCluster_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ucs_cluster" "test" {
  category     = "attachedcluster"
  cluster_type = "privatek8s"
  fleet_id     = hcs_ucs_fleet.test.id
  cluster_name = "%s"
  country      = "US"
  city         = "US"

  annotations = {
    "kubeconfig" = hcs_cce_cluster.test.kube_config_raw
  }
}
`, testCluster_base(name), name)
}

func testCluster_cce(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ucs_cluster" "test" {
  category           = "self"
  cluster_type       = "cce"
  fleet_id           = hcs_ucs_fleet.test.id
  cluster_id         = hcs_cce_cluster.test.id
  cluster_region     = "%s"
  cluster_project_id = "%s"
}
`, testCluster_base(name), acceptance.HCS_REGION_NAME, acceptance.HCS_PROJECT_ID)
}

func testCluster_cce_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ucs_cluster" "test" {
  category           = "self"
  cluster_type       = "cce"
  cluster_id         = hcs_cce_cluster.test.id
  cluster_region     = "%s"
  cluster_project_id = "%s"
}
`, testCluster_base(name), acceptance.HCS_REGION_NAME, acceptance.HCS_PROJECT_ID)
}
