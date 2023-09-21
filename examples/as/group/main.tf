data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_ims_images" "test" {
  name        = "ecs_mini_image"
}

resource "hcs_networking_secgroup" "test" {
  name                 = "secgroup_3196"
  delete_default_rules = true
}

resource "hcs_vpc" "test" {
  name = "tf_test_3200"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "test" {
  name       = "tf_test_3200"
  vpc_id     = hcs_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "hcs_ecs_compute_keypair" "acc_key" {
  name       = "tf_test_3200"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "hcs_elb_loadbalancer" "loadbalancer_1" {
  name          = "tf_test_3200"
  ipv4_subnet_id = hcs_vpc_subnet.test.ipv4_subnet_id
}

resource "hcs_elb_listener" "listener_1" {
  name            = "tf_test_3200"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = hcs_elb_loadbalancer.loadbalancer_1.id
}

resource "hcs_elb_pool" "pool_1" {
  name        = "tf_test_3200"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcs_elb_listener.listener_1.id
}

resource "hcs_as_configuration" "acc_as_config"{
  scaling_configuration_name = "tf_test_3200"
  instance_config {
    image    = data.hcs_ims_images.test.images[0].id
    flavor   = data.hcs_ecs_compute_flavors.test.ids[0]
    key_name = hcs_ecs_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "business_type_01"
      disk_type   = "SYS"
    }
  }
}

resource "hcs_as_group" "acc_as_group"{
  scaling_group_name       = "tf_test_3200"
  scaling_configuration_id = hcs_as_configuration.acc_as_config.id
  vpc_id                   = hcs_vpc.test.id
  max_instance_number      = 5

  networks {
    id = hcs_vpc_subnet.test.id
  }
  security_groups {
    id = hcs_networking_secgroup.test.id
  }
  lbaas_listeners {
    listener_id   = hcs_elb_listener.listener_1.id
    pool_id       = hcs_elb_pool.pool_1.id
    protocol_port = hcs_elb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}