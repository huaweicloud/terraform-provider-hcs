package dew

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getKmsGrantResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getGrant: Query the KMS manual Grant
	var (
		getGrantHttpUrl = "v1.0/{project_id}/kms/list-grants"
		getGrantProduct = "kms"
	)
	getGrantClient, err := cfg.NewServiceClient(getGrantProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS Client: %s", err)
	}

	getGrantPath := getGrantClient.Endpoint + getGrantHttpUrl
	getGrantPath = strings.ReplaceAll(getGrantPath, "{project_id}", getGrantClient.ProjectID)

	getGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		JSONBody: map[string]interface{}{
			"key_id": state.Primary.Attributes["key_id"],
			"limit":  100,
		},
	}
	getGrantResp, err := getGrantClient.Request("POST", getGrantPath, &getGrantOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	grantResponseBody, err := utils.FlattenResponse(getGrantResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	searchPath := fmt.Sprintf("grants[?grant_id=='%s']|[0]", state.Primary.ID)
	r, err := jmespath.Search(searchPath, grantResponseBody)
	if err != nil || r == nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	return r, nil
}

func TestAccKmsGrant_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_kms_grant.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsGrantResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKmsGrant_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "operations.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "key_id", "hcs_kms_key.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKmsGrantImportState(rName),
			},
		},
	})
}

func testKmsGrant_basic(name string) string {
	return fmt.Sprintf(`

resource "hcs_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"

  enterprise_project_id = "0"
}

resource "hcs_kms_grant" "test" {
  key_id            = hcs_kms_key.test.id
  grantee_principal = "%[2]s"
  operations        = ["create-datakey", "encrypt-datakey"]
}
`, name, acceptance.HCS_USER_ID)
}

func testKmsGrantImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["key_id"] == "" {
			return "", fmt.Errorf("attribute (key_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["key_id"] + "/" +
			rs.Primary.ID, nil
	}
}
