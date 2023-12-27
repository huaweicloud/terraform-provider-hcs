package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	hwacceptance "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getClusterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getDwsCluster: Query the DWS cluster.
	var (
		getDwsClusterHttpUrl = "v1.0/{project_id}/clusters/{id}"
		getDwsClusterProduct = "dws"
	)
	getDwsClusterClient, err := cfg.NewServiceClient(getDwsClusterProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsClusterPath := getDwsClusterClient.Endpoint + getDwsClusterHttpUrl
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{project_id}", getDwsClusterClient.ProjectID)
	getDwsClusterPath = strings.ReplaceAll(getDwsClusterPath, "{id}", state.Primary.ID)

	getDwsClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	getDwsClusterResp, err := getDwsClusterClient.Request("GET", getDwsClusterPath, &getDwsClusterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsCluster: %s", err)
	}

	getDwsClusterRespBody, err := utils.FlattenResponse(getDwsClusterResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsCluster: %s", err)
	}

	return getDwsClusterRespBody, nil
}

func TestAccResourceDWS_basic(t *testing.T) {
	var obj interface{}

	resourceName := "hcs_dws_cluster.test"
	name := acceptance.RandomAccResourceName()
	randPwd := fmt.Sprintf("%s@!%d", acctest.RandString(5), acctest.RandIntRange(100, 999))

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basic(name, 3, "auto_assign", randPwd, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccDwsCluster_basic(name, 6, "auto_assign", randPwd, "cat"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "6"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "cat"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn", "volume", "endpoints"},
			},
		},
	})
}

func testAccDwsCluster_basic(rName string, numberOfNode int, publicIpBindType, password, tag string) string {
	baseNetwork := common.TestBaseNetwork(rName)

	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dws_cluster" "test" {
  name              = "%s"
  node_type         = "dwsk2.xlarge"
  number_of_node    = %d
  vpc_id            = hcs_vpc.test.id
  network_id        = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = data.hcs_availability_zones.test.names[0]
  user_name         = "test_cluster_admin"
  user_pwd          = "%s"

  public_ip {
    public_bind_type = "%s"
  }

  tags = {
    key = "val"
    foo = "%s"
  }
}
`, baseNetwork, rName, numberOfNode, password, publicIpBindType, tag)
}

func TestAccResourceDWS_basicV2(t *testing.T) {
	var obj interface{}

	resourceName := "hcs_dws_cluster.test"
	name := acceptance.RandomAccResourceName()
	randPwd := fmt.Sprintf("%s@!%d", acctest.RandString(5), acctest.RandIntRange(100, 999))

	rc := hwacceptance.InitResourceCheck(
		resourceName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basicV2(name, 3, "auto_assign", randPwd, "bar", 100),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "version", "8.2.0.103"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.capacity", "100"),
				),
			},
			{
				Config: testAccDwsCluster_basicV2(name, 6, "auto_assign", randPwd, "cat", 150),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "6"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "cat"),
					resource.TestCheckResourceAttr(resourceName, "version", "8.2.0.103"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.capacity", "150"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn", "volume", "endpoints"},
			},
		},
	})
}

func testAccDwsCluster_basicV2(rName string, numberOfNode int, publicIpBindType, password, tag string, volumeCap int) string {
	baseNetwork := common.TestBaseNetwork(rName)

	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_dws_cluster" "test" {
  name              = "%s"
  node_type         = "dwsk2.xlarge"
  number_of_node    = %d
  vpc_id            = hcs_vpc.test.id
  network_id        = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = data.hcs_availability_zones.test.names[0]
  user_name         = "test_cluster_admin"
  user_pwd          = "%s"
  version           = "8.2.0.103"
  number_of_cn      = 3

  public_ip {
    public_bind_type = "%s"
  }

  volume {
    type     = "SSD"
    capacity = %d
  }

  tags = {
    key = "val"
    foo = "%s"
  }
}
`, baseNetwork, rName, numberOfNode, password, publicIpBindType, volumeCap, tag)
}
