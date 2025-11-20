package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceAZs_basic(t *testing.T) {
	dataSourceName := "data.hcs_drs_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAZs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
			{
				Config: testAccDataSourceAZs_update1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
			{
				Config: testAccDataSourceAZs_update2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
			{
				Config: testAccDataSourceAZs_update3,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
		},
	})
}

const testAccDataSourceAZs_basic string = `data "hcs_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "migration"
  direction   = "up"
  node_type   = "micro"
}`

const testAccDataSourceAZs_update1 string = `data "hcs_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "sync"
  direction   = "down"
  node_type   = "small"
}`

const testAccDataSourceAZs_update2 string = `data "hcs_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "cloudDataGuard"
  direction   = "non-dbs"
  node_type   = "medium"
}`

const testAccDataSourceAZs_update3 string = `data "hcs_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "cloudDataGuard"
  direction   = "non-dbs"
  node_type   = "high"
}`

