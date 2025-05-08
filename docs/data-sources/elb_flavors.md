---
subcategory: Dedicated Load Balance (Dedicated ELB)
---

# hcs_elb_flavors

Use this data source to get the available ELB Flavors(Qos).

## Example Usage

```hcl
data "hcs_elb_flavors" "flavors" {
  type            = "l7"
  max_connections = 20000
  cps             = 1000
  bandwidth       = 20
  qps             = 80
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be
  used.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `type` - (Optional, String) Specifies the flavor type. Values options:
    + **l4**: indicates Layer-4 flavor.
    + **l7**: indicates Layer-7 flavor.

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

* `flavor_sold_out` - Indicates whether the flavor has ten or more load balancers bound.
  + **true**: The flavor has ten or more load balancers bound.
  + **false**: The flavor has fewer than ten load balancers bound.

* `status` - Specifies the binding status of the flavor. Allowed values are:
  + **bind** – Indicates the flavor is currently bound.
  + **unbind** – Indicates the flavor is currently unbound.

* `max_connections` - Indicates the maximum connections of the flavor.

* `cps` - Indicates the cps of the flavor.

* `qps` - Indicates the qps of the L7 flavor.

* `bandwidth` - Indicates the bandwidth size(Mbit/s) of the flavor.