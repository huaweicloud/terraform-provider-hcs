package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getEipResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}
	return eips.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcEip_basic(t *testing.T) {
	var (
		eip eips.PublicIp

		resourceName = "hcs_vpc_eip.test"
		randName     = acceptance.RandomAccResourceName()
		udpateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "5"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
				),
			},
			{
				Config: testAccVpcEip_update(udpateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", udpateName),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.size", "8"),
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

func TestAccVpcEip_share(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceName()
		resourceName = "hcs_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_share(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME),
					resource.TestCheckResourceAttr(resourceName, "bandwidth.0.name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
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

func TestAccVpcEip_deprecated(t *testing.T) {
	var (
		eip eips.PublicIp

		randName     = acceptance.RandomAccResourceName()
		resourceName = "hcs_vpc_eip.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEip_deprecated(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "UNBOUND"),
					resource.TestCheckResourceAttr(resourceName, "publicip.0.type", acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME),
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

func testAccVpcEip_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_eip" "test" {
  publicip {
    type       = "%[2]s"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 5
  }
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccVpcEip_update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_eip" "test" {
  publicip {
    type       = "%[2]s"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 8
  }
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccVpcEip_share(rName string) string {
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
}
`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccVpcEip_deprecated(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_networking_vip" "test" {
  name       = "%[2]s"
  network_id = hcs_vpc_subnet.test.id
}

resource "hcs_vpc_eip" "test" {
  publicip {
    type    = "%[3]s"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
  }
}
`, common.TestVpc(rName), rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}
