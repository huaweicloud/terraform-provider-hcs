---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_snapshot_rollback

Rollback an ECS Snapshot for an Instance.

## Example Usage

### Rollback an ECS Snapshot

```hcl
resource "hcs_ecs_compute_snapshot_rollback" "snapshot" {
  instance_id = "6e6da0c2-6ade-41ce-bd31-62fd222ec115"
  snapshot_id = "5c1892f3-87e7-4ca0-a111-f8093eed2bcc"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to rollback ECS snapshot.

* `snapshot_id` - (Required, String, ForceNew) The ID of the ECS snapshot.

## Note

ECS snapshots can only be rolled back. Local update and deletion are not supported.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
