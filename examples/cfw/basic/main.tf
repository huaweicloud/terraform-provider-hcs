data "hcs_cfw_firewalls" "test" {
  instance_id= var.firewall_instance_id
}

resource "hcs_cfw_eip_protection" "test" {
  object_id = data.hcs_cfw_firewalls.test.records[0].protect_objects[0].object_id

  protected_eip {
    id          = var.protected_eip_id
    public_ipv4 = var.protected_eip_address
  }
}

resource "hcs_cfw_protection_rule" "test" {
  depends_on = [hcs_cfw_eip_protection.test]

  name                = var.protection_rule_name
  object_id           = data.hcs_cfw_firewalls.test.records[0].protect_objects[0].object_id
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "1.1.1.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
