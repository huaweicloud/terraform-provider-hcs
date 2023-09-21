resource "hcs_ims_image_share" "test" {
  source_image_id    = var.source_image_id
  target_project_ids = var.target_project_ids
}
