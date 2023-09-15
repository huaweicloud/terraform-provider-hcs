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

resource "hcs_ecs_compute_keypair" "mykey" {
  name       = "terraform-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDy+49hbB9Ni2SttHcbJU+ngQXUhiGDVsflp2g5A3tPrBXq46kmm/nZv9JQqxlRzqtFi9eTI7OBvn2A34Y+KCfiIQwtgZQ9LF5ROKYsGkS2o9ewsX8Hghx1r0u5G3wvcwZWNctgEOapXMD0JEJZdNHCDSK8yr+btR4R8Ypg0uN+Zp0SyYX1iLif7saiBjz0zmRMmw5ctAskQZmCf/W5v/VH60fYPrBU8lJq5Pu+eizhou7nFFDxXofr2ySF8k/yuA9OnJdVF9Fbf85Z59CWNZBvcTMaAH2ALXFzPCFyCncTJtc/OVMRcxjUWU1dkBhOGQ/UnhHKcflmrtQn04eO8xDr root@terra-dev"
}

resource "hcs_ecs_compute_instance" "basic" {
  name            = "basic"
  image_id        = data.hcs_ims_images.myimage.id
  flavor_id       = data.hcs_ecs_compute_flavors.myflavor.ids[0]
  security_groups = ["default"]

  # NOTE: admin_pass doesn't work with user_data, use key_pair instead.
  key_pair          = hcs_ecs_compute_keypair.mykey.name
  availability_zone = data.hcs_availability_zones.myaz.names[0]

  # NOTE: can also use file("userdata.sh") to fetch the content.
  user_data = "#!/bin/bash\necho hello > /home/terraform.txt"

  network {
    uuid = data.hcs_vpc_subnet.mynet.id
  }
}
