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

func getBandwidthAssociateResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloudStack Network client: %s", err)
	}
	bwID := state.Primary.Attributes["bandwidth_id"]
	return bandwidths.Get(client, bwID).Extract()
}

func TestAccBandWidthAssociate_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthAssociate_basic(randName, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "hcs_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "hcs_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcs_vpc_eip.test.0", "address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_basic(randName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "hcs_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcs_vpc_eip.test.1", "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBandWidthAssociate_migrate(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthAssociate_migrate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "hcs_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "hcs_vpc_eip.source", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcs_vpc_eip.source", "address"),
				),
			},
			{
				Config: testAccBandWidthAssociate_owner(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", "hcs_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "hcs_vpc_eip.test", "address"),
				),
			},
		},
	})
}

func testAccBandWidthAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "hcs_vpc_eip" "test" {
  count = 2
  name  = "%[1]s-${count.index}"

  publicip {
    type = "%[2]s"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s-${count.index}"
    size        = 5
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccBandWidthAssociate_basic(rName string, index int) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpc_bandwidth_associate" "test" {
  bandwidth_id = hcs_vpc_bandwidth.test.id
  eip_id       = hcs_vpc_eip.test.%d.id
}
`, testAccBandWidthAssociate_base(rName), index)
}

func testAccBandWidthAssociate_migrate(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "hcs_vpc_bandwidth" "source" {
  name = "%[1]s-source"
  size = 5
}

resource "hcs_vpc_eip" "source" {
  publicip {
    type = "%[2]s"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = hcs_vpc_bandwidth.source.id
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}

resource "hcs_vpc_bandwidth_associate" "test" {
  bandwidth_id = hcs_vpc_bandwidth.test.id
  eip_id       = hcs_vpc_eip.source.id
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccBandWidthAssociate_owner(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "%[2]s"
  }

  bandwidth {
    share_type = "WHOLE"
    id         = hcs_vpc_bandwidth.test.id
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}

resource "hcs_vpc_bandwidth_associate" "test" {
  bandwidth_id = hcs_vpc_bandwidth.test.id
  eip_id       = hcs_vpc_eip.test.id
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}
