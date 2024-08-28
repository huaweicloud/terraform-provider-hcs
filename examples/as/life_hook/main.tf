variable "hook_name" {
  type = string
  default = "as-policy-77e3"
}

variable "smn_topic_urn" {}

data "hcs_smn_topics" "tops" {
  name = "topic_1"
}

variable "group_name" {
  type = string
  default = "as-group-fb25"
}

data "hcs_as_groups" "groups" {
  name = var.group_name
}

resource "hcs_as_lifecycle_hook" "lifecycle_hook1" {
  scaling_group_id       = data.hcs_as_groups.groups.groups[0].scaling_group_id
  name                   = var.hook_name
  type                   = "ADD"
  default_result         = "ABANDON"
  notification_topic_urn = var.smn_topic_urn
  notification_message   = "This is a test message"
}