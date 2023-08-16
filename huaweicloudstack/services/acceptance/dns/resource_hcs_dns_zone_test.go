package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dns/v2/zones"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getDNSZoneResourceFunc(c *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	dnsClient, err := c.DnsV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}
	return zones.Get(dnsClient, state.Primary.ID).Extract()
}

func TestAccDNSZone_basic(t *testing.T) {
	var zone zones.Zone
	resourceName := "hcs_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "description", "a zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccDNSZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated zone"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
		},
	})
}

func TestAccDNSZone_private(t *testing.T) {
	var zone zones.Zone
	resourceName := "hcs_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_private(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "description", "a private zone"),
					resource.TestCheckResourceAttr(resourceName, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
				),
			},
		},
	})
}

func TestAccDNSZone_readTTL(t *testing.T) {
	var zone zones.Zone
	resourceName := "hcs_dns_zone.zone_1"
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&zone,
		getDNSZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZone_readTTL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func testAccDNSZone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "default" {
  name = "%s"
  cidr = "1.0.0.0/24"
}

resource "hcs_dns_zone" "zone_1" {
  name        = "%s"
  description = "a zone"
  ttl         = 300
  zone_type             = "private"

  router {
    router_id = resource.hcs_vpc.default.id
  }
}
`, zoneName, zoneName)
}

func testAccDNSZone_update(zoneName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "default" {
  name = "%s"
  cidr = "1.0.0.0/24"
}

resource "hcs_dns_zone" "zone_1" {
  name        = "%s"
  description = "an updated zone"
  ttl         = 600
  zone_type             = "private"

  router {
    router_id = resource.hcs_vpc.default.id
  }
}
`, zoneName, zoneName)
}

func testAccDNSZone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "default" {
  name = "%s"
  cidr = "1.0.0.0/24"
}

resource "hcs_dns_zone" "zone_1" {
  name  = "%s"
  email = "email1@example.com"
  zone_type             = "private"

  router {
    router_id = resource.hcs_vpc.default.id
  }
}
`, zoneName, zoneName)
}

func testAccDNSZone_private(zoneName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "default" {
  name = "%s"
  cidr = "1.0.0.0/24"
}

resource "hcs_dns_zone" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  router {
    router_id = resource.hcs_vpc.default.id
  }
}
`, zoneName, zoneName)
}
