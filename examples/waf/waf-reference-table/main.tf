# create vpc and subnet
resource "hcs_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = "192.168.0.0/24"
}

resource "hcs_vpc_subnet" "vpc_subnet_1" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc_1.id
}

# create security group
resource "hcs_networking_secgroup" "secgroup" {
  name        = var.security_group_name
  description = "terraform security group"
}

data "hcs_availability_zones" "zones" {}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.zones.names[1]
  performance_type  = "normal"
  cpu_core_count    = 2
}

# create a waf dedicated instance
resource "hcs_waf_dedicated_instance" "instance_1" {
  name               = var.waf_dedicated_instance_name
  available_zone     = data.hcs_availability_zones.zones.names[1]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.hcs_ecs_compute_flavors.flavors.ids[0]
  vpc_id             = hcs_vpc.vpc_1.id
  subnet_id          = hcs_vpc_subnet.vpc_subnet_1.id

  security_group = [
    hcs_networking_secgroup.secgroup.id
  ]
}

resource "hcs_waf_reference_table" "ref_table" {
  name        = var.waf_table_name
  type        = "url"
  description = "tf acc"

  conditions = [
    "/admin",
    "/manage"
  ]

  depends_on = [
    hcs_waf_dedicated_instance.instance_1
  ]
}
