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

## Note

ECS snapshots can only be created, but cannot be updated or deleted locally.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 14 hour.
