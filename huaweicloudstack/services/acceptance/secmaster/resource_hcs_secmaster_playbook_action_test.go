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

func getPlaybookActionResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getPlaybookAction: Query the SecMaster playbook action detail
	var (
		getPlaybookActionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions"
		getPlaybookActionProduct = "secmaster"
	)
	getPlaybookActionClient, err := cfg.NewServiceClient(getPlaybookActionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookActionPath := getPlaybookActionClient.Endpoint + getPlaybookActionHttpUrl
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{project_id}", getPlaybookActionClient.ProjectID)
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{version_id}", state.Primary.Attributes["version_id"])
	getPlaybookActionPath += "?limit=1000"

	getPlaybookActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookActionResp, err := getPlaybookActionClient.Request("GET", getPlaybookActionPath, &getPlaybookActionOpt)
	if err != nil {
		return nil, err
	}

	getPlaybookActionRespBody, err := utils.FlattenResponse(getPlaybookActionResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", state.Primary.ID)
	getPlaybookActionRespBody = utils.PathSearch(jsonPath, getPlaybookActionRespBody, nil)
	if getPlaybookActionRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getPlaybookActionRespBody, nil
}

func TestAccPlaybookAction_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_secmaster_playbook_action.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookActionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
			acceptance.TestAccPreCheckSecMasterPlaybookVersion(t)
			acceptance.TestAccPreCheckSecMasterPlaybookAction(t)

		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPlaybookAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "action_id", acceptance.HCS_SECMASTER_ACTION_ID),
					resource.TestCheckResourceAttr(rName, "name", "action1"),
					resource.TestCheckResourceAttrPair(rName, "version_id", "hcs_secmaster_playbook_version.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPlaybookAction_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "version_id", "hcs_secmaster_playbook_version.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookActionImportState(rName),
			},
		},
	})
}

func testPlaybookAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_secmaster_playbook_version" "test" {
  workspace_id = "%[2]s"
  playbook_id  = hcs_secmaster_playbook.test.id
  dataclass_id = "%[3]s"
  description  = "created by terraform"
}

resource "hcs_secmaster_playbook_action" "test" {
  workspace_id = "%[2]s"
  version_id   = hcs_secmaster_playbook_version.test.id
  action_id    = "%[4]s"
  name         = "action1"
  description  = "created by terraform"
}
`, testPlaybook_basic(name), acceptance.HCS_SECMASTER_WORKSPACE_ID, acceptance.HCS_SECMASTER_DATACLASS_ID, acceptance.HCS_SECMASTER_ACTION_ID)
}

func testPlaybookAction_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_secmaster_playbook_version" "test" {
  workspace_id = "%[2]s"
  playbook_id  = hcs_secmaster_playbook.test.id
  dataclass_id = "%[3]s"
  description  = "created by terraform"
}

resource "hcs_secmaster_playbook_action" "test" {
  workspace_id = "%[2]s"
  version_id   = hcs_secmaster_playbook_version.test.id
  action_id    = "%[4]s"
  name         = "action2"
  description  = ""
}
`, testPlaybook_basic(name), acceptance.HCS_SECMASTER_WORKSPACE_ID, acceptance.HCS_SECMASTER_DATACLASS_ID, acceptance.HCS_SECMASTER_ACTION_ID)
}

func testPlaybookActionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["version_id"], rs.Primary.ID), nil
	}
}
