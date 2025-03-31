package secmaster

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

func getPlaybookResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	workspaceID := state.Primary.Attributes["workspace_id"]
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return GetPlaybook(client, workspaceID, state.Primary.ID)
}

func TestAccPlaybook_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_secmaster_playbook.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPlaybook_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPlaybook_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookImportState(rName),
			},
		},
	})
}

func testPlaybook_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_secmaster_playbook" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = "created by terraform"
}
`, acceptance.HCS_SECMASTER_WORKSPACE_ID, name)
}

func testPlaybook_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_secmaster_playbook" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_update"
  description  = ""
}
`, acceptance.HCS_SECMASTER_WORKSPACE_ID, name)
}

func testPlaybookImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}

func GetPlaybook(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	getPlaybookHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}"

	getPlaybookPath := client.Endpoint + getPlaybookHttpUrl
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{project_id}", client.ProjectID)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{workspace_id}", workspaceId)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{playbook_id}", id)

	getPlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookResp, err := client.Request("GET", getPlaybookPath, &getPlaybookOpt)
	if err != nil {
		return nil, err
	}

	getPlaybookRespBody, err := utils.FlattenResponse(getPlaybookResp)
	if err != nil {
		return nil, err
	}

	playbook := utils.PathSearch("data", getPlaybookRespBody, nil)
	if playbook == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return playbook, nil
}
