---
subcategory: "NAT Gateway (NAT)"
---

# hcs_nat_dnat_rule

Manages a DNAT rule resource of the **public** NAT within HuaweiCloudStack(hcs).

## Example Usage

### DNAT rule in VPC scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "hcs_compute_instance" "test" {
  ...
}

resource "hcs_nat_dnat_rule" "test" {
  nat_gateway_id        = var.gateway_id
  floating_ip_id        = var.publicip_id
  port_id               = hcs_compute_instance.test.network[0].port
  protocol              = "tcp"
  internal_service_port = 23
  external_service_port = 8023
}
```

### DNAT rule in VPC scenario and specify the port ranges

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "hcs_compute_instance" "test" {
  ...
}

resource "hcs_nat_dnat_rule" "test" {
  nat_gateway_id              = var.gateway_id
  floating_ip_id              = var.publicip_id
  port_id                     = hcs_compute_instance.test.network[0].port
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String, ForceNew) Specifies the ID of the NAT gateway to which the DNAT rule belongs.  
  Changing this will create a new resource.

* `floating_ip_id` - (Required, String) Specifies the ID of the floating IP address.

* `protocol` - (Required, String) Specifies the protocol type.  
  The valid values are **tcp**, **udp**.

* `internal_service_port` - (Optional, Int) Specifies port used by Floating IP provide services for external
  systems.  
  Exactly one of `internal_service_port` and `internal_service_port_range` must be set.

* `external_service_port` - (Optional, Int) Specifies port used by ECSs or BMSs to provide services for
  external systems.  
  Exactly one of `external_service_port` and `external_service_port_range` must be set.  
  Required if `internal_service_port` is set.

* `internal_service_port_range` - (Optional, String) Specifies port range used by Floating IP provide services
  for external systems.  
  This parameter and `external_service_port_range` are mapped **1:1** in sequence(, ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.

* `external_service_port_range` - (Optional, String) Specifies port range used by ECSs or BMSs to provide
  services for external systems.  
  This parameter and `internal_service_port_range` are mapped **1:1** in sequence(ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.  
  Required if `internal_service_port_range` is set.

* `port_id` - (Optional, String) Specifies the port ID of network. This parameter is mandatory in VPC scenario.

* `description` - (Optional, String) Specifies the description of the DNAT rule.  
  The value is a string of no more than `255` characters, and angle brackets (<>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The creation time of the DNAT rule.

* `status` - The current status of the DNAT rule.

* `floating_ip_address` - The actual floating IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

DNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import hcs_nat_dnat_rule.test f4f783a7-b908-4215-b018-724960e5df4a
```
