package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/swr/v2/namespaces"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getResourceOrganization(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	return namespaces.Get(swrClient, state.Primary.ID).Extract()
}

func TestAccSWROrganization_basic(t *testing.T) {
	var org namespaces.Namespace
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_swr_organization.test"
	loginServer := fmt.Sprintf("swr.%s.%s", acceptance.HCS_REGION_NAME, acceptance.HCS_CLOUD)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&org,
		getResourceOrganization,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrOrganization(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWROrganization_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "permission", "Manage"),
					resource.TestCheckResourceAttr(resourceName, "login_server", loginServer),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSWROrganization_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_swr_organization" "test" {
  name = "%s"
}
`, rName)
}
