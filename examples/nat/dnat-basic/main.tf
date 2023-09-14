data "hcs_availability_zones" "newAZ_Example" {}

resource "hcs_ecs_compute_instance" "newCompute_Example" {
  name              = var.ecs_name
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  security_groups   = [hcs_networking_secgroup.newSecgroup_Example.name]
  availability_zone = data.hcs_availability_zones.newAZ_Example.names[0]

  network {
    uuid = hcs_vpc_subnet.newSubnet_Example.id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = var.image_id
    volume_type = "business_type_01"
    volume_size = 20
  }
}

resource "hcs_vpc_eip" "newEIP_Example" {
  publicip {
    type = var.eip_external_network_name
  }
  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER"
  }
}

resource "hcs_vpc" "newVPC_Example" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_vpc_subnet" "newSubnet_Example" {
  name          = var.subnet_name
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  vpc_id        = hcs_vpc.newVPC_Example.id
  primary_dns   = "100.125.129.250"
  secondary_dns = "100.125.1.250"
}

resource "hcs_networking_secgroup" "newSecgroup_Example" {
  name        = var.secgroup_name
  description = "This is a security group"
}

resource "hcs_networking_secgroup_rule" "newSecgroupRule_Example" {
  count = length(var.security_group_rule)

  direction         = lookup(var.security_group_rule[count.index], "direction", null)
  ethertype         = lookup(var.security_group_rule[count.index], "ethertype", null)
  protocol          = lookup(var.security_group_rule[count.index], "protocol", null)
  port_range_min    = lookup(var.security_group_rule[count.index], "port_range_min", null)
  port_range_max    = lookup(var.security_group_rule[count.index], "port_range_max", null)
  remote_ip_prefix  = lookup(var.security_group_rule[count.index], "remote_ip_prefix", null)
  security_group_id = hcs_networking_secgroup.newSecgroup_Example.id
}

resource "hcs_nat_gateway" "newNet_gateway_Example" {
  name                = var.net_gateway_name
  description         = "example for net test"
  spec                = "1"
  vpc_id              = hcs_vpc.newVPC_Example.id
  subnet_id = hcs_vpc_subnet.newSubnet_Example.id
}

resource "hcs_nat_dnat_rule" "newDNATRule_Example" {
  count = length(var.example_dnat_rule)

  floating_ip_id = hcs_vpc_eip.newEIP_Example.id
  nat_gateway_id = hcs_nat_gateway.newNet_gateway_Example.id
  port_id        = hcs_ecs_compute_instance.newCompute_Example.network[0].port

  internal_service_port = lookup(var.example_dnat_rule[count.index], "internal_service_port", null)
  protocol              = lookup(var.example_dnat_rule[count.index], "protocol", null)
  external_service_port = lookup(var.example_dnat_rule[count.index], "external_service_port", null)
}
