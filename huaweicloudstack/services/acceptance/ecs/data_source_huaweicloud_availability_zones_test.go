package ecs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccAvailabilityZones_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailabilityZonesConfig_all,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.hcs_availability_zones.all", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccAvailabilityZonesConfig_all = `
data "hcs_availability_zones" "all" {}
`
