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

resource "hcs_networking_port" "vip_ass_test_port" {
  name           = var.port_name
  network_id     = hcs_vpc_subnet.test.id
  admin_state_up = "true"
}

resource "hcs_networking_vip" "test" {
  network_id = hcs_vpc_subnet.test.id
}

# associate ports to the vip
resource "hcs_networking_vip_associate" "vip_associated" {
  vip_id   = hcs_networking_vip.test.id
  port_ids = [
    hcs_networking_port.vip_ass_test_port.id
  ]
}
