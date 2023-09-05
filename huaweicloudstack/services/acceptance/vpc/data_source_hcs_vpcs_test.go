package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpcs_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpcs" "test" {
  id = hcs_vpc.test.id
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpcs" "test" {
  cidr = hcs_vpc.test.cidr

  depends_on = [
    hcs_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpcs" "test" {
  name = hcs_vpc.test.name

  depends_on = [
    hcs_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byAll(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.hcs_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byAll(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${hcs_vpc.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byAll(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "hcs_vpcs" "test" {
  id                    = hcs_vpc.test.id
  name                  = hcs_vpc.test.name
  cidr                  = hcs_vpc.test.cidr
  status                = "OK"

  depends_on = [
    hcs_vpc.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}
