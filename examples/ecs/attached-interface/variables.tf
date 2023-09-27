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
  default = "ecs-attached-interface"
}

variable "ecs_description" {
  default = ""
}

variable "disk_type" {
  default = "business_type_01"
}

variable "system_disk_size" {
  type = number
  default = 10
}

variable "data_disk_size" {
  type = number
  default = 10
}

variable "subnet_attach_name" {
  default = "subnet-attach"
}

variable "subnet_cidr" {
  default = "192.168.1.0/24"
}

variable "subnet_gateway" {
  default = "192.168.1.1"
}

variable "ecs_fix_ip" {
  default = "192.168.1.100"
}