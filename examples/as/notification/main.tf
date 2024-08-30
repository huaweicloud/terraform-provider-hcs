variable "group_name" {
  type = string
  default = "as-group-fb25"
}

data "hcs_as_groups" "groups" {
  name = var.group_name
}

variable "smn_topic_urn" {}

resource "hcs_as_notification" "as_notification" {
  scaling_group_id       = data.hcs_as_groups.groups.groups[0].scaling_group_id
  topic_urn = var.smn_topic_urn
  events = [ "SCALING_UP", "SCALING_UP_FAIL", "SCALING_DOWN", "SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL" ]
}