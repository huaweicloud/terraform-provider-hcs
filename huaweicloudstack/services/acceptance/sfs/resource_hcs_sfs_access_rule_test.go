package sfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/sfs/v2/shares"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getSfsAccessRuleResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS client: %s", err)
	}

	resourceID := state.Primary.ID
	sfsID := state.Primary.Attributes["sfs_id"]
	rules, err := shares.ListAccessRights(client, sfsID).ExtractAccessRights()
	if err != nil {
		return nil, err
	}

	for _, item := range rules {
		if item.ID == resourceID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("the sfs access rule %s does not exist", resourceID)
}

func TestAccSFSAccessRuleV2_basic(t *testing.T) {
	var rule shares.AccessRight
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_sfs_access_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&rule,
		getSfsAccessRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configAccSFSAccessRuleV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
		},
	})
}

func configAccSFSAccessRuleV2_basic(rName string) string {
	return fmt.Sprintf(`
data "hcs_availability_zones" "test" {}

resource "hcs_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = hcs_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "hcs_sfs_file_system" "test" {
  share_proto       = "NFS"
  size              = 10
  name              = "%[1]s"
  description       = "sfs_c2c_test-file"
  access_to         = hcs_vpc.test.id
  access_level      = "rw"
  availability_zone = data.hcs_availability_zones.test.names[0]
}

resource "hcs_sfs_access_rule" "test" {
  sfs_id    = hcs_sfs_file_system.test.id
  access_to = hcs_vpc.test.id
}`, rName)
}
