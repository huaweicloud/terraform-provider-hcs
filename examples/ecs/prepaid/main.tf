data "hcs_availability_zones" "myaz" {}

data "hcs_ecs_compute_flavors" "myflavor" {
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "hcs_ims_images" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "hcs_ecs_compute_instance" "basic" {
  name              = "basic"
  image_id          = data.hcs_ims_images.myimage.id
  flavor_id         = data.hcs_ecs_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  admin_pass        = "******"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1

  network {
    uuid = data.hcs_vpc_subnet.mynet.id
  }
}
