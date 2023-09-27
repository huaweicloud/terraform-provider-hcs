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
  default = "ecs-attached-volume"
}

variable "ecs_description" {
  default = ""
}

variable "disk_type" {
  default = "business_type_01"
}

variable "volume_name" {
  default = "ecs-data-disk"
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

variable "attach_disk_size" {
  type = number
  default = 10
}

variable "volume_device" {
  default = "/dev/vdb"
}