package sfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func TestAccSFSFileSystemV2DataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.hcs_sfs_file_system.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "available"),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
				),
			},
		},
	})
}

func testAccSFSFileSystemV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
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

data "hcs_availability_zones" "test" {}

resource "hcs_sfs_file_system" "test" {
  share_proto       = "NFS"
  size              = 10
  name              = "%[1]s"
  description       = "sfs_c2c_test-file"
  access_to         = hcs_vpc.test.id
  access_level      = "ro"
  availability_zone = data.hcs_availability_zones.test.names[0]
}

data "hcs_sfs_file_system" "test" {
  id = hcs_sfs_file_system.test.id
}
`, rName)
}
