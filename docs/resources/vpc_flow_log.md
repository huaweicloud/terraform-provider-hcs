---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_flow_log

Manages a VPC flow log resource.

## Example Usage

```hcl
var "flowlog_name" {}
var "vpc_id" {}
var "log_group_id" {}
var "log_stream_id" {}

resource "hcs_vpc_flow_log" "test_flowlog" {
  name          = var.flowlog_name
  resource_type = "vpc"
  resource_id   = var.vpc_id
  traffic_type  = "all"
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new VPC flow log.

* `name` - (Required, String) Specifies the VPC flow log name. The value can contain no more than 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `resource_type` - (Required, String, ForceNew) Specifies the resource type for which that the logs to be collected.
  The value can be:
  + *port*: Select this to record traffic information about a specific NIC.
  + *network*: Select this to record traffic information about all NICs in a specific subnet.
  + *vpc*: Select this to record traffic information about all NICs in a specific VPC.

  Changing this creates a new VPC flow log.

* `resource_id` - (Required, String, ForceNew) Specifies the resource ID for which that the logs to be collected.
  Changing this creates a new VPC flow log.

* `log_group_id` - (Required, String, ForceNew) Specifies the LTS log group ID.
  Changing this creates a new VPC flow log.

* `log_stream_id` - (Required, String, ForceNew) Specifies the LTS log stream ID.
  Changing this creates a new VPC flow log.

* `traffic_type` - (Optional, String, ForceNew) Specifies the type of traffic to log. The value can be:
  + *all*: Specifies that both accepted and rejected traffic of the specified resource will be logged.
  + *accept*: Specifies that only accepted inbound and outbound traffic of the specified resource will be logged.
  + *reject*: Specifies that only rejected inbound and outbound traffic of the specified resource will be logged.

  Defaults to *all*. Changing this creates a new VPC flow log.

* `enabled` - (Optional, Bool) Specifies whether to enable the flow log function, the default value is *true*.

* `description` - (Optional, String) Specifies description about the VPC flow log.
  The value can contain no more than 255 characters and cannot contain angle brackets (< or >).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC flow log ID in UUID format.

* `status` - The status of the flow log. The value can be `ACTIVE`, `DOWN` or `ERROR`.

## Import

VPC flow logs can be imported using the `id`, e.g.

```bash
$ terraform import hcs_vpc_flow_log.flowlog1 41b9d73f-eb1c-4795-a100-59a99b062513
```
