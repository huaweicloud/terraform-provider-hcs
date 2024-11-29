---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_snapshot

Creating an ECS Snapshot for an Instance.

## Example Usage

### Creating an ECS Snapshot

```hcl
resource "hcs_ecs_compute_snapshot" "snapshot" {
  instance_id = "6e6da0c2-6ade-41ce-bd31-62fd222ec115"
  name        = "ecs_snapshot_02"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to create ECS snapshot.

* `name` - (Required, String, ForceNew) The snapshot name.

## Import

Snapshot can be imported using the Instance ID and Snapshot ID separated by a slash, e.g.
```
$ terraform import hcs_ecs_compute_snapshot.test 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf/aa4a8f8d-160d-4643-9ab4-087657b307ba
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 14 hour.
