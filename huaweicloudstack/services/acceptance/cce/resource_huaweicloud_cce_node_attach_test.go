package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/nodes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getAttachedNodeFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return nodes.Get(client, state.Primary.Attributes["cluster_id"], state.Primary.ID).Extract()
}

func testAccPreCheck(t *testing.T) {
	if acceptance.HCS_CCE_CLUSTER_ID == "" || acceptance.HCS_KEYPAIR_NAME == "" || acceptance.HCS_ECS_INSTANCE_ID == "" {
		t.Skip("This environment [HCS_CCE_CLUSTER_ID,HCS_KEYPAIR_NAME,HCS_ECS_INSTANCE_ID] must be set for test")
	}
}

func TestAccNodeAttach_basic(t *testing.T) {
	var (
		node nodes.Nodes

		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_attach.test"

		baseConfig = testAccNodeAttach_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getAttachedNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			testAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodeAttach_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.10"),
				),
			},
			{
				Config: testAccNodeAttach_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				Config: testAccNodeAttach_basic_step3(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// TODO: Change to other os types
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.10"),
				),
			},
		},
	})
}

func testAccNodeAttach_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[3]s

data "hcs_ims_images" "test" {
  name        = "EulerOS 2.5 64bit"
}
`, common.TestVpc(name), name, common.TestVariables())
}

func testAccNodeAttach_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_attach" "test" {
  cluster_id = var.cluster_id
  server_id  = var.ecs_instance_id
  key_pair   = var.keypair_name
  os         = "EulerOS 2.10"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, baseConfig, name)
}

func testAccNodeAttach_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_attach" "test" {
  cluster_id = var.cluster_id
  server_id  = var.ecs_instance_id
  key_pair   = var.keypair_name
  os         = "EulerOS 2.10"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, baseConfig, name)
}

func testAccNodeAttach_basic_step3(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_attach" "test" {
  cluster_id = var.cluster_id
  server_id  = "febce7bc-43db-410e-83fa-6dcabb806c1d"
  key_pair   = var.keypair_name
  os         = "EulerOS 2.10"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, baseConfig, name)
}
