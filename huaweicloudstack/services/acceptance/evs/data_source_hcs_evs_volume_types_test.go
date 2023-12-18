package evs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccEvsVolumeTypesDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcs_evs_volume_types.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEvsVolumeTypesConfigAll,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "volume_types.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testEvsVolumeTypesConfigAll = `
data "hcs_evs_volume_types" "all" {}
`
