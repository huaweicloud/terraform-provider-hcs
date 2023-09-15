data "hcs_availability_zones" "default" {}

data "hcs_ims_images" "default" {
  name        = var.image_name
  most_recent = true
}

data "hcs_ecs_compute_flavors" "default" {
  availability_zone = data.hcs_availability_zones.default.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_keypair" "default" {
  name     = var.keypair_name
  key_file = var.private_key_path
}

resource "hcs_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_vpc_subnet" "default" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  vpc_id     = hcs_vpc.default.id
  gateway_ip = var.gateway_ip
}

resource "hcs_networking_secgroup" "default" {
  name = var.security_group_name
}

resource "hcs_vpc_eip" "default" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcs_ecs_compute_instance" "default" {
  name              = var.ecs_instance_name
  image_id          = data.hcs_ims_images.default.id
  flavor_id         = data.hcs_ecs_compute_flavors.default.ids[0]
  availability_zone = data.hcs_availability_zones.default.names[0]
  key_pair          = hcs_ecs_compute_keypair.default.name
  user_data         = <<-EOF
#!/bin/bash
echo '${file("./test.txt")}' > /home/test.txt
EOF

  system_disk_type = "SAS"
  system_disk_size = 50

  security_groups = [
    hcs_networking_secgroup.default.name
  ]

  network {
    uuid = hcs_vpc_subnet.default.id
  }
}

resource "hcs_ecs_compute_eip_associate" "default" {
  public_ip   = hcs_vpc_eip.default.address
  instance_id = hcs_ecs_compute_instance.default.id
}

resource "null_resource" "provision" {
  depends_on = [hcs_ecs_compute_eip_associate.default]

  provisioner "remote-exec" {
    connection {
      user        = "root"
      private_key = file(var.private_key_path)
      host        = hcs_vpc_eip.default.address
    }

    inline = [
      "cat /home/test.txt"
    ]
  }
}
