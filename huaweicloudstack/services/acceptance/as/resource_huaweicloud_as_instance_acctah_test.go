package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/autoscaling/v1/instances"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

func getASInstanceAttachResourceFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCS_REGION_NAME
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := state.Primary.Attributes["scaling_group_id"]
	instanceID := state.Primary.Attributes["instance_id"]
	page, err := instances.List(client, groupID, nil).AllPages()
	if err != nil {
		return nil, err
	}

	allInstances, err := page.(instances.InstancePage).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to fetching instances in AS group %s: %s", groupID, err)
	}

	for _, ins := range allInstances {
		if ins.ID == instanceID {
			return &ins, nil
		}
	}

	return nil, fmt.Errorf("can not find the instance %s in AS group %s", instanceID, groupID)
}

func TestAccASInstanceAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcs_as_instance_attach.test0"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getASInstanceAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASInstanceAttach_conf(name, "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id", "hcs_as_group.acc_as_group", "id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "hcs_ecs_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "protected", "false"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "true"),
					resource.TestCheckResourceAttr(rName, "status", "STANDBY"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"append_instance"},
			},
		},
	})
}

func testASInstanceAttach_conf(name, protection, standby string) string {
	return fmt.Sprintf(`
%s

resource "hcs_ecs_compute_instance" "test" {
  count              = 2
  name               = "%s-${count.index}"
  description        = "instance for AS attach"
  image_id           = data.hcs_ims_images.test.images[0].id
  flavor_id          = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids = [hcs_networking_secgroup.test.id]

  network {
    uuid = hcs_vpc_subnet.test.id
  }
}

resource "hcs_as_instance_attach" "test0" {
  scaling_group_id = hcs_as_group.acc_as_group.id
  instance_id      = hcs_ecs_compute_instance.test[0].id
  protected        = %[3]s
  standby          = %[4]s
}

resource "hcs_as_instance_attach" "test1" {
  scaling_group_id = hcs_as_group.acc_as_group.id
  instance_id      = hcs_ecs_compute_instance.test[1].id
  protected        = %[3]s
  standby          = %[4]s
}
`, testASGroup_basic(name), name, protection, standby)
}
