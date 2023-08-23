package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccImsImagesDataSource_basic(t *testing.T) {
	imageName := "CentOS_7.4_64bit"
	dataSourceName := "data.hcs_ims_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_publicName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_osVersion("CentOS 7.4 64bit"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_nameRegex("^CentOS_7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_public(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
		},
	})
}

func TestAccImsImagesDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.hcs_ims_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_base(rName),
			},
			{
				Config: testAccImsImagesDataSource_queryName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccImsImagesDataSource_publicName(imageName string) string {
	return fmt.Sprintf(`
data "hcs_ims_images" "test" {
  name       = "%s"
  visibility = "public"
}
`, imageName)
}

func testAccImsImagesDataSource_nameRegex(regexp string) string {
	return fmt.Sprintf(`
data "hcs_ims_images" "test" {
  architecture = "x86_64"
  name_regex   = "%s"
  visibility   = "public"
}
`, regexp)
}

func testAccImsImagesDataSource_osVersion(osVersion string) string {
	return fmt.Sprintf(`
data "hcs_ims_images" "test" {
  architecture = "x86_64"
  os_version   = "%s"
  visibility   = "public"
}
`, osVersion)
}

func testAccImsImagesDataSource_public() string {
	return `
data "hcs_ims_images" "test" {
  os         = "CentOS"
  visibility = "public"
}
`
}

func testAccImsImagesDataSource_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_availability_zones" "test" {}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_instance" "test" {
  name       = "%[2]s"
  image_name = "CentOS_7.4_64bit"
  flavor_id  = data.hcs_ecs_compute_flavors.test.ids[0]

  security_group_ids = [
    hcs_networking_secgroup.test.id
  ]

  availability_zone = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}

resource "hcs_ims_image" "test" {
  name        = "%[2]s"
  instance_id = hcs_ecs_compute_instance.test.id
  description = "created by Terraform AccTest"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccImsImagesDataSource_queryName(rName string) string {
	return fmt.Sprintf(`
%s

data "hcs_ims_images" "test" {
  name = hcs_ims_image.test.name
}
`, testAccImsImagesDataSource_base(rName))
}

func testAccImsImagesDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s
data "hcs_ims_images" "test" {
  visibility = "private"
}
`, testAccImsImagesDataSource_base(rName))
}
