---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_server_group

Manages Server Group resource within HuaweiCloudStack.

## Example Usage

```hcl
data "hcs_ecs_compute_instance" "instance_demo" {
  name = "ecs-servergroup-demo"
}

resource "hcs_ecs_compute_server_group" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
  members  = [
    data.hcs_ecs_compute_instance.instance_demo.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the server group resource. If omitted,
  the provider-level region will be used. Changing this creates a new server group.

* `name` - (Required, String, ForceNew) Specifies a unique name for the server group. This parameter can contain a
  maximum of 255 characters, which may consist of letters, digits, underscores (_), and hyphens (-). Changing this
  creates a new server group.

* `policies` - (Required, List, ForceNew) Specifies the policy set for the server group. The value must belong 
  to *anti-affinity*, *affinity*, *soft-affinity*, and *soft-anti-affinity*. *anti-affinity* and *affinity* 
  cannot belong to the same group.
  + `anti-affinity`: The servers in this group must be scheduled to different hosts

  + `affinity`: The servers in this group must be scheduled on the same host.

  + `soft-affinity`: If possible, servers in this group are scheduled on the same host. However, 
  if this function cannot be implemented, they should still be scheduled instead of causing task failure.

  + `soft-anti-affinity`: If possible, servers in this group should be scheduled to different hosts. 
  If this function cannot be implemented, they should still be scheduled instead of causing task failure.

* `members` - (Optional, List) Specifies an array of one or more instance ID to attach server group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import hcs_ecs_compute_server_group.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
