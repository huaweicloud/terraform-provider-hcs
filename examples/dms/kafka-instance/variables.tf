variable "vpc_name" {
  description = "The name of the HuaweiCloudStack VPC"
  default     = "tf_vpc_demo"
}

variable "subnet_name" {
  description = "The name of the HuaweiCloudStack Subnet"
  default     = "tf_subnet_demo"
}

variable "security_group_name" {
  description = "The name of the HuaweiCloudStack Security Group"
  default     = "tf_secgroup_demo"
}

variable "access_user_name" {
  description = "The access user of the Kafka instance"
  default     = "user"
}

variable "manager_user" {
  description = "The manager user of the Kafka instance"
  default     = "kafka-user"
}
