package eip

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccEIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcs_vpc_eip_associate.test"
	resourceName := "hcs_vpc_eip.test"
	partten := `^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`

	// hcs_vpc_eip_associate and hcs_vpc_eip have the same ID
	// and call the same API to get resource
	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
					resource.TestMatchOutput("public_ip_address", regexp.MustCompile(partten)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_port(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcs_vpc_eip_associate.test"
	resourceName := "hcs_vpc_eip.test"

	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_port(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(associateName, "status", "BOUND"),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEIPAssociate_compatible(t *testing.T) {
	var eip eips.PublicIp
	rName := acceptance.RandomAccResourceName()
	associateName := "hcs_networking_eip_associate.test"
	resourceName := "hcs_vpc_eip.test"

	rc := acceptance.InitResourceCheck(
		associateName,
		&eip,
		getEipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociate_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(
						associateName, "port_id", &eip.PortID),
					resource.TestCheckResourceAttrPair(
						associateName, "public_ip", resourceName, "address"),
				),
			},
			{
				ResourceName:      associateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEIPAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc_eip" "test" {
  publicip {
    type = "%[2]s"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s"
  }
}`, rName, acceptance.HCS_EIP_EXTERNAL_NETWORK_NAME)
}

func testAccEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

data "hcs_availability_zones" "test" {}

resource "hcs_ecs_compute_instance" "test" {
  name               = "%[3]s"
  image_id           = "%[4]s"
  flavor_id          = "%[5]s"
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = "%[4]s"
    volume_type = "business_type_01"
    volume_size = 20
  }
}

resource "hcs_vpc_eip_associate" "test" {
  public_ip  = hcs_vpc_eip.test.address
  network_id = hcs_ecs_compute_instance.test.network[0].uuid
  fixed_ip   = hcs_ecs_compute_instance.test.network[0].fixed_ip_v4
}

data "hcs_ecs_compute_instance" "test" {
  depends_on = [hcs_vpc_eip_associate.test]

  name = "%[3]s"
}

output "public_ip_address" {
  value = data.hcs_ecs_compute_instance.test.public_ip
}
`, testAccEIPAssociate_base(rName), common.TestBaseNetwork(rName), rName, acceptance.HCS_IMAGE_ID, acceptance.HCS_FLAVOR_ID)
}

func testAccEIPAssociate_port(rName string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "hcs_networking_vip" "test" {
  name       = "%[3]s"
  network_id = hcs_vpc_subnet.test.id
}

resource "hcs_vpc_eip_associate" "test" {
  public_ip = hcs_vpc_eip.test.address
  port_id   = hcs_networking_vip.test.id
}
`, common.TestVpc(rName), testAccEIPAssociate_base(rName), rName)
}

func testAccEIPAssociate_compatible(rName string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "hcs_networking_vip" "test" {
  name       = "%[3]s"
  network_id = hcs_vpc_subnet.test.id
}
  
resource "hcs_networking_eip_associate" "test" {
  public_ip = hcs_vpc_eip.test.address
  port_id   = hcs_networking_vip.test.id
}
`, common.TestVpc(rName), testAccEIPAssociate_base(rName), rName)
}
