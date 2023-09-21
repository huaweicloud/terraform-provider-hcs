variable "ecs_name" {
  type    = string
  default = "ECS_EP_Example"
}

variable "image_id" {
  type    = string
}

variable "flavor_id" {
  type    = string
}

variable "vpc_name" {
  type    = string
  default = "vpc_EP_Example"
}

variable "vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

variable "subnet_name" {
  type    = string
  default = "subnet_EP_Example"
}

variable "subnet_cidr" {
  type    = string
  default = "192.168.64.0/18"
}

variable "subnet_gateway_ip" {
  type    = string
  default = "192.168.64.1"
}

variable "vpcep_name" {
  type    = string
  default = "vpcep_Example"
}

variable "secgroup_name" {
  type    = string
  default = "secgroup_Example"
}

variable "security_group_rule" {
  type = list(object({
    direction        = string
    ethertype        = string
    protocol         = string
    port_range_min   = number
    port_range_max   = number
    remote_ip_prefix = string
  }))
  default = [
    { direction = "ingress", ethertype = "IPv4", protocol = "tcp", port_range_min = 80, port_range_max = 80, remote_ip_prefix = "0.0.0.0/0" },
    { direction = "ingress", ethertype = "IPv4", protocol = "tcp", port_range_min = 22, port_range_max = 22, remote_ip_prefix = "0.0.0.0/0" },
  ]
}
