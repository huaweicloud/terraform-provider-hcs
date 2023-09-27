variable "firewall_instance_id" {
  type        = string
  description = "The firewall instance ID."
}

variable "protected_eip_id" {
  type        = string
  description = "The ID of the protected EIP."
}

variable "protected_eip_address" {
  type        = string
  description = "The IPv4 address of the protected EIP."
}

variable "protection_rule_name" {
  type        = string
  description = "The protection rule name."
}
