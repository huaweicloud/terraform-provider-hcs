---
subcategory: "Elastic IP (EIP)"
---

# hcs_vpc_bandwidth

Manages a **Shared** Bandwidth resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_vpc_bandwidth" "bandwidth_1" {
  name = "bandwidth_1"
  size = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the Shared Bandwidth.
  If omitted, the provider-level region will be used. Changing this creates a new bandwidth.

* `name` - (Required, String) Specifies the bandwidth name. The value is a string of 1 to 64 characters that
  can contain letters, digits, underscores (_), hyphens (-), and periods (.).

* `size` - (Required, Int) Specifies the size of the Shared Bandwidth. The value ranges from 5 Mbit/s to 2000 Mbit/s.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the Shared Bandwidth.
  Changing this creates a new bandwidth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Shared Bandwidth.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `bandwidth_type` - Indicates the bandwidth type.

* `status` - Indicates the bandwidth status.

* `publicips` - An array of EIPs that use the bandwidth. The object includes the following:
  + `id` - The ID of the EIP or IPv6 port that uses the bandwidth.
  + `type` - The EIP type. Possible values are *5_bgp* (dynamic BGP) and *5_sbgp* (static BGP).
  + `ip_address` - The IPv4 or IPv6 address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Shared Bandwidths can be imported using the `id`, e.g.

```
$ terraform import hcs_vpc_bandwidth.bandwidth_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
