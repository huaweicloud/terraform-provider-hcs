package ddm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getDdmInstanceResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getInstance: Query DDM instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/instances/{instance_id}"
		getInstanceProduct = "ddm"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{instance_id}", fmt.Sprintf("%v", state.Primary.ID))

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmInstance: %s", err)
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", getInstanceRespBody, nil)
	if status == "DELETED" {
		return nil, fmt.Errorf("error get DDM instance")
	}
	return getInstanceRespBody, nil
}

func TestAccDdmInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	rName := "hcs_ddm_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						acceptance.HCS_DDM_FLAVOR_ID, "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						acceptance.HCS_DDM_ENGINE_ID, "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "4"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						acceptance.HCS_DDM_FLAVOR_ID, "flavors.1.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						acceptance.HCS_DDM_ENGINE_ID, "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"hcs_networking_secgroup.test_update", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update_reduce(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						acceptance.HCS_DDM_FLAVOR_ID, "flavors.1.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						acceptance.HCS_DDM_ENGINE_ID, "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"hcs_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"hcs_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"hcs_networking_secgroup.test_update", "id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine_id", "flavor_id"},
			},
		},
	})
}

func testDdmInstance_base(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

`, common.TestBaseNetwork(name))
}

func testDdmInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = "%[3]s"
  node_num          = 2
  engine_id         = "%[4]s"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), name, acceptance.HCS_DDM_FLAVOR_ID, acceptance.HCS_DDM_ENGINE_ID)
}

func testDdmInstance_basic_update(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_networking_secgroup" "test_update" {
  name = "%[2]s"
}

resource "hcs_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = "%[3]s"
  node_num          = 4
  engine_id         = "%[4]s"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test_update.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), updateName, acceptance.HCS_DDM_FLAVOR_ID, acceptance.HCS_DDM_ENGINE_ID)
}

func testDdmInstance_basic_update_reduce(name, updateName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup" "test_update" {
  name = "%[2]s"
}

resource "hcs_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = "%[3]s"
  node_num          = 2
  engine_id         = "%[4]s"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test_update.id

  availability_zones = [
    data.hcs_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), updateName, acceptance.HCS_DDM_FLAVOR_ID, acceptance.HCS_DDM_ENGINE_ID)
}
