---
subcategory: "NAT Gateway (NAT)"
---

# hcs_nat_snat_rule

Manages an SNAT rule resource of the **public** NAT within HuaweiCloudStack(hcs).

## Example Usage

### SNAT rule in VPC scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}
variable "subent_id" {}

resource "hcs_nat_snat_rule" "test" {
  nat_gateway_id = var.gateway_id
  floating_ip_id = var.publicip_id
  subnet_id      = var.subent_id
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String, ForceNew) Specifies the ID of the gateway to which the SNAT rule belongs.  
  Changing this will create a new resource.

* `floating_ip_id` - (Required, String, ForceNew) Specifies the IDs of floating IPs connected by SNAT rule.  
  Multiple floating IPs are separated using commas (,). The number of floating IP IDs cannot exceed `20`.
  Changing this will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies the network IDs of subnet connected by SNAT rule (VPC side).  
  Changing this will create a new resource.

* `source_type` - (Optional, Int) Specifies the resource scenario.  
  The valid value is **0** (VPC scenario), and the default value is `0`.
  Changing this will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the SNAT rule.
  The value is a string of no more than `255` characters, and angle brackets (<>) are not allowed.
  Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `floating_ip_address` - The actual floating IP address.

* `status` - The status of the SNAT rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

SNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import hcs_nat_snat_rule.test 9e0713cb-0a2f-484e-8c7d-daecbb61dbe4
```
