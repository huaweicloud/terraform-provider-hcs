variable "volume_name" {
  type    = string
  default = "volume-name"
}

variable "volume_size" {
  type    = number
  default = 1
}

variable "volume_type" {
  type    = string
  default = "SSD"
}

variable "volume_description" {
  type    = string
  default = "volume-description"
}

variable "snapshot_name" {
  type    = string
  default = "snapshot-name"
}

variable "snapshot_description" {
  type    = string
  default = "snapshot-description"
}