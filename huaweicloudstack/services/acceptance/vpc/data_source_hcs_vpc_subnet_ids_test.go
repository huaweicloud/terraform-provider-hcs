package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcSubnetIdsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_vpc_subnet_ids.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetIdsDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "2"),
				),
			},
		},
	})
}
func testAccSubnetIdsDataSource_basic(rName string) string {
	return fmt.Sprintf(`

resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "172.16.128.0/20"
}

resource "hcs_vpc_subnet" "test1" {
  name       = "%s"
  cidr       = "172.16.140.0/22"
  gateway_ip = "172.16.140.1"
  vpc_id     = hcs_vpc.test.id
}

resource "hcs_vpc_subnet" "test2" {
  name       = "%s"
  cidr       = "172.16.136.0/22"
  gateway_ip = "172.16.136.1"
  vpc_id     = hcs_vpc.test.id
}

data "hcs_vpc_subnet_ids" "test" {
  vpc_id = hcs_vpc_subnet.test1.vpc_id
}
`, rName, rName, rName)
}
