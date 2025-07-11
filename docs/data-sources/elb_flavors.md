---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# hcs_elb_flavors

Use this data source to get the available ELB Flavors(QoS).

## Example Usage

```hcl
variable "flavor_name" {}

data "hcs_elb_flavors" "test" {
  name = var.flavor_name
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be
  used.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `name` - (Optional, String) Specifies the flavor name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ids` - Indicates the list of flavor IDs.

* `flavors` - Indicates the list of flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the flavor.

* `name` - Indicates the name of the flavor.

* `type` - Indicates the type of the flavor.

* `shared` - Indicates whether the flavor is available to all users(*Not supported currently*).

* `flavor_sold_out` - Indicates whether the flavor has ten or more load balancers bound.

* `status` - Indicates the binding status of the flavor.

* `info` - A map of flavor specifications, which may include:

   + **connection** - Indicates the maximum connections of the flavor.
   + **cps** - Indicates the cps of the flavor.
   + **bandwidth** - Indicates the bandwidth size (Mbit/s) of the flavor.
   + **qps** - Indicates the qps of the L7 flavor.

