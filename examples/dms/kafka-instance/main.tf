# random password
resource "random_password" "access_user_password" {
  length           = 16
  special          = true
  override_special = "!@#$%*"
}

# random password
resource "random_password" "manager_password" {
  length           = 16
  special          = true
  override_special = "!@#$%*"
}

# Create a VPC.
resource "hcs_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = "192.168.0.0/24"
}

# Create a subnet under the VPC that created above.
resource "hcs_vpc_subnet" "vpc_subnet_1" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc_1.id
}

# Create a security group.
resource "hcs_networking_secgroup" "secgroup" {
  name        = var.security_group_name
  description = "terraform security group"
}

# List the availability zones in the current region.
data "hcs_availability_zones" "zones" {}

# Query flavor information based on flavorID and storage I/O specification.
# Make sure the flavors are available in the availability zone.
data "hcs_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"

  availability_zones = [
    data.hcs_availability_zones.zones.names[0],
    data.hcs_availability_zones.zones.names[1],
    data.hcs_availability_zones.zones.names[2],
  ]
  storage_spec_code = "dms.physical.storage.ultra.v2"
}

# Create the DMS Kafka instance.
resource "hcs_dms_kafka_instance" "kafka_instance_1" {
  name        = "instance_1"
  description = "kafka instance demo"

  availability_zones = [
    data.hcs_availability_zones.zones.names[0],
    data.hcs_availability_zones.zones.names[1],
    data.hcs_availability_zones.zones.names[2],
  ]

  engine_version    = "2.7"
  flavor_id         = data.hcs_dms_kafka_flavors.test.flavor_id
  storage_spec_code = data.hcs_dms_kafka_flavors.test.flavors[0].ios[0].storage_spec_code
  storage_space     = 600
  broker_num        = 3

  vpc_id            = hcs_vpc.vpc_1.id
  network_id        = hcs_vpc_subnet.vpc_subnet_1.id
  security_group_id = hcs_networking_secgroup.secgroup.id

  access_user      = var.access_user_name
  password         = random_password.access_user_password
  manager_user     = var.manager_user
  manager_password = random_password.manager_password
}
