package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccVpcEipDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_vpc_eip.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipConfig_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_share_type", "PER"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEipConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_eip" "test" {
  publicip {
    type = "%[2]s"
  }
  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
  }
}

data "hcs_vpc_eip" "test" {
  public_ip = hcs_vpc_eip.test.address
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}
