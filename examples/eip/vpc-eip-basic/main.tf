resource "hcs_vpc_eip" "eip_1" {
  publicip {
    type = var.eip_external_network_name
  }
  bandwidth {
    name        = var.eip_name
    size        = var.bandwidth_size
    share_type  = "PER"
  }
}
