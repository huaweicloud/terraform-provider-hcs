package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

const (
	volumeType = "business_type_01"
	size       = "1"
	updateSize = "2"
)

func getVolumeResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.BlockStorageV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloudStack block storage v2 client: %s", err)
	}
	return volumes.Get(c, state.Primary.ID).Extract()
}

func TestAccEvsVolume_basic(t *testing.T) {
	var volume volumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_evs_volume.test"
	resourceName1 := "hcs_evs_volume.test.0"
	resourceName2 := "hcs_evs_volume.test.1"
	resourceName3 := "hcs_evs_volume.test.2"
	resourceName4 := "hcs_evs_volume.test.3"
	resourceName5 := "hcs_evs_volume.test.4"
	resourceName6 := "hcs_evs_volume.test.5"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(6),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.hcs_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Created by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName1, "size", size),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume"),

					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume"),

					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume"),

					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume"),

					resource.TestCheckResourceAttr(resourceName5, "name", rName+"_gpssd2_normal_volume"),
					resource.TestCheckResourceAttr(resourceName5, "volume_type", volumeType),

					resource.TestCheckResourceAttr(resourceName6, "name", rName+"_essd2_normal_volume"),
					resource.TestCheckResourceAttr(resourceName6, "volume_type", volumeType),
				),
			},
			{
				Config: testAccEvsVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(6),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.hcs_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Updated by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName1, "size", updateSize),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume_update"),
					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume_update"),
					resource.TestCheckResourceAttr(resourceName5, "name", rName+"_gpssd2_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName6, "name", rName+"_essd2_normal_volume_update"),
				),
			},
		},
	})
}

func testAccEvsVolume_base() string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    volume_type = string
  }))
  default = [
    {
      suffix = "vbd_normal_volume",
      volume_type = "%[1]s"
    },
    {
      suffix = "vbd_share_volume",
      volume_type = "%[1]s"
    },
    {
      suffix = "scsi_normal_volume",
      volume_type = "%[1]s"
    },
    {
      suffix = "scsi_share_volume",
      volume_type = "%[1]s"
    },
    {
      suffix = "gpssd2_normal_volume",
      volume_type = "%[1]s"
    },
    {
      suffix = "essd2_normal_volume",
      volume_type = "%[1]s"
    },
  ]
}

data "hcs_availability_zones" "test" {}
`, volumeType)
}

func testAccEvsVolume_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.hcs_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}"
  size              = %s
  description       = "Created by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
}
`, testAccEvsVolume_base(), rName, size)
}

func testAccEvsVolume_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.hcs_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}_update"
  size              = %s
  description       = "Updated by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
}
`, testAccEvsVolume_base(), rName, updateSize)
}
