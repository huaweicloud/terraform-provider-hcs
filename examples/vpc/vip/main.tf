data "hcs_availability_zones" "myaz" {}

data "hcs_ecs_compute_instance" "hcs_instance" {
  name = var.ecs_name
}

resource "hcs_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_vpc_subnet" "test" {
  vpc_id      = hcs_vpc.test.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "hcs_networking_vip" "test" {
  network_id = hcs_vpc_subnet.test.id
}

# associate ports to the vip
resource "hcs_networking_vip_associate" "vip_associated" {
  vip_id   = hcs_networking_vip.test.id
  port_ids = [
    data.hcs_ecs_compute_instance.hcs_instance.network[0].port
  ]
}
