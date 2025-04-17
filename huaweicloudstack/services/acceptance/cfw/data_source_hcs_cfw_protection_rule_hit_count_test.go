package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDatasourceHitCount_basic(t *testing.T) {
	rName := "data.hcs_cfw_protection_rule_hit_count.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceHitCount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "records.0.fw_instance_id", acceptance.HCS_CFW_INSTANCE_ID),
				),
			},
		},
	})
}

func TestAccDatasourceHitCount_eps(t *testing.T) {
	rName := "data.hcs_cfw_protection_rule_hit_count.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceHitCount_eps(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "records.0.fw_instance_id", acceptance.HCS_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccDatasourceHitCount_fw_instance(t *testing.T) {
	rName := "data.hcs_cfw_protection_rule_hit_count.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceHitCount_fw_instance(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "records.0.fw_instance_id", acceptance.HCS_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccDatasourceHitCount_base(name string) string {
	return fmt.Sprintf(`
data "hcs_cfw_firewalls" "test" {
  fw_instance_id = "%[1]s"
}

resource "hcs_cfw_protection_rule" "test" {
  name                = "%[2]s"
  object_id           = data.hcs_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "1.1.1.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, acceptance.HCS_CFW_INSTANCE_ID, name)
}

func testAccDatasourceHitCount_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_cfw_protection_rule_hit_count" test {
  rule_ids = [hcs_cfw_protection_rule.test.id]
}
`, testAccDatasourceHitCount_base(name))
}

func testAccDatasourceHitCount_eps(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_cfw_protection_rule_hit_count" test {
  rule_ids = [hcs_cfw_protection_rule.test.id]

  enterprise_project_id = "%[2]s"
}
`, testAccDatasourceHitCount_base(name), acceptance.HCS_ENTERPRISE_PROJECT_ID)
}

func testAccDatasourceHitCount_fw_instance(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_cfw_protection_rule_hit_count" test {
  rule_ids = [hcs_cfw_protection_rule.test.id]

  fw_instance_id = "%[2]s"
}
`, testAccDatasourceHitCount_base(name), acceptance.HCS_CFW_INSTANCE_ID)
}
