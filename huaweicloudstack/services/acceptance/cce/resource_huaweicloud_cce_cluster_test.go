package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/clusters"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func TestAccCluster_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_users.0.client_certificate_data", "kube_config_raw",
				},
			},
			{
				Config: testAccCluster_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "new description"),
				),
			},
		},
	})
}

func TestAccCluster_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEipAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip", "certificate_users.0.client_certificate_data", "kube_config_raw",
				},
			},
		},
	})
}

// Untested
func TestAccCluster_turbo(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_turbo(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
			{
				Config: testAccCluster_turbo(rName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "eni"),
					resource.TestCheckOutput("is_eni_subnet_id_different", "false"),
				),
			},
		},
	})
}

func TestAccCluster_hibernate(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
			{
				Config: testAccCluster_hibernate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Hibernation"),
				),
			},
			{
				Config: testAccCluster_awake(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
				),
			},
		},
	})
}

func TestAccCluster_multiContainerNetworkCidrs(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_multiContainerNetworkCidrs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "vpc-router"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "container_network_cidr", "172.16.0.0/24,172.16.1.0/24"),
				),
			},
		},
	})
}

func TestAccCluster_secGroup(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "hcs_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_secGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
					resource.TestCheckResourceAttr(resourceName, "service_network_cidr", "10.248.0.0/16"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "hcs_networking_secgroup", "id"),
				),
			},
			{
				Config: testAccCluster_secGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "hcs_networking_secgroup", "id"),
				),
			},
		},
	})
}

func testAccCheckClusterDestroy(s *terraform.State) error {
	config := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
	cceClient, err := config.CceV3Client(acceptance.HCS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CCE v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcs_cce_cluster" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

func testAccCheckClusterExists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		config := config.GetHcsConfig(acceptance.TestAccProvider.Meta())
		cceClient, err := config.CceV3Client(acceptance.HCS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CCE v3 client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCluster_basic(rName string) string {
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
`, common.TestVpc(rName), rName)
}

func testAccCluster_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "new description"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_withEip(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = "%[3]s"
}
`, common.TestVpc(rName), rName, acceptance.HCS_EIP_ADDRESS)
}

func testAccCluster_turbo(rName string, eniNum int) string {
	return fmt.Sprintf(`
%s

resource "hcs_vpc_subnet" "eni_test" {
  count      = %[3]d

  name       = "%[2]s-eni-${count.index}"
  cidr       = cidrsubnet(hcs_vpc.test.cidr, 8, count.index + 1)
  gateway_ip = cidrhost(cidrsubnet(hcs_vpc.test.cidr, 8, count.index + 1), 1)
  vpc_id     = hcs_vpc.test.id
}

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "eni"
  eni_subnet_id          = join(",", hcs_vpc_subnet.eni_test[*].ipv4_subnet_id)
}

output "is_eni_subnet_id_different" {
  value = length(setsubtract(split(",", hcs_cce_cluster.test.eni_subnet_id),
  hcs_vpc_subnet.eni_test[*].ipv4_subnet_id)) != 0
}
`, common.TestVpc(rName), rName, eniNum)
}

func testAccCluster_hibernate(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_awake(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  hibernate              = false

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_multiContainerNetworkCidrs(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "vpc-router"
  container_network_cidr = "172.16.0.0/24,172.16.1.0/24"
  service_network_cidr   = "10.248.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_secGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup" "test" {
  name        = "secgroup_1"
  description = "My security group"
}

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform"
  security_group_id      = hcs_networking_secgroup.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccCluster_secGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcs_networking_secgroup" "test2" {
  name        = "secgroup_1"
  description = "My security group"
}

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  service_network_cidr   = "10.248.0.0/16"
  description            = "created by terraform update"
  security_group_id      = hcs_networking_secgroup.test2.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(rName), rName)
}
