---
subcategory: "Network ACL"
---

# hcs_network_acl

Manages a network ACL resource within HuaweiCloudStack.

## Example Usage

```hcl
data "hcs_vpc_subnet" "subnet" {
  name = "subnet-default"
}

resource "hcs_network_acl_rule" "rule_1" {
  name             = "my-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_ports= ["23"]
  enabled          = "true"
}

resource "hcs_network_acl_rule" "rule_2" {
  name             = "my-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_ports= ["123"]
  enabled          = "false"
}

resource "hcs_network_acl" "fw_acl" {
  name          = "my-fw-acl"
  subnets       = [data.hcs_vpc_subnet.subnet.id]
  inbound_rules = [
    hcs_network_acl_rule.rule_1.id,
    hcs_network_acl_rule.rule_2.id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the network acl resource. If omitted, the
  provider-level region will be used. Changing this creates a new network acl resource.

* `name` - (Required, String) Specifies the network ACL name. This parameter can contain a maximum of 64 characters,
  which may consist of letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional, String) Specifies the supplementary information about the network ACL. This parameter can
  contain a maximum of 255 characters and cannot contain angle brackets (< or >).

* `inbound_rules` - (Optional, List)  A list of the IDs of ingress rules associated with the network ACL.

* `outbound_rules` - (Optional, List) A list of the IDs of egress rules associated with the network ACL.

* `subnets` - (Optional, List) A list of the IDs of networks associated with the network ACL.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the network ACL.
* `inbound_policy_id` - The ID of the ingress firewall policy for the network ACL.
* `outbound_policy_id` - The ID of the egress firewall policy for the network ACL.
* `ports` - A list of the port IDs of the subnet gateway.
* `status` - The status of the network ACL.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.
