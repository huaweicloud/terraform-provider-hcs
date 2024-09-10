data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = var.subnet_name
}

data "hcs_ims_images" "test" {
  name       = var.image_name
}

data "hcs_networking_secgroups" "test" {
  name = var.secgroup_name
}

# create ecs instance with encrypt data volume
resource "hcs_ecs_compute_instance" "ecs-test" {
  name                  = join("-", [var.ecs_name, "-data-volume-encrypt"])
  description           = var.ecs_description
  image_id              = data.hcs_ims_images.test.images[0].id
  flavor_id             = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids    = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }
  system_disk_type = var.disk_type
  system_disk_size = var.system_disk_size

  data_disks {
    kms_key_id = var.kms_key_id
    encrypt_cipher = var.encrypt_cipher
    type = var.disk_type
    size = var.data_disk_size
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}

# create ecs instance with encrypt system volume, flavor must be use cloud disk flavor
resource "hcs_ecs_compute_instance" "ecs-test" {
  name                  = join("-", [var.ecs_name, "-sys-volume-encrypt"])
  description           = var.ecs_description
  image_id              = data.hcs_ims_images.test.images[0].id
  flavor_id             = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids    = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }
  system_disk_type = var.disk_type
  system_disk_size = var.system_disk_size
  kms_key_id = var.kms_key_id
  encrypt_cipher = var.encrypt_cipher

  data_disks {
    kms_key_id = var.kms_key_id
    encrypt_cipher = var.encrypt_cipher
    type = var.disk_type
    size = var.data_disk_size
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}