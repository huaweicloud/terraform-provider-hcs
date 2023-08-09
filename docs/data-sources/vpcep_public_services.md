---
subcategory: "VPC Endpoint (VPCEP)"
---

# hcs_vpcep_public_services

Use this data source to get available public VPC endpoint services.

## Example Usage

```hcl
data "hcs_vpcep_public_services" "all_services" {
}

```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the public VPC endpoint services. If omitted, the
  provider-level region will be used.

* `service_name` - (Optional, String) Specifies the name of the public VPC endpoint service. The value is not
  case-sensitive and supports fuzzy match.

* `service_id` - (Optional, String) Specifies the unique ID of the public VPC endpoint service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `services` - Indicates the public VPC endpoint services information. Structure is documented below.

The `services` block contains:

* `id` - The unique ID of the public VPC endpoint service.
* `service_name` - The name of the public VPC endpoint service.
* `service_type` - The type of the VPC endpoint service.
