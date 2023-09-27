variable "vpc_name" {
  default = "vpc-default"
}

variable "subnet_name" {
  default = "subnet-32a8"
}

variable "secgroup_name" {
  default = "default"
}

variable "image_name" {
  default = "cirros-arm"
}

variable "ecs_name" {
  default = "ecs-associated-eip"
}

variable "ecs_description" {
  default = ""
}

variable "disk_type" {
  default = "business_type_01"
}

variable "volume_description" {
  default = ""
}

variable "system_disk_size" {
  type = number
  default = 10
}

variable "data_disk_size" {
  type = number
  default = 10
}

variable "epi_type" {
  default = "eip"
}

variable "bandwidth_name" {
  default = "test"
}

variable "bandwidth_size" {
  type = number
  default = 100
}

variable "bandwidth_share_type" {
  default = "PER"
}

