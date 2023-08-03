package cce

import (
	"fmt"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/acceptance/common"
)

func testAccCceCluster_config(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "hcs_availability_zones" "test" {}

resource "hcs_compute_keypair" "test" {
  name = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "hcs_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = hcs_vpc.test.id
  subnet_id              = hcs_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}

resource "hcs_cce_node" "test" {
  cluster_id        = hcs_cce_cluster.test.id
  name              = "%[2]s"
  flavor_id         = "s6.large.2"
  availability_zone = data.hcs_availability_zones.test.names[0]
  key_pair          = hcs_compute_keypair.test.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}`, common.TestBaseNetwork(rName), rName)
}
