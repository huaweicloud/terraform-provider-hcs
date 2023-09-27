data "hcs_ecs_compute_instance" "test" {
  name = var.instance_name
}

resource "hcs_ims_image" "test" {
  name        = var.image_name
  instance_id = data.hcs_ecs_compute_instance.test.id
  description = "created by Terraform"
}

resource "hcs_ims_image" "ims_test_file" {
  name        = "ims_test_file"
  image_url   = "bucket:test.qcow2"
  min_disk    = 10
  os_version  = "Other(64 bit) 64bit"
  description = "Create an image from the OBS bucket."
}
