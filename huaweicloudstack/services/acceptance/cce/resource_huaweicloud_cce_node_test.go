package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/nodes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func getNodeFunc(cfg *config.HcsConfig, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return nodes.Get(client, state.Primary.Attributes["cluster_id"], state.Primary.ID).Extract()
}

func TestAccNode_basic(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HCS_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "s2.xlarge.2"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "ca1.dc1"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.10"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.volumetype", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.volumetype", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENodeImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"tags",
					"password",
				},
			},
			{
				Config: testAccNode_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
		},
	})
}

func testAccNode_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(name), name)
}

func testAccNode_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node" "test" {
  cluster_id        = hcs_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  root_volume {
    volumetype = "SSD"
    size       = 40
  }
  data_volumes {
    volumetype = "SSD"
    size       = 100
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccNode_base(name), name)
}

func TestAccNode_auto_assign_eip(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccSkipUnsupportedTest(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_auto_assign_eip(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)),
				),
			},
		},
	})
}

func testAccNode_auto_assign_eip(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node" "test" {
  cluster_id        = hcs_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  // Assign EIP
  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 100

  root_volume {
    volumetype = "SSD"
    size       = 40
  }
  data_volumes {
    volumetype = "SSD"
    size       = 100
  }
}
`, testAccNode_base(name), name)
}

func TestAccNode_existing_eip(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_existing_eip(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestMatchResourceAttr(resourceName, "public_ip",
						regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)),
				),
			},
		},
	})
}

func testAccNode_existing_eip(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_node" "test" {
  cluster_id        = hcs_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  // Assign existing EIP
  eip_id = "%[3]s"

  root_volume {
    volumetype = "SSD"
    size       = 40
  }
  data_volumes {
    volumetype = "SSD"
    size       = 100
  }
}
`, testAccNode_basic(name), name, acceptance.HCS_EIP_ID)
}

func TestAccNode_volume_extendParams(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_volume_extendParams(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.extend_params.test_key", "test_val"),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.extend_params.test_key", "test_val"),
				),
			},
		},
	})
}

func testAccNode_volume_extendParams(name string) string {
	return fmt.Sprintf(`
resource "hcs_cce_node" "test" {
  cluster_id        = "%[1]s"
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  root_volume {
    size       = 40
	volumetype = "SSD"

	extend_params = {
	  test_key = "test_val"
	}
  }
  data_volumes {
    size       = 100
	volumetype = "SSD"

	extend_params = {
	  test_key = "test_val"
	}
  }
}
`, acceptance.HCS_CCE_CLUSTER_ID, name)
}

func TestAccNode_volume_encryption(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccSkipUnsupportedTest(t)
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckKmsKey(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_volume_encryption(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "root_volume.0.kms_key_id", acceptance.HCS_KMS_KEY_ID),
					resource.TestCheckResourceAttr(resourceName, "data_volumes.0.kms_key_id", acceptance.HCS_KMS_KEY_ID),
				),
			},
		},
	})
}

func testAccNode_volume_encryption(rName string) string {
	return fmt.Sprintf(`
resource "hcs_cce_node" "test" {
  cluster_id        = "%[1]s"
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  root_volume {
    size       = 40
    volumetype = "SSD"
    kms_key_id = "%[3]s"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = "%[3]s"
  }
}
`, acceptance.HCS_CCE_CLUSTER_ID, rName, acceptance.HCS_KMS_KEY_ID)
}

func TestAccNode_storage(t *testing.T) {
	var (
		node nodes.Nodes

		resourceName = "hcs_cce_node.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccSkipUnsupportedTest(t)
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckKms(t)
			acceptance.TestAccPreCheckKmsKey(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNode_storage(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccNode_storage(name string) string {
	return fmt.Sprintf(`
resource "hcs_cce_node" "test" {
  cluster_id        = "%[1]s"
  name              = "%[2]s"
  flavor_id         = "s2.xlarge.2"
  availability_zone = "ca1.dc1"
  password          = "HuaweCloud@Test123"
  os                = "EulerOS 2.10"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = "%[3]s"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = "%[3]s"
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = 1
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "100"
      match_label_metadata_encrypted = "1"
      match_label_metadata_cmkid     = "%[3]s"
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
`, acceptance.HCS_CCE_CLUSTER_ID, name, acceptance.HCS_KMS_KEY_ID)
}

func testAccCCENodeImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		node, ok := s.RootModule().Resources["hcs_cce_node.test"]
		if !ok {
			return "", fmt.Errorf("node not found: %s", node)
		}
		if node.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", acceptance.HCS_CCE_CLUSTER_ID, node.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", acceptance.HCS_CCE_CLUSTER_ID, node.Primary.ID), nil
	}
}
