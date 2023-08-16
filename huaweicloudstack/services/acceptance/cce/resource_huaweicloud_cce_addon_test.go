package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/addons"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getAddonFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceAddonV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return addons.Get(client, state.Primary.ID, state.Primary.Attributes["cluster_id"]).Extract()
}

func TestAccAddon_basic(t *testing.T) {
	var (
		addon addons.Addon

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_addon.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&addon,
			getAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id",
						acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "metrics-server"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAddonImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAddonImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterId, addonId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of CCE add-on is not found in the tfstate", resName)
		}
		clusterId = rs.Primary.Attributes["cluster_id"]
		addonId = rs.Primary.ID
		if clusterId == "" || addonId == "" {
			return "", fmt.Errorf("the CCE add-on ID is not exist or related CCE cluster ID is missing")
		}
		return fmt.Sprintf("%s/%s", clusterId, addonId), nil
	}
}

func testAccAddon_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_addon" "test" {
  cluster_id    = "%[1]s"
  template_name = "metrics-server"
  depends_on    = [hcs_cce_node.test]
}
`, acceptance.HCS_CCE_CLUSTER_ID)
}

func TestAccAddon_values(t *testing.T) {
	var (
		addon addons.Addon

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_addon.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&addon,
			getAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_values_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id",
						acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "version", "1.25.21"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "autoscaler"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccAddon_values_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// the values not set, only check if the updating request is successful
					resource.TestCheckResourceAttr(resourceName, "cluster_id",
						acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "version", "1.25.21"),
					resource.TestCheckResourceAttr(resourceName, "template_name", "autoscaler"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccAddon_values_step1(name string) string {
	return fmt.Sprintf(`
data "hcs_cce_addon_template" "test" {
  cluster_id = "%[1]s"
  name       = "autoscaler"
  version    = "1.25.21"
}

resource "hcs_cce_addon" "test" {
  cluster_id    = "%[1]s"
  template_name = "autoscaler"
  version       = "1.25.21"

  values {
    basic       = jsondecode(data.hcs_cce_addon_template.test.spec).basic
    custom_json = jsonencode(merge(
      jsondecode(data.hcs_cce_addon_template.test.spec).parameters.custom,
      {
        cluster_id = "%[1]s"
        tenant_id  = "%[2]s"
        logLevel   = 3
      }
    ))
    flavor_json = jsonencode(jsondecode(data.hcs_cce_addon_template.test.spec).parameters.flavor1)
  }
}
`, acceptance.HCS_CCE_CLUSTER_ID, acceptance.HCS_PROJECT_ID)
}

func testAccAddon_values_step2(rName string) string {
	return fmt.Sprintf(`
data "hcs_cce_addon_template" "test" {
  cluster_id = "%[1]s"
  name       = "autoscaler"
  version    = "1.25.21"
}

resource "hcs_cce_addon" "test" {
  cluster_id    = "%[1]s"
  template_name = "autoscaler"
  version       = "1.25.21"

  values {
    basic       = jsondecode(data.hcs_cce_addon_template.test.spec).basic
    custom_json = jsonencode(merge(
      jsondecode(data.hcs_cce_addon_template.test.spec).parameters.custom,
      {
        cluster_id = "%[1]s"
        tenant_id  = "%[2]s"
        logLevel   = 4
      }
    ))
    flavor_json = jsonencode(jsondecode(data.hcs_cce_addon_template.test.spec).parameters.flavor2)
  }
}
`, acceptance.HCS_CCE_CLUSTER_ID, acceptance.HCS_PROJECT_ID)
}
