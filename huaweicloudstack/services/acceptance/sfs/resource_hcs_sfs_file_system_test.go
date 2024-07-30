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

func getSfsFileSystemResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.SfsV2Client(acceptance.HCS_REGION_NAME)
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

	return nil, fmt.Errorf("the sfs file system %s does not exist", resourceID)
}

func TestAccSFSFileSystemV2_basic(t *testing.T) {
	var share shares.Share
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "-update"
	resourceName := "hcs_sfs_file_system.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&share,
		getSfsFileSystemResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "ro"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "cert"),
				),
			},
			{
				Config: testAccSFSFileSystemV2_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "ro"),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_withEpsId(t *testing.T) {
	var share shares.Share
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_sfs_file_system.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&share,
		getSfsFileSystemResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_withoutRule(t *testing.T) {
	var share shares.Share
	rName := acceptance.RandomAccResourceName()
	resourceName := "hcs_sfs_file_system.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&share,
		getSfsFileSystemResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_withoutRule(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "unavailable"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
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

func testAccSFSFileSystemV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_file_system" "test" {
  share_proto       = "NFS"
  size              = 10
  name              = "%s"
  description       = "sfs_c2c_test-file"
  access_to         = hcs_vpc.test.id
  access_level      = "ro"
  availability_zone = data.hcs_availability_zones.test.names[0]
}
`, rName, rName)
}

func testAccSFSFileSystemV2_epsId(rName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_file_system" "test" {
  share_proto           = "NFS"
  size                  = 10
  name                  = "%s"
  description           = "sfs_c2c_test-file"
  access_to             = hcs_vpc.test.id
  access_type           = "cert"
  access_level          = "ro"
  availability_zone     = data.hcs_availability_zones.test.names[0]
  enterprise_project_id = "%s"
}
`, rName, rName, acceptance.HCS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSFSFileSystemV2_update(rName, updateName string) string {
	return fmt.Sprintf(`
resource "hcs_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_file_system" "test" {
  share_proto       = "NFS"
  size              = 20
  name              = "%s"
  description       = "sfs_c2c_test-file"
  access_to         = hcs_vpc.test.id
  access_type       = "cert"
  availability_zone = data.hcs_availability_zones.test.names[0]
}
`, rName, updateName)
}

func testAccSFSFileSystemV2_withoutRule(rName string) string {
	return fmt.Sprintf(`
resource "hcs_sfs_file_system" "test" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs_c2c_test-file"
}
`, rName)
}
