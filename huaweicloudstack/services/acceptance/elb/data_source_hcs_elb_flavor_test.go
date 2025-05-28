/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccDataSourceELbFlavorV3_basic(t *testing.T) {
	dataSourceName := "data.hcs_elb_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccElbFlavorsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.shared"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.flavor_sold_out"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.max_connections"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.cps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.qps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ids.#"),

					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("shared_filter_is_useful", "true"),
					resource.TestCheckOutput("max_connections_filter_is_useful", "true"),
					resource.TestCheckOutput("cps_filter_is_useful", "true"),
					resource.TestCheckOutput("qps_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbFlavorsDataSource_basic() string {
	return fmt.Sprintf(`
data "hcs_elb_flavors" "test" {}

locals {
  flavor_id = data.hcs_elb_flavors.test.flavors[0].id
}
data "hcs_elb_flavors" "flavor_id_filter" {
  flavor_id = data.hcs_elb_flavors.test.flavors[0].id
}
output "flavor_id_filter_is_useful" {
  value = length(data.hcs_elb_flavors.flavor_id_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.flavor_id_filter.flavors[*].id : v == local.flavor_id]
  )
}

locals {
  type = "l7"
}
data "hcs_elb_flavors" "type_filter" {
  type = "l7"
}
output "type_filter_is_useful" {
  value = length(data.hcs_elb_flavors.type_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.type_filter.flavors[*].type : v == local.type]
  )
}

locals {
  name = data.hcs_elb_flavors.test.flavors[0].name
}
data "hcs_elb_flavors" "name_filter" {
  name = data.hcs_elb_flavors.test.flavors[0].name
}
output "name_filter_is_useful" {
  value = length(data.hcs_elb_flavors.name_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.name_filter.flavors[*].name : v == local.name]
  )
}

locals {
  shared = true
}
data "hcs_elb_flavors" "shared_filter" {
  shared = true
}
output "shared_filter_is_useful" {
  value = length(data.hcs_elb_flavors.shared_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.shared_filter.flavors[*].shared : v == local.shared]
  )
}

locals {
  max_connections = data.hcs_elb_flavors.test.flavors[0].max_connections
}
data "hcs_elb_flavors" "max_connections_filter" {
  max_connections = data.hcs_elb_flavors.test.flavors[0].max_connections
}
output "max_connections_filter_is_useful" {
  value = length(data.hcs_elb_flavors.max_connections_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.max_connections_filter.flavors[*].max_connections : v == local.max_connections]
  )
}

locals {
  cps = data.hcs_elb_flavors.test.flavors[0].cps
}
data "hcs_elb_flavors" "cps_filter" {
  cps = data.hcs_elb_flavors.test.flavors[0].cps
}
output "cps_filter_is_useful" {
  value = length(data.hcs_elb_flavors.cps_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.cps_filter.flavors[*].cps : v == local.cps]
  )
}

locals {
  l7_flavors = [
    for f in data.hcs_elb_flavors.test.flavors : f
    if f.type == "l7"
  ]
  qps = local.l7_flavors[0].qps
}
data "hcs_elb_flavors" "qps_filter" {
  qps  = local.qps
  type = "l7"
}
output "qps_filter_is_useful" {
  value = local.type == "l7" ? (
    length(data.hcs_elb_flavors.qps_filter.flavors) > 0 && alltrue(
      [for v in data.hcs_elb_flavors.qps_filter.flavors[*].qps : v == local.qps]
    )
  ) : null
}

locals {
  bandwidth = data.hcs_elb_flavors.test.flavors[0].bandwidth
}
data "hcs_elb_flavors" "bandwidth_filter" {
  bandwidth = data.hcs_elb_flavors.test.flavors[0].bandwidth
}
output "bandwidth_filter_is_useful" {
  value = length(data.hcs_elb_flavors.bandwidth_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.bandwidth_filter.flavors[*].bandwidth : v == local.bandwidth]
  )
}
`)
}
