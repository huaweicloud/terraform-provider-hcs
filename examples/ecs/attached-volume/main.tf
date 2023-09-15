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

resource "hcs_evs_volume" "myvolume" {
  name              = "myvolume"
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  volume_type       = "SAS"
  size              = 10
}

resource "hcs_ecs_compute_volume_attach" "attached" {
  instance_id = hcs_ecs_compute_instance.myinstance.id
  volume_id   = hcs_evs_volume.myvolume.id
}
