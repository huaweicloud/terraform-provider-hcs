resource "hcs_vpc_bandwidth" "test" {
  name = var.bandwidth_name
  size = var.bandwidth_size
}
