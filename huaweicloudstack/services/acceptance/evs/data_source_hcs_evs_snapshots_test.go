package evs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccEvsSnapshotDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcs_evs_snapshots.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEvsSnapshotConfigAll,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "snapshots.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testEvsSnapshotConfigAll = `
data "hcs_evs_snapshots" "all" {}
`
