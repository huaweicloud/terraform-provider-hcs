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

func getIncidentResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getIncident: Query the SecMaster incident detail
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster Client: %s", err)
	}

	return GetIncident(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccIncident_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_secmaster_incident.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIncidentResourceFunc,
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
				Config: testIncident_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test incident"),
					resource.TestCheckResourceAttr(rName, "type.0.category", "DDoS"),
					resource.TestCheckResourceAttr(rName, "type.0.incident_type", "ACK Flood"),
					resource.TestCheckResourceAttr(rName, "level", "Tips"),
					resource.TestCheckResourceAttr(rName, "status", "Open"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_feature", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_name", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.source_type", "1"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", "2023-04-18T13:00:00.000+08:00"),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", "2023-04-19T14:00:00.000+08:00"),
					resource.TestCheckResourceAttr(rName, "verification_status", "Unknown"),
					resource.TestCheckResourceAttr(rName, "stage", "Preparation"),
					resource.TestCheckResourceAttr(rName, "debugging_data", "false"),
					resource.TestCheckResourceAttr(rName, "labels", "test1,test2"),
				),
			},
			{
				Config: testIncident_update(name + "-update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "test incident update"),
					resource.TestCheckResourceAttr(rName, "type.0.category", "Web Attack"),
					resource.TestCheckResourceAttr(rName, "type.0.incident_type", "Black IP"),
					resource.TestCheckResourceAttr(rName, "level", "Low"),
					resource.TestCheckResourceAttr(rName, "status", "Block"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_feature", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_name", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.source_type", "1"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", "2023-04-18T14:00:00.000+08:00"),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", "2023-04-19T15:00:00.000+08:00"),
					resource.TestCheckResourceAttr(rName, "verification_status", "Positive"),
					resource.TestCheckResourceAttr(rName, "stage", "Detection and analysis"),
					resource.TestCheckResourceAttr(rName, "debugging_data", "true"),
					resource.TestCheckResourceAttr(rName, "labels", "test1,test2,test3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIncidentImportState(rName),
				ImportStateVerifyIgnore: []string{
					"updated_at",
				},
			},
		},
	})
}

func testIncident_basic(name string) string {
	return fmt.Sprintf(`
resource "hcs_secmaster_incident" "test" {
  workspace_id = "%s"
  name         = "%s"
  description  = "test incident"

  type {
    category      = "DDoS"
    incident_type = "ACK Flood"
  }

  level  = "Tips"
  status = "Open"

  data_source {
    product_feature = "hss"
    product_name    = "hss"
    source_type     = 1
  }

  first_occurrence_time = "2023-04-18T13:00:00.000+08:00"
  last_occurrence_time  = "2023-04-19T14:00:00.000+08:00"
  verification_status   = "Unknown"
  stage                 = "Preparation"
  debugging_data        = false
  labels                = "test1,test2"
}
`, acceptance.HCS_SECMASTER_WORKSPACE_ID, name)
}

func testIncident_update(name string) string {
	return fmt.Sprintf(`
resource "hcs_secmaster_incident" "test" {
  workspace_id = "%s"
  name         = "%s"
  description  = "test incident update"

  type {
    category      = "Web Attack"
    incident_type = "Black IP"
  }

  level  = "Low"
  status = "Block"

  data_source {
    product_feature = "hss"
    product_name    = "hss"
    source_type     = 1
  }

  first_occurrence_time = "2023-04-18T14:00:00.000+08:00"
  last_occurrence_time  = "2023-04-19T15:00:00.000+08:00"
  verification_status   = "Positive"
  stage                 = "Detection and analysis"
  debugging_data        = true
  labels                = "test1,test2,test3"
}
`, acceptance.HCS_SECMASTER_WORKSPACE_ID, name)
}

func testIncidentImportState(name string) resource.ImportStateIdFunc {
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

func GetIncident(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	// getIncident: Query the SecMaster incident detail
	getIncidentHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/incidents/{id}"
	getIncidentPath := client.Endpoint + getIncidentHttpUrl
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{project_id}", client.ProjectID)
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{workspace_id}", workspaceId)
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{id}", id)

	getIncidentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIncidentResp, err := client.Request("GET", getIncidentPath, &getIncidentOpt)
	if err != nil {
		return nil, err
	}

	getIncidentRespBody, err := utils.FlattenResponse(getIncidentResp)
	if err != nil {
		return nil, err
	}

	dataObject := utils.PathSearch("data.data_object", getIncidentRespBody, nil)
	if dataObject == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataObject, nil
}
