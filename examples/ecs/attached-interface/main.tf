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
}

resource "hcs_ecs_compute_instance" "myinstance" {
  name              = "basic"
  image_id          = data.hcs_ims_images.myimage.id
  flavor_id         = data.hcs_ecs_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.hcs_availability_zones.myaz.names[0]

  network {
    uuid = data.hcs_vpc_subnet.mynet.id
  }
}

data "hcs_vpc" "myvpc" {
  name = "vpc-default"
}

resource "hcs_vpc_subnet" "attach" {
  name       = "subnet-attach"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = data.hcs_vpc.myvpc.id

  availability_zone = data.hcs_availability_zones.myaz.names[0]
}

resource "hcs_ecs_compute_interface_attach" "attached" {
  instance_id = hcs_ecs_compute_instance.myinstance.id
  network_id  = hcs_vpc_subnet.attach.id

  # This is optional
  fixed_ip = "192.168.1.100"
}
