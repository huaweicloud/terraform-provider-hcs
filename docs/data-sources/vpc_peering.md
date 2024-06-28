---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_peering

The VPC Peering data source provides details about a specific VPC peering.

## Example Usage

```hcl
data "hcs_vpc_peering" "peering" {
  name = "peering"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPC peering. The given filters
must match exactly one VPC peering whose data will be exported as attributes.

* `id` - (Optional, String) The ID of the specific VPC Peering to retrieve.

* `name` - (Optional, String) The name of the specific VPC Peering to retrieve.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpc_id` - The VPC ID of the VPC Peering.

* `peer_vpc_id` - The peer VPC ID of the VPC Peering.
