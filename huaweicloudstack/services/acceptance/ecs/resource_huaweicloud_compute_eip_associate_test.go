package ecs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
	bandwidthsv1 "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccComputeEIPAssociate_basic(t *testing.T) {
	var instance cloudservers.CloudServer
	var eip eips.PublicIp

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_eip_associate.test"
	partten := `^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeEIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeEIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("hcs_compute_instance.test", &instance),
					testAccCheckVpcV1EIPExists("hcs_vpc_eip.test", &eip),
					testAccCheckComputeEIPAssociateAssociated(&eip, &instance, 1),
					resource.TestCheckResourceAttrSet(resourceName, "port_id"),
					resource.TestMatchOutput("public_ip_address", regexp.MustCompile(partten)),
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

func TestAccComputeEIPAssociate_fixedIP(t *testing.T) {
	var instance cloudservers.CloudServer
	var eip eips.PublicIp

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcs_ecs_compute_eip_associate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeEIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeEIPAssociate_fixedIP(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists("hcs_compute_instance.test", &instance),
					testAccCheckVpcV1EIPExists("hcs_vpc_eip.test", &eip),
					testAccCheckComputeEIPAssociateAssociated(&eip, &instance, 1),
					resource.TestCheckResourceAttrSet(resourceName, "port_id"),
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

func testAccCheckComputeEIPAssociateDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	computeClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_eip_associate" {
			continue
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		instance, err := cloudservers.Get(computeClient, instanceId).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// But if the instance still exists, then walk through its known addresses
		// and see if there's a floating IP.
		for _, networkAddresses := range instance.Addresses {
			for _, address := range networkAddresses {
				if address.Type == "floating" || address.Type == "fixed" {
					return fmt.Errorf("EIP %s is still attached to instance %s", address.Addr, instanceId)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeEIPAssociateAssociated(
	eip *eips.PublicIp, instance *cloudservers.CloudServer, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		computeClient, err := cfg.ComputeV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		newInstance, err := cloudservers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i++
			if i != n {
				continue
			}
			for _, address := range networkAddresses {
				if address.Type == "floating" && address.Addr == eip.PublicAddress {
					return nil
				}
			}
		}
		return fmt.Errorf("EIP %s was not attached to instance %s", eip.PublicAddress, instance.ID)
	}
}

func testAccCheckVpcV1EIPExists(n string, eip *eips.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		networkingClient, err := cfg.NetworkingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating networking client: %s", err)
		}

		found, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		*eip = found

		return nil
	}
}

func testAccComputeEIPAssociate_Base(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [data.hcs_networking_secgroups.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }
}

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "eip"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%s"
    size        = 10
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeEIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_eip_associate" "test" {
  public_ip   = hcs_vpc_eip.test.address
  instance_id = hcs_ecs_compute_instance.test.id
}

data "hcs_ecs_compute_instance" "default" {
  depends_on = [hcs_ecs_compute_eip_associate.test]

  name = hcs_ecs_compute_instance.test.name
}

output "public_ip_address" {
  value = data.hcs_ecs_compute_instance.default.public_ip
}
`, testAccComputeEIPAssociate_Base(rName))
}

func testAccComputeEIPAssociate_fixedIP(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_eip_associate" "test" {
  public_ip   = hcs_vpc_eip.test.address
  instance_id = hcs_ecs_compute_instance.test.id
  fixed_ip    = hcs_ecs_compute_instance.test.access_ip_v4
}
`, testAccComputeEIPAssociate_Base(rName))
}

func TestAccComputeEIPAssociate_bandwidth(t *testing.T) {
	var portInfo bandwidthsv1.PublicIpinfo
	randName := acceptance.RandomAccResourceNameWithDash()

	resourceName := "hcs_ecs_compute_eip_associate.test"
	bwResourceName := "hcs_vpc_bandwidth.bandwidth_1"
	ecsResourceName := "hcs_ecs_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBandWidthAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeEIPAssociate_bandwidth(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandWidthAssociateExists(resourceName, &portInfo),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", bwResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "port_id", ecsResourceName, "network.0.port"),
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

func testAccCheckBandWidthAssociateDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	client, err := cfg.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating VPC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ecs_compute_eip_associate" {
			continue
		}

		bwID := rs.Primary.Attributes["bandwidth_id"]
		ipv6PortID := rs.Primary.Attributes["port_id"]
		band, err := bandwidthsv1.Get(client, bwID).Extract()
		if err != nil {
			// ignore 404 status code
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return fmt.Errorf("error fetching bandwidth %s: %s", bwID, err)
			}
		} else {
			for _, item := range band.PublicipInfo {
				if item.PublicipId == ipv6PortID {
					return fmt.Errorf("IPv6 port %s still exists in bandwidth %s", ipv6PortID, bwID)
				}
			}
		}
	}

	return nil
}

func testAccCheckBandWidthAssociateExists(n string, info *bandwidthsv1.PublicIpinfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("bandwidth associate resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		client, err := cfg.NetworkingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating VPC client: %s", err)
		}

		bwID := rs.Primary.Attributes["bandwidth_id"]
		ipv6PortID := rs.Primary.Attributes["port_id"]
		band, err := bandwidthsv1.Get(client, bwID).Extract()
		if err != nil {
			return fmt.Errorf("error fetching bandwidth %s: %s", bwID, err)
		}

		for _, item := range band.PublicipInfo {
			if item.PublicipId == ipv6PortID {
				*info = item
				return nil
			}
		}

		return fmt.Errorf("resource not found: IPv6 port %s does not exist in bandwidth %s",
			ipv6PortID, bwID)
	}
}

func testAccComputeEIPAssociate_bandwidth(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

data "hcs_ims_images" "image_1" {
  name        = "ims-xdoc"
}

resource "hcs_vpc" "vpc_1" {
  name = "%[1]s"
  cidr = "172.16.0.0/16"
}

resource "hcs_vpc_subnet" "subnet_1" {
  vpc_id      = hcs_vpc.vpc_1.id
  name        = "subnet-ipv6"
  cidr        = "172.16.10.0/24"
  gateway_ip  = "172.16.10.1"
  ipv6_enable = true
}

resource "hcs_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "hcs_ecs_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.hcs_ims_images.image_1.id
  availability_zone  = data.hcs_availability_zones.test.names[0]
  flavor_id          = "c6.large.2"
  security_group_ids = [hcs_networking_secgroup.test.id]

  network {
    uuid        = hcs_vpc_subnet.subnet_1.id
    ipv6_enable = true
  }
}

resource "hcs_vpc_bandwidth" "bandwidth_1" {
  name = "%[1]s"
  size = 5
}

resource "hcs_ecs_compute_eip_associate" "test" {
  bandwidth_id = hcs_vpc_bandwidth.bandwidth_1.id
  instance_id  = hcs_ecs_compute_instance.test.id
  fixed_ip     = hcs_ecs_compute_instance.test.access_ip_v6
}
`, rName)
}
