---
subcategory: "Elastic IP (EIP)"
---

# hcs_vpc_eips

Use this data source to get a list of EIPs.

## Example Usage

An example filter by name and tag

```hcl
variable "enterprise_project_id" {}

data "hcs_vpc_eips" "eip" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available EIPs in the current region.
 All EIPs that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the EIP. If omitted, the provider-level region
  will be used.

* `ids` - (Optional, List) Specifies an array of one or more IDs of the desired EIP.

* `public_ips` - (Optional, List) Specifies an array of one or more public ip addresses of the desired EIP.

* `port_ids` - (Optional, List) Specifies an array of one or more port ids which bound to the desired EIP.

  The default value is `4`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired EIP belongs to.

 -> A maximum of 10 tag keys are allowed for each query operation. Each tag key can have up to 10 tag values.
  The tag key cannot be left blank or set to an empty string. Each tag key must be unique, and each tag value in a
  tag must be unique, use commas(,) to separate the multiple values. An empty for values indicates any value.
  The values are in the OR relationship.

## Attributes Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `eips` - Indicates a list of all EIPs found. Structure is documented below.

The `eips` block supports:

* `id` - The ID of the EIP.
* `name` - The name of the EIP.
* `public_ip` - The public ip address of the EIP.
* `private_ip` - The private ip address of the EIP.
* `port_id` - The port id bound to the EIP.
* `status` - The status of the EIP.
* `type` - The type of the EIP.
* `enterprise_project_id` - The the enterprise project ID of the EIP.
* `bandwidth_id` - The bandwidth id of the EIP.
* `bandwidth_name` - The bandwidth name of the EIP.
* `bandwidth_size` - The bandwidth size of the EIP.
