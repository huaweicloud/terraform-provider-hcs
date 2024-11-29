resource "hcs_ecs_compute_instance_clone" "test" {
  instance_id         = var.instance_id
  power_on          = var.power_on
  name = var.ecs_name
  vpc_id = var.vpc_id
  retain_passwd = var.retain_passwd
  admin_pass = "test_admin_pass"

  network {
    subnet_id = "b4db86e5-406e-422b-8744-54a15beb0aba"
    fixed_ip_v4 = "192.167.0.10"
    ipv6_enable = false
    fixed_ip_v6 = "1030::C9B4:FF12:48AA:1A2B"
    security_group_ids {
      security_group_id = "3e8cdc50-b6e3-4f7a-b7b7-88fe18a0e6a0"
    }
  }
}