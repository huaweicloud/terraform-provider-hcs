variable "hook_name" {
  type = string
  default = "as-policy-77e3"
}

data "hcs_smn_topics" "tops" {
  name = "topic_1"
}

variable "smn_topic_urn" {}

resource "hcs_as_lifecycle_hook" "lifecycle_hook1" {
  scaling_group_id       = data.hcs_as_groups.groups.groups[0].scaling_group_id
  name                   = var.hook_name
  type                   = "ADD"
  default_result         = "ABANDON"
  notification_topic_urn = var.smn_topic_urn
  notification_message   = "This is a test message"
}