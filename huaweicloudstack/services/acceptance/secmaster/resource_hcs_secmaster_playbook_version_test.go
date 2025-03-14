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

func getPlaybookVersionResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	workspaceID := state.Primary.Attributes["workspace_id"]
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return GetPlaybookVersion(client, workspaceID, state.Primary.ID)
}

func TestAccPlaybookVersion_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_secmaster_playbook_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
			acceptance.TestAccPreCheckSecMasterPlaybookVersion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPlaybookVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id",
						"hcs_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttr(rName, "trigger_type", "EVENT"),
				),
			},
			{
				Config: testPlaybookVersion_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HCS_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id",
						"hcs_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(rName, "dataobject_create", "true"),
					resource.TestCheckResourceAttr(rName, "dataobject_update", "true"),
					resource.TestCheckResourceAttr(rName, "dataobject_delete", "true"),
					resource.TestCheckResourceAttr(rName, "rule_enable", "true"),
					resource.TestCheckResourceAttr(rName, "TIMER", "EVENT"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookVersionImportState(rName),
			},
		},
	})
}

// The playbook version can be updated after matching a workflow (playbook action).
func testPlaybookVersion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_secmaster_playbook_version" "test" {
  workspace_id = "%[2]s"
  playbook_id  = hcs_secmaster_playbook.test.id
  dataclass_id = "%[3]s"
  description  = "created by terraform"

  dataobject_create = true
  trigger_type      = "EVENT"
}
`, testPlaybook_basic(name), acceptance.HCS_SECMASTER_WORKSPACE_ID, acceptance.HCS_SECMASTER_DATACLASS_ID)
}

func testPlaybookVersion_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_secmaster_playbook_version" "test" {
  workspace_id      = "%[2]s"
  playbook_id       = hcs_secmaster_playbook.test.id
  dataclass_id      = "%[3]s"
  description       = "updated by terraform"
  dataobject_create = true
  dataobject_update = true
  dataobject_delete = true
  rule_enable       = true
  trigger_type      = "TIMER"
}
`, testPlaybook_basic(name), acceptance.HCS_SECMASTER_WORKSPACE_ID, acceptance.HCS_SECMASTER_DATACLASS_ID)
}

func testPlaybookVersionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["playbook_id"], rs.Primary.ID), nil
	}
}

func GetPlaybookVersion(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	getPlaybookVersionHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}"
	getPlaybookVersionPath := client.Endpoint + getPlaybookVersionHttpUrl
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{project_id}", client.ProjectID)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{workspace_id}", workspaceId)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{version_id}", id)

	getPlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookVersionResp, err := client.Request("GET", getPlaybookVersionPath, &getPlaybookVersionOpt)
	if err != nil {
		return nil, err
	}

	getPlaybookVersionRespBody, err := utils.FlattenResponse(getPlaybookVersionResp)
	if err != nil {
		return nil, err
	}

	playbookVersion := utils.PathSearch("data", getPlaybookVersionRespBody, nil)
	if playbookVersion == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return playbookVersion, nil
}
