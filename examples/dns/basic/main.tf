resource "hcs_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_dns_zone" "test" {
  name        = var.zone_name
  description = "a zone"
  ttl         = 300
  zone_type             = "private"

  router {
    router_id = hcs_vpc.test.id
  }
}

resource "hcs_dns_recordset" "test" {
  zone_id     = hcs_dns_zone.test.id
  name        = var.record_name
  type        = "TXT"
  description = "an updated record set"
  status      = "ENABLE"
  ttl         = 6000
  records     = ["\"test records\""]
}