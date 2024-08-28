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
  default = "ecs-server"
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

variable "enterprise_project_id" {
  default = "default"
}

variable "kms_key_id" {
  type = string
}

variable "encrypt_cipher" {
  type = string
}