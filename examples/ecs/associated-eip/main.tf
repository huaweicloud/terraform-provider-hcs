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

resource "hcs_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "mybandwidth"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcs_ecs_compute_eip_associate" "associated" {
  public_ip   = hcs_vpc_eip.myeip.address
  instance_id = hcs_ecs_compute_instance.myinstance.id
}
