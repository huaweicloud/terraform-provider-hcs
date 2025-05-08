---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# hcs_elb_flavor

Manages an ELB Flavor(Qos) resource within HCS.

## Example Usage

```hcl
resource "hcs_elb_flavor" "flavor_1" {
    name = "flavor1"
    type = "l4"
    info {
      flavor_type = "cps"
      value       = 5000
    }
    info {
      flavor_type = "connection"
      value       = 100
    }
    info {
      flavor_type = "bandwidth"
      value       = 20
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB flavor resource. If omitted, the
  provider-level region will be used. Changing this creates a new flavor.

* `type` - (Required, String, ForceNew) Specifies the flavor type. Values options:
-- l4: indicates Layer-4 flavor.
-- l7: indicates Layer-7 flavor.

* `name` - (Optional, String) Human-readable name for the Flavor. Does not have to be unique.

* `info` - (Required, Set) Specifies the details of each flavor metric. Multiple info blocks can be provided.
  The [object](#info) structure is documented below.

<a name="info"></a>
The `info` block supports:
* `flavor_type` - (Required, String) Specifies the type of the flavor metric.  
  Value options depend on the parent resource's `type` argument:
  + If `type` is **l4**, allowed values are **cps**, **connection**, and **bandwidth**.
  + If `type` is **l7**, allowed values are **cps**, **connection**, **bandwidth**, and **qps**.

* `value` - (Required, Int) Specifies the value for the specified `flavor_type`.  
  The value must be in the range `[0, cluster maximum]`.  
  If the value exceeds the maximum permitted for a single cluster, it will be automatically set to the cluster maximum.  
  The actual supported maximum value is proportional to the number of network elements (NEs) in the cluster. The per-NE single maximums are:
  + If `type` is **l4**:
    - **cps:** `87,000`
    - **connection:** `2,200,000`
    - **bandwidth:** `16,384` (Mbit/s)
  + If `type` is **l7**:
    - **cps:** `20,000`
    - **connection:** `320,000`
    - **bandwidth:** `10,384` (Mbit/s)
    - **qps:** `300,000`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the flavor.

* `flavor_sold_out` - Indicates whether the flavor has ten or more loadbalancers bound.
  + **true**: The flavor has ten or more loadbalancers bound.
  + **false**: The flavor has fewer than ten loadbalancers bound.

* `status` - Specifies the binding status of the flavor. Allowed values are:
  + **bind** – Indicates the flavor is currently bound.
  + **unbind** – Indicates the flavor is currently unbound.

## Import

ELB flavor can be imported using the flavor ID, e.g.

```
$ terraform import hcs_elb_flavor.flavor_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
