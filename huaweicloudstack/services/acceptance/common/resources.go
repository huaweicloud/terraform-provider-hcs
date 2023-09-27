package common

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance"
)

// TestSecGroup can be referred as `hcs_networking_secgroup.test`
func TestSecGroup(name string) string {
	return fmt.Sprintf(`
resource "hcs_networking_secgroup" "test" {
  name                 = "%s"
  delete_default_rules = true
}
`, name)
}

// TestVpc can be referred as `hcs_vpc.test` and `hcs_vpc_subnet.test`
func TestVpc(name string) string {
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
`, name)
}

// TestBaseNetwork vpc, subnet, security group
func TestBaseNetwork(name string) string {
	return fmt.Sprintf(`
# base security group without default rules
%s

# base vpc and subnet
%s
`, TestSecGroup(name), TestVpc(name))
}

// TestBaseComputeResources vpc, subnet, security group, availability zone, keypair, image, flavor
func TestBaseComputeResources(name string) string {
	return fmt.Sprintf(`
# base test resources
%s

data "hcs_availability_zones" "test" {}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_ims_images" "test" {
  name        = "ecs_mini_image"
}
`, TestBaseNetwork(name))
}

// TestVariables can be referred as `hcs_vpc.test` and `hcs_vpc_subnet.test`
func TestVariables() string {
	return fmt.Sprintf(`
variable "availability_zone" {
  type    = string
  default = "%s"
}

variable "keypair_name" {
  type    = string
  default = "%s"
}

variable "kms_key_id" {
  type    = string
  default = "%s"
}

variable "server_group_id" {
  type    = string
  default = "%s"
}

variable "ecs_instance_id" {
  type    = string
  default = "%s"
}

variable "cluster_id" {
  type    = string
  default = "%s"
}
`, acceptance.HCS_AVAILABILITY_ZONE, acceptance.HCS_KEYPAIR_NAME,
		acceptance.HCS_KMS_KEY_ID, acceptance.HCS_SERVER_GROUP_ID, acceptance.HCS_ECS_INSTANCE_ID, acceptance.HCS_CCE_CLUSTER_ID)
}
