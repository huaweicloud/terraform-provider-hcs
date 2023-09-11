package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func TestAccVpcDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpc.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpc_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "id", "${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpc_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpc_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpc" "test" {
  id = hcs_vpc.test.id

  depends_on = [
	hcs_vpc.test
  ]
}
`, testAccDataSourceVpc_base(rName, cidr))
}

func testAccDataSourceVpc_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpc" "test" {
  cidr = hcs_vpc.test.cidr

  depends_on = [
	hcs_vpc.test
  ]
}
`, testAccDataSourceVpc_base(rName, cidr))
}

func testAccDataSourceVpc_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpc" "test" {
  name = hcs_vpc.test.name

  depends_on = [
	hcs_vpc.test
  ]
}
`, testAccDataSourceVpc_base(rName, cidr))
}
