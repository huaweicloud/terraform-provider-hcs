---
subcategory: "Elastic IP (EIP)"
---

# hcs_vpc_bandwidth

Provides details about a specific bandwidth.

## Example Usage

```hcl
variable "bandwidth_name" {}

data "hcs_vpc_bandwidth" "bandwidth_1" {
  name = var.bandwidth_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available bandwidth in the current tenant. The
following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the bandwidth. If omitted, the provider-level region will
  be used.

* `name` - (Required, String) The name of the Shared Bandwidth to retrieve.

* `size` - (Optional, Int) The size of the Shared Bandwidth to retrieve. The value ranges from 5 Mbit/s to 2000 Mbit/s.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the Shared Bandwidth to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Shared Bandwidth.

* `share_type` - Indicates whether the bandwidth is shared or dedicated.

* `status` - Indicates the bandwidth status.

* `publicips` - An array of EIPs that use the bandwidth. The object includes the following:
  + `id` - The ID of the EIP or IPv6 port that uses the bandwidth.
  + `type` - The EIP type. Possible values are *5_bgp* (dynamic BGP) and *5_sbgp* (static BGP)
  + `ip_address` - The IPv4 or IPv6 address.
