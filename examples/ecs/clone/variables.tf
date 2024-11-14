variable "vpc_id" {
  default = "ba69be72-406b-4d3a-b552-5f93fb47c1d9"
}

variable "instance_id" {
  default = "c9a6edf4-ba26-48cb-857f-9aff9ffa25dd"
}

variable "ecs_name" {
  default = "ecs-clone-test"
}

variable "power_on" {
  type = bool
  default = false
}

variable "retain_passwd" {
  type = bool
  default = false
}
