data "hcs_vpc" "myvpc" {
  name = "vpc-default"
}

data "hcs_vpc_eip" "myeip" {
  public_ip = "98.100.16.139"
}

data "hcs_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "hcs_networking_secgroup" "mysecgroup" {
  name = "secgroup-default"
}

resource "hcs_bms_instance" "test" {
  name              = var.instance_name
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  user_id           = var.user_id
  security_groups   = [data.hcs_networking_secgroup.mysecgroup.id]
  availability_zone = "bms_x86"
  vpc_id            = data.hcs_vpc.myvpc.id
  eip_id            = data.hcs_vpc_eip.myeip.id
  admin_pass        = "Huawei12#$"

  data_disks {
    type = "bms_ipsan"
    size = 1
  }

  nics {
    subnet_id  = data.hcs_vpc_subnet.mynet.id
    ip_address = ""
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
