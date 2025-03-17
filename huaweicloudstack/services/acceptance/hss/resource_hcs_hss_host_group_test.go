package hss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getHostGroupResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HCS_REGION_NAME
		epsId   = acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST
		groupId = state.Primary.ID
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return QueryHostGroupById(client, region, epsId, groupId)
}

func TestAccHostGroup_basic(t *testing.T) {
	var (
		hostGroup interface{}
		name      = acceptance.RandomAccResourceName()
		rName     = "hcs_hss_host_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&hostGroup,
		getHostGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
				),
			},
			{
				Config: testAccHostGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostGroupImportStateIDFunc(rName),
				// The field `unprotect_host_ids` will be filled in during the creation and editing operations.
				// We only need to add ignore to the test case and do not need to make special instructions in the document.
				ImportStateVerifyIgnore: []string{
					"unprotect_host_ids",
				},
			},
		},
	})
}

const testAccCompute_data = `
data "hcs_availability_zones" "test" {}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 1
  memory_size       = 1
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-c7bb"
}

data "hcs_ims_images" "test" {
  name       = "ecs_mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}
`

func testAccHostGroup_base(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  count = 2

  name                = "%[2]s_${count.index}"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}
`, testAccCompute_data, name)
}

func testAccHostGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_hss_host_group" "test" {
  name                  = "%[2]s"
  host_ids              = slice(hcs_compute_instance.test[*].id, 0, 1)
  enterprise_project_id = "%[3]s"
}
`, testAccHostGroup_base(name), name, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_hss_host_group" "test" {
  name                  = "%[2]s-update"
  host_ids              = hcs_compute_instance.test[*].id
  enterprise_project_id = "%[3]s"
}
`, testAccHostGroup_base(name), name, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		id := rs.Primary.ID
		if epsId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<id>', but got '%s/%s'",
				epsId, id)
		}
		return fmt.Sprintf("%s/%s", epsId, id), nil
	}
}

func QueryHostGroupById(client *golangsdk.ServiceClient, region, epsId, groupId string) (interface{}, error) {
	allHostGroups, err := queryHostGroupsByName(client, region, epsId, "")
	if err != nil {
		return nil, err
	}

	hostGroup := filterHostGroupById(allHostGroups, groupId)
	if hostGroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return hostGroup, nil
}

func queryHostGroupsByName(client *golangsdk.ServiceClient, region, epsId, groupName string) ([]interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/host-management/groups"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildQueryHostGroupsByNameQueryParams(epsId, groupName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	var (
		offset = 0
		result = make([]interface{}, 0)
	)

	// The `name` parameter is a fuzzy match, so pagination must be used to retrieve all data related to that name.
	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving HSS host groups, %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		hostGroupsResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(hostGroupsResp) == 0 {
			break
		}

		result = append(result, hostGroupsResp...)
		offset += len(hostGroupsResp)
	}

	return result, nil
}

func buildQueryHostGroupsByNameQueryParams(epsId, groupName string) string {
	queryParams := "?limit=20"
	if epsId != "" {
		queryParams += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}
	if groupName != "" {
		queryParams += fmt.Sprintf("&group_name=%v", groupName)
	}

	return queryParams
}

func filterHostGroupById(allHostGroups []interface{}, groupId string) interface{} {
	for _, hostGroup := range allHostGroups {
		if utils.PathSearch("group_id", hostGroup, "").(string) == groupId {
			return hostGroup
		}
	}

	return nil
}
