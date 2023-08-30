package dns

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dns/v2/zones"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func getDNSRecordsetResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
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
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a recordset description update"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.0"),
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
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "a record set"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "3000"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_publicZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "TXT"),
					resource.TestCheckResourceAttr(rName, "description", "an updated record set"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "6000"),
					resource.TestCheckResourceAttr(rName, "records.0", "\"test records\""),
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
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "600"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.3"),
					resource.TestCheckResourceAttrSet(rName, "zone_name"),
				),
			},
			{
				Config: testDNSRecordset_privateZone_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s", name)),
					resource.TestCheckResourceAttr(rName, "type", "A"),
					resource.TestCheckResourceAttr(rName, "description", "a private record set update"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLE"),
					resource.TestCheckResourceAttr(rName, "ttl", "900"),
					resource.TestCheckResourceAttr(rName, "records.0", "10.1.0.3"),
				),
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
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a recordset description update"
  status      = "ENABLE"
  ttl         = 600
  records     = ["10.1.0.0"]
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_publicZone(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "TXT"
  description = "a record set"
  status      = "ENABLE"
  ttl         = 3000
  records     = ["\"test records\""]
}
`, testAccDNSZone_basic(name), name)
}

func testDNSRecordset_publicZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "TXT"
  description = "an updated record set"
  status      = "ENABLE"
  ttl         = 6000
  records     = ["\"test records\""]
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
  status      = "ENABLE"
  ttl         = 600
  records     = ["10.1.0.3"]
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_privateZone_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set update"
  status      = "ENABLE"
  ttl         = 900
  records     = ["10.1.0.3"]
}
`, testAccDNSZone_private(name), name)
}
