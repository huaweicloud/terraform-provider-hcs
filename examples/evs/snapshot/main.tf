data "hcs_availability_zones" "myaz" {}

resource "hcs_evs_volume" "myvolume" {
  availability_zone = data.hcs_availability_zones.myaz.names[0]
  name              = var.volume_name
  description       = var.volume_description
  volume_type       = var.volume_type
  size              = var.volume_size
}

resource "hcs_evs_snapshot" "mysnapshot" {
  volume_id   = hcs_evs_volume.myvolume.id
  name        = var.snapshot_name
  description = var.snapshot_description
}