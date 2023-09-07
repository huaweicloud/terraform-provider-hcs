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
  container_network_type = "overlay_l2"
  cluster_type           = "ARM64"
  eip                    = hcs_vpc_eip.cce.address
}


resource "hcs_vpc_eip" "cce" {
  publicip {
    type = "eip_external_net"
  }
  bandwidth {
    name        = "cce-apiserver"
    size        = 20
    share_type  = "PER"
  }
}

resource "hcs_cce_node" "cce-node" {
  cluster_id        = hcs_cce_cluster.mycce.id
  name              = "node"
  flavor_id         = "rc6.large.2"
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

provider "kubernetes" {
  host                   = "https://${hcs_vpc_eip.cce.address}:5443"
  cluster_ca_certificate = base64decode(hcs_cce_cluster.mycce.certificate_clusters[0].certificate_authority_data)
  client_certificate     = base64decode(hcs_cce_cluster.mycce.certificate_users[0].client_certificate_data)
  client_key             = base64decode(hcs_cce_cluster.mycce.certificate_users[0].client_key_data)
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "my-first-namespace"
  }
}
