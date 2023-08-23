package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccBandWidthDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_vpc_bandwidth.test"
	eipResourceName := "hcs_vpc_eip.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "publicips.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.id",
						eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.ip_address",
						eipResourceName, "address"),
				),
			},
		},
	})
}

func testAccBandWidthDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 10
}

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "%[2]s"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = hcs_vpc_bandwidth.test.id
  }
}

data "hcs_vpc_bandwidth" "test" {
  depends_on = [hcs_vpc_eip.test]

  name = hcs_vpc_bandwidth.test.name
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}
