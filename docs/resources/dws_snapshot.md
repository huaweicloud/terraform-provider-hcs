---
subcategory: "Data Warehouse Service (DWS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dws_snapshot"
description: |-
  Manages a DWS snapshot resource within HuaweiCloudStack.  
---

# hcs_dws_snapshot

Manages a DWS snapshot resource within HuaweiCloudStack.  

## Example Usage

```hcl
variable "cluster_id" {}

resource "hcs_dws_snapshot" "test" {
  name        = "demo"
  cluster_id  = var.cluster_id
  description = "This is a demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of snapshot.

  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the cluster for which you want to create a snapshot.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of snapshot.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `started_at` - The time when a snapshot starts to be created.  
  The format is **ISO8601**: **YYYY-MM-DDThh:mm:ssZ**.

* `finished_at` - The time when a snapshot is complete.  
  The format is **ISO8601**: **YYYY-MM-DDThh:mm:ssZ**.

* `size` - The snapshot size, in GB.

* `status` - The snapshot status.  
  The valid values are as follows:
  + **CREATING**
  + **AVAILABLE**
  + **UNAVAILABLE**

* `type` - The snapshot type.  
  The valid values are as follows:
  + **MANUAL**
  + **AUTOMATED**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The dws snapshot can be imported using the `id`, e.g.

```bash
$ terraform import hcs_dws_snapshot.test e87192d9-b592-4658-b23f-bdc0bb69ec2c
```
