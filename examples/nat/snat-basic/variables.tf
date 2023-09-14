variable "nat_gatway_name" {
  default = "nat-gateway-basic"
}

variable "vpc_name" {
  default = "vpc-basic"
}

variable "vpc_cidr" {
  default = "172.16.0.0/16"
}

variable "subnet_name" {
  default = "subent-basic"
}

variable "subnet_cidr" {
  default = "172.16.10.0/24"
}

variable "subnet_gateway" {
  default = "172.16.10.1"
}

variable "primary_dns" {
  default = "100.125.1.250"
}

variable "eip_external_network_name" {
  default = "eip"
}
