package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccDataPublicGateway_basic(t *testing.T) {
	var (
		name            = acceptance.RandomAccResourceName()
		nameFilter      = acceptance.InitDataSourceCheck("data.hcs_nat_gateway.name_filter")
		idFilter        = acceptance.InitDataSourceCheck("data.hcs_nat_gateway.id_filter")
		allParamsFilter = acceptance.InitDataSourceCheck("data.hcs_nat_gateway.all_params_filter")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPublicGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					nameFilter.CheckResourceExists(),
					idFilter.CheckResourceExists(),
					allParamsFilter.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataPublicGateway_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  subnet_id             = hcs_vpc_subnet.test.id
  vpc_id                = hcs_vpc.test.id
}

data "hcs_nat_gateway" "name_filter" {
  name = hcs_nat_gateway.test.name
}

data "hcs_nat_gateway" "id_filter" {
  id = hcs_nat_gateway.test.id
}

data "hcs_nat_gateway" "all_params_filter" {
  name                  = hcs_nat_gateway.test.name
  spec                  = hcs_nat_gateway.test.spec
  subnet_id             = hcs_nat_gateway.test.subnet_id
  vpc_id                = hcs_nat_gateway.test.vpc_id
}
`, common.TestBaseNetwork(name), name)
}
