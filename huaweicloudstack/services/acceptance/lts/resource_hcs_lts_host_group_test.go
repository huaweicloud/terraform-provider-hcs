package lts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getHostGroupResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getHostGroup: Query the LTS HostGroup detail
	var (
		getHostGroupHttpUrl = "v3/{project_id}/lts/host-group-list"
		getHostGroupProduct = "lts"
	)
	getHostGroupClient, err := cfg.NewServiceClient(getHostGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS Client: %s", err)
	}

	getHostGroupPath := getHostGroupClient.Endpoint + getHostGroupHttpUrl
	getHostGroupPath = strings.ReplaceAll(getHostGroupPath, "{project_id}", getHostGroupClient.ProjectID)

	getHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getHostGroupOpt.JSONBody = utils.RemoveNil(buildGetOrDeleteHostGroupBodyParams(state.Primary.ID))
	getHostGroupResp, err := getHostGroupClient.Request("POST", getHostGroupPath, &getHostGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	getHostGroupRespBody, err := utils.FlattenResponse(getHostGroupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?host_group_id=='%s']|[0]", state.Primary.ID)
	getHostGroupRespBody = utils.PathSearch(jsonPath, getHostGroupRespBody, nil)
	if getHostGroupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getHostGroupRespBody, nil
}

func TestAccHostGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_lts_host_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
				),
			},
			{
				Config: testHostGroup_basic_update(name + "-update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
				),
			},
		},
	})
}

func testHostGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_lts_host_group" "test" {
 name = "%s"
 type = "linux"
}

`, name)
}

func testHostGroup_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_lts_host_group" "test" {
 name = "%s"
 type = "linux"
}
`, name)
}

func buildGetOrDeleteHostGroupBodyParams(id string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_group_id_list": []string{
			id,
		},
	}
	return bodyParams
}
