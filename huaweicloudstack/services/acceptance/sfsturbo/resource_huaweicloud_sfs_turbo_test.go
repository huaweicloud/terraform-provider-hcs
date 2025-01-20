package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/sfs_turbo/v1/shares"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getSfsTurboResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS client: %s", err)
	}

	resourceID := state.Primary.ID
	share, err := shares.Get(client, resourceID).Extract()
	if err != nil {
		return nil, err
	}

	if share.ID == resourceID {
		return &share, nil
	}

	return nil, fmt.Errorf("the sfs turbo %s does not exist", resourceID)
}

func TestAccSFSTurbo_basic(t *testing.T) {
	var turbo shares.Turbo
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "sfsturbo.hdd"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "150"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
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

func TestAccSFSTurbo_ssd(t *testing.T) {
	var turbo shares.Turbo
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_ssd(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "sfsturbo.ssd"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "350"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSFSTurbo_ssd_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "350"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcs_networking_secgroup.test_update", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "232"),
				),
			},
		},
	})
}

func testAccSFSTurbo_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_turbo" "test" {
  name              = "%s"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = data.hcs_availability_zones.test.names[0]

  share_proto = "NFS"
  share_type  = "sfsturbo.hdd"
  size        = 50
  bandwidth   = 150

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_ssd(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_turbo" "test" {
  name              = "%s"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = data.hcs_availability_zones.test.names[0]

  share_proto = "NFS"
  share_type  = "sfsturbo.ssd"
  size        = 50
  bandwidth   = 350

  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_ssd_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_networking_secgroup" "test_update" {
  name                 = "%[2]s_update"
  delete_default_rules = true
}

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_turbo" "test" {
  name              = "%[2]s_update"
  vpc_id            = hcs_vpc.test.id
  subnet_id         = hcs_vpc_subnet.test.id
  security_group_id = hcs_networking_secgroup.test.id
  availability_zone = data.hcs_availability_zones.test.names[0]

  share_proto = "NFS"
  share_type  = "sfsturbo.ssd"
  size        = 100
  bandwidth   = 350

  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}
