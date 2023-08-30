package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getBandwidthResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloudStack Network client: %s", err)
	}
	return bandwidths.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcBandWidth_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidth_basic(randName, 5),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "size", "5"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
					resource.TestCheckResourceAttr(resourceName, "publicips.#", "0"),
				),
			},
			{
				Config: testAccVpcBandWidth_basic(randName+"_update", 6),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "size", "6"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
				),
			},
		},
	})
}

func testAccVpcBandWidth_basic(rName string, size int) string {
	return fmt.Sprintf(`
resource "hcs_vpc_bandwidth" "test" {
  name = "%s"
  size = "%d"
}
`, rName, size)
}
