package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceWafDedicatedInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName1 := "data.hcs_waf_dedicated_instances.instance_1"
	resourceName2 := "data.hcs_waf_dedicated_instances.instance_2"

	dc := acceptance.InitDataSourceCheck(resourceName1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstances_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "name", name),
					resource.TestCheckResourceAttr(resourceName1, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.available_zone"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.ecs_flavor"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.server_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.service_ip"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.run_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.access_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.upgradable"),

					resource.TestCheckResourceAttr(resourceName2, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName2, "name"),
					resource.TestCheckResourceAttrSet(resourceName2, "instances.0.available_zone"),
				),
			},
		},
	})
}

func TestAccDataSourceWafDedicatedInstances_withEpsId(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.hcs_waf_dedicated_instances.instance_1"

	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstances_epsId(name, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccWafDedicatedInstances_conf(name string) string {
	return fmt.Sprintf(`
%s

data "hcs_waf_dedicated_instances" "instance_1" {
  name = hcs_waf_dedicated_instance.instance_1.name

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}

data "hcs_waf_dedicated_instances" "instance_2" {
  id   = hcs_waf_dedicated_instance.instance_1.id
  name = hcs_waf_dedicated_instance.instance_1.name

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name))
}

func testAccWafDedicatedInstances_epsId(name string, epsId string) string {
	return fmt.Sprintf(`
%s

data "hcs_waf_dedicated_instances" "instance_1" {
  id                    = hcs_waf_dedicated_instance.instance_1.id
  name                  = hcs_waf_dedicated_instance.instance_1.name
  enterprise_project_id = "%s"

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsId), epsId)
}
