package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/vpcs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "hcs_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "193.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
				),
			},
			{
				Config: testAccVpcV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "cidr", "193.169.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
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

func TestAccVpcV1_secondaryCIDR(t *testing.T) {
	var vpc vpcs.Vpc

	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_updated"
	resourceName := "hcs_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_secondaryCIDR(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "193.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
				),
			},
			{
				Config: testAccVpcV1_secondaryCIDR_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cidr", "193.169.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secondary_cidr"},
			},
			{
				Config: testAccVpcV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
				),
			},
		},
	})
}

func TestAccVpcV1_WithEpsId(t *testing.T) {
	var vpc vpcs.Vpc

	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "193.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
				),
			},
		},
	})
}

// TestAccVpcV1_WithCustomRegion this case will run a test for resource-level region. Before run this case,
// you shoule set `HCS_CUSTOM_REGION_NAME` in your system and it should be different from `HCS_REGION_NAME`.
func TestAccVpcV1_WithCustomRegion(t *testing.T) {

	vpcName1 := fmt.Sprintf("test_vpc_region_%s", acctest.RandString(5))
	vpcName2 := fmt.Sprintf("test_vpc_region_%s", acctest.RandString(5))

	resName1 := "hcs_vpc.test1"
	resName2 := "hcs_vpc.test2"

	var vpc1, vpc2 vpcs.Vpc

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPrecheckCustomRegion(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_WithCustomRegion(vpcName1, vpcName2, acceptance.HCS_CUSTOM_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRegionVpcV1Exists(resName1, &vpc1, acceptance.HCS_REGION_NAME),
					testAccCheckCustomRegionVpcV1Exists(resName2, &vpc2, acceptance.HCS_CUSTOM_REGION_NAME),
				),
			},
		},
	})
}

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	vpcClient, err := hcsConfig.NetworkingV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating hcs vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_vpc" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckCustomRegionVpcV1Exists(name string, vpc *vpcs.Vpc, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmtp.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		vpcClient, err := hcsConfig.NetworkingV1Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("vpc not found")
		}

		*vpc = *found
		return nil
	}
}

func testAccCheckVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		hcsConfig := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		vpcClient, err := hcsConfig.NetworkingV1Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating hcs vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

func testAccVpcV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name        = "%s"
  cidr        = "193.168.0.0/16"

}
`, rName)
}

func testAccVpcV1_update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name        = "%s"
  cidr        ="193.169.0.0/16"

}
`, rName)
}

func testAccVpcV1_secondaryCIDR(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name           = "%s"
  cidr           = "193.168.0.0/16"
}
`, rName)
}

func testAccVpcV1_secondaryCIDR_update(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name           = "%s"
  cidr           = "193.169.0.0/16"
}
`, rName)
}

func testAccVpcV1_epsId(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name                  = "%s"
  cidr                  = "193.168.0.0/16"
}
`, rName)
}

func testAccVpcV1_WithCustomRegion(name1, name2, region string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test1" {
  name = "%s"
  cidr = "193.168.0.0/16"
}

resource "hcs_vpc" "test2" {   
  region = "%s"
  name   = "%s"
  cidr   = "193.168.0.0/16"
}
`, name1, region, name2)
}
