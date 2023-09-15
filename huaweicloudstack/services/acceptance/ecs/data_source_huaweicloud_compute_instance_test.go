package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getEcsInstanceResourceFunc(conf *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ComputeV1Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute v1 client: %s", err)
	}

	resourceID := state.Primary.ID
	found, err := cloudservers.Get(client, resourceID).Extract()
	if err == nil && found.Status == "DELETED" {
		return nil, fmt.Errorf("the resource %s has been deleted", resourceID)
	}

	return found, err
}

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.hcs_ecs_compute_instance.this"
	var instance cloudservers.CloudServer

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rc := acceptance.InitResourceCheck(
		"hcs_ecs_compute_instance.test",
		&instance,
		getEcsInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network.#"),
					resource.TestCheckResourceAttrSet("data.hcs_ecs_compute_instance.byID", "status"),
				),
			},
		},
	})
}

func testAccComputeInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  name               = "%s"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.test.ids[0]
  security_groups = [data.hcs_networking_secgroups.test.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = data.hcs_ims_images.test.images[0].id
    volume_type = "business_type_01"
    volume_size = 20
  }
}
data "hcs_ecs_compute_instance" "this" {
  name = hcs_ecs_compute_instance.test.name
}

data "hcs_ecs_compute_instance" "byID" {
  instance_id = hcs_ecs_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
