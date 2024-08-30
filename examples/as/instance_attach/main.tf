variable "ecs_id" {
  type = string
  default = "1833f1d0-9250-4054-bc30-8f6bd7469b95"
}

variable "group_name" {
  type = string
  default = "as-group-fb25"
}

data "hcs_as_groups" "groups" {
  name = var.group_name
}

resource "hcs_as_instance_attach" "as_instance1" {
  scaling_group_id = data.hcs_as_groups.groups.groups[0].scaling_group_id
  instance_id      = var.ecs_id
}