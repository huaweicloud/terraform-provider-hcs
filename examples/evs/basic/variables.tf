variable "volume_configuration" {
  type = list(object({
    suffix      = string
    volume_type = string
    description = string
    size        = number
  }))
  default = [
    {
      suffix      = "volume1",
      volume_type = "SSD",
      description = "volume1_description",
      size        = 1
    },
    {
      suffix      = "volume2",
      volume_type = "SSD",
      description = "volume2_description",
      size        = 2
    },
  ]
}