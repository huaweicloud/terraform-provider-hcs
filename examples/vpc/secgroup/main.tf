resource "hcs_networking_secgroup" "test" {
  name        = "secgroup-basic"
  description = "basic security group"
}

# allow ping
resource "hcs_networking_secgroup_rule" "allow_ping" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = hcs_networking_secgroup.test.id
}

# allow https
resource "hcs_networking_secgroup_rule" "allow_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = hcs_networking_secgroup.test.id
}
