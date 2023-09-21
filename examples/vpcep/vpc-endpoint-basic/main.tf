data "hcs_availability_zones" "test" {
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

resource "hcs_ecs_compute_instance" "newCompute_Example" {
  name              = var.ecs_name
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  security_groups   = [hcs_networking_secgroup.newSecgroup_Example.name]
  availability_zone = data.hcs_availability_zones.test.names[0]

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

resource "hcs_networking_secgroup" "newSecgroup_Example" {
  name        = var.secgroup_name
  description = "This is a security group"
}

resource "hcs_vpcep_service" "test" {
  name        = var.vpcep_name
  server_type = "VM"
  vpc_id      = hcs_vpc.newVPC_Example.id
  port_id     = hcs_ecs_compute_instance.newCompute_Example.network[0].port
  approval    = false

  port_mapping {
    service_port  = 1111
    terminal_port = 1111
  }
}

resource "hcs_vpcep_endpoint" "test" {
  service_id  = hcs_vpcep_service.test.id
  vpc_id      = hcs_vpc.newVPC_Example.id
  network_id  = hcs_vpc_subnet.newSubnet_Example.id
  enable_dns  = false
}