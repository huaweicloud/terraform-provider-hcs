package vdc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceVdcUser_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	displayName := "terraform_test"
	vdcId := "05432d17-29ed-4ef8-bb4b-f300f47270ec"
	pwd := ""

	dataSourceName := "data.hcs_vdc_user.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVdcUserTest(name, vdcId, pwd, displayName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "vdc_id", vdcId),
					resource.TestCheckResourceAttr(dataSourceName, "auth_type", "LOCAL_AUTH"),
					resource.TestCheckResourceAttr(dataSourceName, "access_mode", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", displayName),
				),
			},
		},
	})
}

func testAccDataSourceVdcUserTest(name string, vdcId string, pwd string, displayName string) string {
	return fmt.Sprintf(`
%s
data "hcs_vdc_user" "test" {
  vdc_id     = hcs_vdc_user.user01.vdc_id
  name       = hcs_vdc_user.user01.name
}
`, testAccVdcUserBasic(name, vdcId, pwd, displayName))
}

func testAccVdcUserBasic(rName string, vdcId string, pwd string, displayName string) string {
	return fmt.Sprintf(`

resource "hcs_vdc_user" "user01" {
  vdc_id = "%s"
  name = "%s"
  password = "%s"
  display_name = "%s"
}
`, vdcId, rName, pwd, displayName)
}
