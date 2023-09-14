data "hcs_availability_zones" "myaz" {}

resource "hcs_evs_volume" "myvolume" {
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  count             = length(var.volume_configuration)
  name              = "myvolume_${var.volume_configuration[count.index].suffix}"
  size              = "${var.volume_configuration[count.index].size}"
  description       = "${var.volume_configuration[count.index].description}"
  volume_type       = var.volume_configuration[count.index].volume_type
}