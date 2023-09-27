resource "hcs_vpc" "vpc_1" {
  name = "vpc-te"
}

resource "hcs_vpc_subnet" "subnet_1" {
  vpc_id      = hcs_vpc.vpc_1.id
  name        = "sub-te"
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
}

resource "hcs_elb_loadbalancer" "elb_1" {
  name          = "elb-te"
  ipv4_subnet_id = hcs_vpc_subnet.subnet_1.ipv4_subnet_id
}

resource "hcs_elb_listener" "listener_1" {
  name            = "listener_http"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = hcs_elb_loadbalancer.elb_1.id
}

resource "hcs_elb_pool" "group_1" {
  name        = "group_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcs_elb_listener.listener_1.id
}

resource "hcs_elb_monitor" "health_check" {
  protocol    = "HTTP"
  interval    = 30
  timeout     = 15
  max_retries = 10
  url_path    = "/api"
  port        = 8888
  pool_id        = hcs_elb_pool.group_1.id
}

resource "hcs_elb_member" "member_1" {
  address       = var.member_ip
  protocol_port = 80
  weight        = 1
  pool_id       = hcs_elb_pool.group_1.id
  subnet_id     = hcs_vpc_subnet.subnet_1.ipv4_subnet_id
}