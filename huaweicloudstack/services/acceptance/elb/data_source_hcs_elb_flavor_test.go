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
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.info.connection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.info.cps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.info.qps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.info.bandwidth"),
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
data "hcs_elb_flavors" "test" {
  name = "qos-test"
  depends_on = [
    hcs_elb_flavor.flavor_1
  ]
}

resource "hcs_elb_flavor" "flavor_1" {
    name            = "qos-test"
    type            = "l7"
  info {
    flavor_type = "cps"
    value       = 0
  }
  info {
    flavor_type = "connection"
    value       = 100
  }
  info {
    flavor_type = "bandwidth"
    value       = 300
  }
  info {
    flavor_type = "qps"
    value       = 300
  }
}

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
  connection = data.hcs_elb_flavors.test.flavors[0].info["connection"]
}
data "hcs_elb_flavors" "max_connections_filter" {
  max_connections = data.hcs_elb_flavors.test.flavors[0].info["connection"]
}
output "max_connections_filter_is_useful" {
  value = length(data.hcs_elb_flavors.max_connections_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.max_connections_filter.flavors[*].info["connection"] : v == local.connection]
  )
}

locals {
  cps = data.hcs_elb_flavors.test.flavors[0].info["cps"]
}
data "hcs_elb_flavors" "cps_filter" {
  cps = data.hcs_elb_flavors.test.flavors[0].info["cps"]
}
output "cps_filter_is_useful" {
  value = length(data.hcs_elb_flavors.cps_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.cps_filter.flavors[*].info["cps"] : v == local.cps]
  )
}

locals {
  l7_flavors = [
    for f in data.hcs_elb_flavors.test.flavors : f
    if f.type == "l7"
  ]
  qps = local.l7_flavors[0].info["qps"]
}
data "hcs_elb_flavors" "qps_filter" {
  qps  = local.qps
  type = "l7"
}
output "qps_filter_is_useful" {
  value = local.type == "l7" ? (
    length(data.hcs_elb_flavors.qps_filter.flavors) > 0 && alltrue(
      [for v in data.hcs_elb_flavors.qps_filter.flavors[*].info["qps"] : v == local.qps]
    )
  ) : null
}

locals {
  bandwidth = data.hcs_elb_flavors.test.flavors[0].info["bandwidth"]
}
data "hcs_elb_flavors" "bandwidth_filter" {
  bandwidth = data.hcs_elb_flavors.test.flavors[0].info["bandwidth"]
}
output "bandwidth_filter_is_useful" {
  value = length(data.hcs_elb_flavors.bandwidth_filter.flavors) > 0 && alltrue(
  [for v in data.hcs_elb_flavors.bandwidth_filter.flavors[*].info["bandwidth"] : v == local.bandwidth]
  )
}
`)
}
