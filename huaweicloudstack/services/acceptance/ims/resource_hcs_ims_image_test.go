package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ims/v2/cloudimages"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/ims"
)

func TestAccImsImage_basic(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "hcs_ims_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
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

func testAccCheckImsImageDestroy(s *terraform.State) error {
	cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	imageClient, err := cfg.ImageV2Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_ims_image" {
			continue
		}

		_, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("image still exists")
		}
	}

	return nil
}

func testAccCheckImsImageExists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		imageClient, err := cfg.ImageV2Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Image: %s", err)
		}

		found, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

func testAccImsImage_basic(rName string) string {
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
  name               = "%[2]s"
  image_name         = "CentOS_7.4_64bit"
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

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

func testAccImsImage_update(rName string) string {
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
  name               = "%[2]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [hcs_networking_secgroup.test.id]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}

resource "hcs_ims_image" "test" {
  name        = "%[2]s"
  instance_id = hcs_compute_instance.test.id
  description = "updated by Terraform AccTest"
}
`, common.TestBaseNetwork(rName), rName)
}
