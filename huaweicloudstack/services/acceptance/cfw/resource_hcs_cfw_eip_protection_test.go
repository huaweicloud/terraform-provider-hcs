package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getEipProtectionResourceFunc(hcsCfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	var (
		queryHttpUrl            = "v1/{project_id}/eips/protect"
		getProtectedEipsProduct = "cfw"
	)
	cfg := hcsCfg.Config
	client, err := cfg.NewServiceClient(getProtectedEipsProduct, acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	resp, err := cfw.QuerySyncedEips(client, queryHttpUrl, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting protected EIPs: %s", err)
	}
	if !cfw.ProtectedEipExist(resp) {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccEipProtection_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "hcs_cfw_eip_protection.test"

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getEipProtectionResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckEip(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEipProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testEipProtection_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cfw_eip_protection" "test" {
  object_id = data.hcs_cfw_firewalls.test.records[0].protect_objects[0].object_id

  protected_eip {
    id          = "%[2]s"
    public_ipv4 = "%[3]s"
  }
}
`, testAccDatasourceFirewalls_basic(), acceptance.HCS_EIP_ID, acceptance.HCS_EIP_ADDRESS)
}
