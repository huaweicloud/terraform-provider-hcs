package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccEvsVolumesDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcs_evs_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volume.#", "0"),
				),
			},
		},
	})
}

func testAccEvsVolumesDataSource_base(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    size        = number
    multiattach = bool
  }))
  default = [
    {suffix = "vbd_normal_volume", size = 1, device_type = "VBD", multiattach = false},
    {suffix = "vbd_share_volume", size = 1, device_type = "VBD", multiattach = true},
    {suffix = "scsi_normal_volume", size = 1, device_type = "SCSI", multiattach = false},
    {suffix = "scsi_share_volume", size = 1, device_type = "SCSI", multiattach = true},
  ]
}

%[1]s

resource "hcs_evs_volume" "test" {
  count = length(var.volume_configuration)
  
  availability_zone = data.hcs_availability_zones.test.names[0]
  volume_type       = "business_type_01"
  name              = "%[2]s_${var.volume_configuration[count.index].suffix}"
  size              = var.volume_configuration[count.index].size
  multiattach       = var.volume_configuration[count.index].multiattach
}
`, common.TestBaseComputeResources(rName), rName)
}

func testAccEvsVolumesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_evs_volumes" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
}
`, testAccEvsVolumesDataSource_base(rName))
}
