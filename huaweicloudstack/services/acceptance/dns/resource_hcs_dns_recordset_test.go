package dns

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dns/v2/zones"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"regexp"
	"strings"
	"testing"
)

func getDNSRecordsetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	// getDNSRecordset: Query DNS recordset
	getDNSRecordsetClient, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DNS Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <zone_id>/<recordset_id>")
	}
	zoneID := parts[0]
	recordsetID := parts[1]

	zoneInfo, err := zones.Get(getDNSRecordsetClient, zoneID).Extract()
	if err != nil {
		return "", fmt.Errorf("error getting zone: %s", err)
	}
	zoneType := zoneInfo.ZoneType
	version := "v2.1"
	if zoneType == "private" {
		version = "v2"
	}
	getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	getDNSRecordsetPath := getDNSRecordsetClient.Endpoint + getDNSRecordsetHttpUrl
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", zoneID)
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", recordsetID)

	getDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := getDNSRecordsetClient.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DNS recordset: %s", err)
	}
	return utils.FlattenResponse(getDNSRecordsetResp)
}

func TestAccDNSRecordset_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))
	rName := "hcs_dns_recordset.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a recordset description"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "300"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.0"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "a recordset description update"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "weight", "5"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordset_publicZone(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))
	rName := "hcs_dns_recordset.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_publicZone(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a record set"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "3000"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_publicZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "an updated record set"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "6000"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_updated"),
				),
			},
		},
	})
}

func TestAccDNSRecordset_privateZone(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))
	rName := "hcs_dns_recordset.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSRecordsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSRecordset_privateZone(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a private record set"),
					resource.TestCheckResourceAttr(rName, "status", "DISABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.3"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_private"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_private"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_privateZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("update.%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "a private record set update"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "900"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_private_update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_private_update"),
				),
			},
			{
				Config:      testDNSRecordset_privateZone_updateWeight(name),
				ExpectError: regexp.MustCompile(`private zone do not support.`),
			},
			{
				Config:      testDNSRecordset_privateZone_updateLineID(name),
				ExpectError: regexp.MustCompile(`private zone do not support.`),
			},
		},
	})
}

func testDNSRecordset_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a recordset description"
  status      = "ENABLE"
  ttl         = 300
  records     = ["10.1.0.0"]

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a recordset description update"
  status      = "DISABLE"
  ttl         = 600
  records     = ["\"test records\""]

  tags = {
    key1 = "value1_update"
    key2 = "value2_update"
  }
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_publicZone(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  status      = "ENABLE"
  ttl         = 3000
  records     = ["10.1.0.0"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_publicZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "an updated record set"
  status      = "DISABLE"
  ttl         = 6000
  records     = ["\"test records\""]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_privateZone(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set"
  status      = "DISABLE"
  ttl         = 600
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar_private"
    key = "value_private"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_updateWeight(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]
  weight      = 3

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_updateLineID(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "update.%s"
  type        = "TXT"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["\"test records\""]
  line_id     = "Dianxin_Shanxi"

  tags = {
    foo = "bar_private_update"
    key = "value_private_update"
  }
}
`, testAccDNSZone_private(name), name)
}
