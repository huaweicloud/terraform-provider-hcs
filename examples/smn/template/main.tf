resource "hcs_smn_message_template" "template_1" {
  name         = "smn_template"
  protocol      = "email"
  content      = "template_content"
}