package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/nodepools"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getNodePoolFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return nodepools.Get(client, state.Primary.Attributes["cluster_id"], state.Primary.ID).Extract()
}

func TestAccNodePool_basic(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		baseConfig = testAccNodePool_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_basic_step1(name, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
				),
			},
			{
				Config: testAccNodePool_basic_step2(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "scall_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "9"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
			{
				Config: testAccNodePool_basic_step3(updateName, baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key1", "test_val1"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.1.extend_params.test_key2", "test_val2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNodePoolImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"initial_node_count", "extend_params",
				},
			},
		},
	})
}

func testAccNodePool_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "hcs_cce_cluster" "test" {
  name                   = "%[3]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, common.TestVpc(rName), common.TestVariables(), rName)
}

func testAccNodePool_basic_step1(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.10"
  flavor_id                = "s2.xlarge.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }
}
`, baseConfig, name)
}

func testAccNodePool_basic_step2(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.10"
  flavor_id                = "s2.xlarge.2"
  initial_node_count       = 2
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = true
  min_node_count           = 2
  max_node_count           = 9
  scale_down_cooldown_time = 100
  priority                 = 1
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }

  lifecycle {
    ignore_changes = [
      extend_params
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePool_basic_step3(name, baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.10"
  flavor_id                = "s2.xlarge.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size          = 40
    volumetype    = "SSD"
    extend_params = {
      test_key = "test_val"
    }
  }

  data_volumes {
    size          = 100
    volumetype    = "SSD"
    extend_params = {
      test_key1 = "test_val1"
    }
  }

  data_volumes {
    size          = 100
    volumetype    = "SSD"
    extend_params = {
      test_key2 = "test_val2"
    }
  }

  extend_params {
    docker_base_size = 20
    postinstall      = <<EOF
#! /bin/bash
date
EOF
  }

  lifecycle {
    ignore_changes = [
      extend_params
    ]
  }
}
`, baseConfig, name)
}

func testAccNodePoolImportStateIdFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterId, nodePoolId string
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of CCE node pool is not found in the tfstate", resName)
		}
		clusterId = rs.Primary.Attributes["cluster_id"]
		nodePoolId = rs.Primary.ID
		if clusterId == "" || nodePoolId == "" {
			return "", fmt.Errorf("the CCE node pool ID is not exist or related CCE cluster ID is missing")
		}
		return fmt.Sprintf("%s/%s", clusterId, nodePoolId), nil
	}
}

func TestAccNodePool_tagsLabelsTaints(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_tagsLabelsTaints_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2", "val2"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
				),
			},
			{
				Config: testAccNodePool_tagsLabelsTaints_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test1", "val1_update"),
					resource.TestCheckResourceAttr(resourceName, "labels.test2_update", "val2_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.key", "test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.value", "test_value_update"),
					resource.TestCheckResourceAttr(resourceName, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.key", "new_test_key"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.value", "new_test_value"),
					resource.TestCheckResourceAttr(resourceName, "taints.1.effect", "NoSchedule"),
				),
			},
		},
	})
}

func testAccNodePool_tagsLabelsTaints_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  tags = {
    test1 = "val1"
    test2 = "val2"
  }

  labels = {
    test1 = "val1"
    test2 = "val2"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

}
`, testAccNodePool_base(name), name)
}

func testAccNodePool_tagsLabelsTaints_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  tags = {
    test1        = "val1_update"
    test2_update = "val2_update"
  }

  labels = {
    test1        = "val1_update"
    test2_update = "val2_update"
  }

  taints {
    key    = "test_key"
    value  = "test_value_update"
    effect = "NoSchedule"
  }

  taints {
    key    = "new_test_key"
    value  = "new_test_value"
    effect = "NoSchedule"
  }
}
`, common.TestVariables(), name)
}

func TestAccNodePool_volume_encryption(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_volume_encryption(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "root_volume.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "data_volumes.0.kms_key_id"),
				),
			},
		},
	})
}

func testAccNodePool_volume_encryption(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
    kms_key_id = var.kms_key_id
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = var.kms_key_id
  }
}
`, testAccNodePool_base(name), name)
}

func TestAccNodePool_prePaid(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_prePaid(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "current_node_count", "1"),
				),
			},
		},
	})
}

func testAccNodePool_prePaid(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.availability_zones.test.names[0]
  key_pair                 = hcs_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccNodePool_base(rName), rName)
}

func TestAccNodePool_SecurityGroups(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_SecurityGroups(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.0",
						"hcs_networking_secgroup.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups.1",
						"hcs_networking_secgroup.test.1", "id"),
				),
			},
		},
	})
}

func testAccNodePool_SecurityGroups(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_networking_secgroup_rule" "rule1" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = hcs_vpc_subnet.eni_test[0].cidr
}

resource "hcs_networking_secgroup_rule" "rule2" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "10250"
  protocol          = "tcp"
  remote_ip_prefix  = hcs_vpc_subnet.test.cidr
}

resource "hcs_networking_secgroup_rule" "rule3" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "udp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup_rule" "rule4" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup_rule" "rule5" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup_rule" "rule6" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = hcs_networking_secgroup.test[0].id
}

resource "hcs_networking_secgroup_rule" "rule7" {
  security_group_id = hcs_networking_secgroup.test[0].id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "22"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = data.availability_zones.test.names[0]
  key_pair                 = hcs_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  security_groups = [
    hcs_networking_secgroup.test[0].id,
    hcs_networking_secgroup.test[1].id
  ]

  pod_security_groups = [
    hcs_networking_secgroup.test[2].id,
    hcs_networking_secgroup.test[3].id
  ]

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccCluster_turbo(name, 1), name)
}

func TestAccNodePool_serverGroup(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_serverGroup(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_group_id",
						"hcs_compute_servergroup.test", "id"),
				),
			},
		},
	})
}

func testAccNodePool_serverGroup(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_compute_servergroup" "test" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = var.vailability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  ecs_group_id             = var.server_group_id

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, testAccNodePool_base(rName), rName)
}

func TestAccNodePool_storage(t *testing.T) {
	var (
		nodePool nodepools.NodePool

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "hcs_cce_node_pool.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&nodePool,
			getNodePoolFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodePool_storage(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.selectors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "storage.0.groups.#"),
				),
			},
		},
	})
}

func testAccNodePool_storage(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node_pool" "test" {
  cluster_id               = hcs_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.5"
  flavor_id                = "s6.large.2"
  initial_node_count       = 1
  availability_zone        = var.availability_zone
  key_pair                 = var.keypair_name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = var.kms_key_id
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = var.kms_key_id
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = "1"
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "100"
      match_label_metadata_encrypted = "1"
      match_label_metadata_cmkid     = var.kms_key_id
      match_label_count              = "1"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "kubernetes"
        size        = "10%%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name        = "runtime"
        size        = "90%%"
      }
    }

    groups {
      name           = "vguser"
      selector_names = ["user"]

      virtual_spaces {
        name        = "user"
        size        = "100%%"
        lvm_lv_type = "linear"
        lvm_path    = "/workspace"
      }
    }
  }
}
`, testAccNodePool_base(rName), rName)
}
