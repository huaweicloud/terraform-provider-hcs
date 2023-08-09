data "hcs_availability_zones" "myaz" {}

data "hcs_vpc" "myvpc" {
  id = var.vpc_id
}

data "hcs_vpc_subnet" "mysubnet" {
  id = var.subnet_id
}

resource "hcs_vpc_eip" "cce" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "cce-apiserver"
    size        = 20
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcs_cce_cluster" "cluster" {
  name                   = var.cce_name
  cluster_type           = "VirtualMachine"
  cluster_version        = "v1.19"
  flavor_id              = "cce.s1.small"
  vpc_id                 = data.hcs_vpc.myvpc.id
  subnet_id              = data.hcs_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = hcs_vpc_eip.cce.address
  delete_all             = "true"
}

resource "hcs_cce_node" "cce-node1" {
  cluster_id        = hcs_cce_cluster.cluster.id
  name              = "node1"
  flavor_id         = "s6.large.2"
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "hcs_cce_node" "cce-node2" {
  cluster_id        = hcs_cce_cluster.cluster.id
  name              = "node2"
  flavor_id         = "s6.large.2"
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "local_file" "kube_config" {
    content = hcs_cce_cluster.cluster.kube_config_raw
    filename = " ~/.kube/config"
}

provider "kubernetes" {
    config_path    = local_file.kube_config.filename
    config_context = "external"
}
