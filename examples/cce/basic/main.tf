resource "hcs_vpc" "myvpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_vpc_subnet" "mysubnet" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway

  # dns is required for cce node installing
  primary_dns   = var.primary_dns
  secondary_dns = var.secondary_dns
  vpc_id        = hcs_vpc.myvpc.id
}

data "hcs_availability_zones" "myaz" {}

resource "hcs_ecs_compute_keypair" "mykeypair" {
  name = var.key_pair_name
}
resource "hcs_cce_cluster" "mycce" {
  name                   = var.cce_cluster_name
  flavor_id              = var.cce_cluster_flavor
  vpc_id                 = hcs_vpc.myvpc.id
  subnet_id              = hcs_vpc_subnet.mysubnet.id
  cluster_type           = "ARM64"
  container_network_type = "overlay_l2"
}

resource "hcs_cce_node" "mynode" {
  cluster_id        = hcs_cce_cluster.mycce.id
  name              = var.node_name
  flavor_id         = var.node_flavor
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  key_pair          = hcs_ecs_compute_keypair.mykeypair.name

  root_volume {
    size       = var.root_volume_size
    volumetype = var.root_volume_type
  }
  data_volumes {
    size       = var.data_volume_size
    volumetype = var.data_volume_type
  }
}