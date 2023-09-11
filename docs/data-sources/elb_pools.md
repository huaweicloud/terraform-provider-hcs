---
subcategory: Dedicated Load Balance (Dedicated ELB)
---

# hcs_elb_pools

Use this data source to get the list of ELB pools.

## Example Usage

```hcl
variable "pool_name" {}

data "hcs_elb_pools" "test" {
  name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the ELB pool.

* `pool_id` - (Optional, String) Specifies the ID of the ELB pool.

* `description` - (Optional, String) Specifies the description of the ELB pool.

* `loadbalancer_id` - (Optional, String) Specifies the loadbalancer ID of the ELB pool.

* `healthmonitor_id` - (Optional, String) Specifies the health monitor ID of the ELB pool.

* `protocol` - (Optional, String) Specifies the protocol of the ELB pool. This can either be TCP, UDP or HTTP.

* `lb_method` - (Optional, String) Specifies the method of the ELB pool. Must be one of ROUND_ROBIN, LEAST_CONNECTIONS,
  or SOURCE_IP.

* `listener_id` - (Optional, String) Specifies the listener ID of the ELB pool.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `pools` - Pool list. For details, see data structure of the pool field.
  The [object](#pools_object) structure is documented below.

<a name="pools_object"></a>
The `pools` block supports:

* `id` - The pool ID.

* `name` - The pool name.

* `description` - The description of pool.

* `protocol` - The protocol of pool.

* `lb_method` - The load balancing algorithm to distribute traffic to the pool's members.

* `healthmonitor_id` - The health monitor ID of the LB pool.

* `ip_version` - The IP version of the LB pool.

* `listeners` - The listener list. The [object](#elem_object) structure is documented below.

* `loadbalancers` - The loadbalancer list. The [object](#elem_object) structure is documented below.

* `members` - The member list. The [object](#elem_object) structure is documented below.

* `persistence` - Indicates whether connections in the same session will be processed by the same pool member or not.
  The [object](#persistence_object) structure is documented below.

<a name="elem_object"></a>
The `listeners`,  `loadbalancers` or `members` block supports:

* `id` - The listener, loadbalancer or member ID.

<a name="persistence_object"></a>
The `persistence` block supports:

* `type` - The type of persistence mode.

* `cookie_name` - The name of the cookie if persistence mode is set appropriately.
