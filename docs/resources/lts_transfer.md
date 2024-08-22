---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_lts_transfer"
description: ""
---

# hcs_lts_transfer

Manages an LTS transfer task resource within HuaweiCloudStack.  

## Example Usage

### Create an OBS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "obs_bucket" {}

resource "hcs_lts_transfer" "test" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = var.obs_bucket
      obs_dir_prefix_name = "dir_prefix_"
      obs_prefix_name     = "prefix_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}
```

### Create a DIS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "dis_stream_id" {}
variable "dis_stream_name" {}

resource "hcs_lts_transfer" "test" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "DIS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      dis_id   = var.dis_stream_id
      dis_name = var.dis_stream_name
    }
  }
}
```

### Create a DMS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "kafka_instance_id" {}
variable "kafka_topic" {}

resource "hcs_lts_transfer" "test" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "DMS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      kafka_id    = var.kafka_instance_id
      kafka_topic = var.kafka_topic
    }
  }
}
```

### Create a delegated OBS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "obs_bucket" {}
variable "agency_domain_id" {}
variable "agency_domain_name" {}
variable "agency_name" {}
variable "agency_project_id" {}

resource "hcs_lts_transfer" "obs_agency" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = var.obs_bucket
      obs_dir_prefix_name = "dir_prefix_"
      obs_prefix_name     = "prefix_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }

    log_agency_transfer {
      agency_domain_id   = var.agency_domain_id
      agency_domain_name = var.agency_domain_name
      agency_name        = var.agency_name
      agency_project_id  = var.agency_project_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `log_group_id` - (Required, String, ForceNew) Log group ID.  

  Changing this parameter will create a new resource.

* `log_streams` - (Required, List, ForceNew) The list of log streams.  

  Changing this parameter will create a new resource.
  The [log_streams](#LtsTransfer_LogStreams) structure is documented below.

* `log_transfer_info` - (Required, List) Log transfer information.
  The [log_transfer_info](#LtsTransfer_LogTransferInfo) structure is documented below.

<a name="LtsTransfer_LogStreams"></a>
The `log_streams` block supports:

* `log_stream_id` - (Required, String, ForceNew) Log stream ID.

  Changing this parameter will create a new resource.

* `log_stream_name` - (Optional, String, ForceNew) Log stream name.

  Changing this parameter will create a new resource.

<a name="LtsTransfer_LogTransferInfo"></a>
The `log_transfer_info` block supports:

* `log_transfer_type` - (Required, String, ForceNew) Log transfer type.  
  The valid values is **OBS**.

  Changing this parameter will create a new resource.

* `log_transfer_mode` - (Required, String, ForceNew) Log transfer mode.  
  Value option is as follows:
    + **cycle**: Periodical transfer, which is available to OBS transfer tasks.

  Changing this parameter will create a new resource.

* `log_storage_format` - (Required, String) Log transfer format.  
  Value options are as follows:
    + **JSON**: JSON format, which is available to OBS transfer tasks.
    + **RAW**: Raw log format, which is available to OBS transfer tasks.

* `log_transfer_status` - (Required, String) Log transfer status.  
  Value options are as follows:
    + **ENABLE**: Log transfer is enabled.
    + **DISABLE**: Log transfer is disabled

* `log_agency_transfer` - (Optional, List, ForceNew) Information about agency which lets an account delegate resource management
   to other accounts.
  This parameter is mandatory if you transfer logs for another account.
  The [log_agency_transfer](#LtsTransfer_LogAgency) structure is documented below.

  Changing this parameter will create a new resource.

* `log_transfer_detail` - (Required, List) Log transfer details.  
  The [log_transfer_detail](#LtsTransfer_LogTransferDetail) structure is documented below.

<a name="LtsTransfer_LogAgency"></a>
The `log_agency_transfer` block supports:

* `agency_domain_id` - (Required, String, ForceNew) Delegator account ID.

  Changing this parameter will create a new resource.

* `agency_domain_name` - (Required, String, ForceNew) Delegator account name.

  Changing this parameter will create a new resource.

* `agency_name` - (Required, String, ForceNew) The agency name created by the delegator.

  Changing this parameter will create a new resource.

* `agency_project_id` - (Required, String, ForceNew) Project ID of the delegator.

  Changing this parameter will create a new resource.

<a name="LtsTransfer_LogTransferDetail"></a>
The `log_transfer_detail` block supports:

* `obs_period` - (Optional, Int) Length of the transfer interval for an OBS transfer task.  
  This parameter is mandatory when you create an OBS transfer task.  
  The log transfer interval is specified by the combination of the values of **obs_period** and **obs_period_unit**,
  and must be set to one of the following: 2 min, 5 min, 30 min, 1 hour, 3 hours, 6 hours, and 12 hours.
  Value options are as follows:
    + **2**: 2 minutes, the **obs_period_unit** must be **min**.
    + **5**: 5 minutes, the **obs_period_unit** must be **min**.
    + **30**: 30 minutes, the **obs_period_unit** must be **min**.
    + **1**: 1 hour, the **obs_period_unit** must be **hour**.
    + **3**: 3 hours, the **obs_period_unit** must be **hour**.
    + **6**: 6 hours, the **obs_period_unit** must be **hour**.
    + **12**: 12 hours, the **obs_period_unit** must be **hour**.

* `obs_period_unit` - (Optional, String) Unit of the transfer interval for an OBS transfer task.  
  This parameter is mandatory when you create an OBS transfer task.
  The log transfer interval is specified by the combination of the values of **obs_period** and **obs_period_unit**,
  and must be set to one of the following: 2 min, 5 min, 30 min, 1 hour, 3 hours, 6 hours, and 12 hours.
  Value options are as follows:
    + **min**: minute.
    + **hour**: hour.

* `obs_bucket_name` - (Optional, String) OBS bucket name.  
  This parameter is mandatory when you create an OBS transfer task.

* `obs_transfer_path` - (Optional, String) OBS bucket path, which is the log transfer destination.  

* `obs_dir_prefix_name` - (Optional, String) Custom transfer path of an OBS transfer task.  

* `obs_prefix_name` - (Optional, String) Transfer file prefix of an OBS transfer task.  

* `obs_eps_id` - (Optional, String) Enterprise project ID of an OBS transfer task.  

* `obs_encrypted_enable` - (Optional, Bool) Whether OBS bucket encryption is enabled.  

* `obs_encrypted_id` - (Optional, String) KMS key ID for an OBS transfer task.  
  This parameter is mandatory if encryption is enabled for the target OBS bucket.  

* `obs_time_zone` - (Optional, String) Time zone for an OBS transfer task.  
  If this parameter is specified, **obs_time_zone_id** must also be specified.

* `obs_time_zone_id` - (Optional, String) ID of the time zone for an OBS transfer task.  
  If this parameter is specified, **obs_time_zone** must also be specified.

* `dis_id` - (Optional, String) DIS stream ID.  
  This parameter is mandatory when you create a DIS transfer task.

* `dis_name` - (Optional, String) DIS stream name.  
  This parameter is mandatory when you create a DIS transfer task.

* `kafka_id` - (Optional, String) Kafka ID.  
  This parameter is mandatory when you create a DMS transfer task.

* `kafka_topic` - (Optional, String) Kafka topic.  
  This parameter is mandatory when you create a DMS transfer task.

  -> Before creating a DMS transfer task, register your Kafka instance with Kafka ID and Kafka topic first.

* `delivery_tags` - (Optional, List) The list of tag fields will be delivered when transferring.
  This field must contain the following host information: **regionName**, **projectId**, **logStreamName** and 
  **logGroupName**.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `log_group_name` - Log group name.  

## Import

The LTS transfer task can be imported using the `id`, e.g.

```bash
$ terraform import hcs_lts_transfer.test 0ce123456a00f2591fabc00385ff1234
```
