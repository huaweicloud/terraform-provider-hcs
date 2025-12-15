---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_kms_key"
description: |-
  Use this data source to get the ID of an available HuaweiCloudStack KMS key.
---

# hcs_kms_key

-> **NOTE:** This data source can only be used in HCS **8.5.0** and **later** version.

Use this data source to get the ID of an available HuaweiCloudStack KMS key.

## Example Usage

```hcl
data "hcs_kms_key" "test" {
  key_alias        = "test_key"
  key_description  = "test key description"
  key_state        = "2"
  key_id           = "af650527-a0ff-4527-aef3-c493df1f3012"
  default_key_flag = "0"
  domain_id        = "b168fe00ff56492495a7d22974df2d0b"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the keys. If omitted, the provider-level region will be
  used.

* `key_alias` - (Optional, String) The alias in which to create the key. It is required when we create a new key.
  Changing this gets the new key.

* `key_description` - (Optional, String) The description of the key as viewed in Huawei console. Changing this gets a
  new key.

* `key_id` - (Optional, String) The globally unique identifier for the key. Changing this gets the new key.

* `default_key_flag` - (Optional, String) Identification of a Master Key. The valid values are as follows:
  + **1**: Default master key
  + **0**: A key

* `key_state` - (Optional, String) The state of a key. The valid values are as follows:
  + **1**: Pending activation
  + **2**: Enabled
  + **3**: Disabled
  + **4**: Scheduled for deletion
  + **5**: Pending import

* `domain_id` - (Optional, String) The ID of a user domain for the key. Changing this gets a new key.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the kms key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `scheduled_deletion_date` - The scheduled deletion time (time stamp) of a key.

* `expiration_time` - The expiration time.

* `creation_date` - The creation time (time stamp) of a key.

* `tags` - The key/value pairs to associate with the kms key.

* `rotation_enabled` - Whether the key rotation is enabled or not.

* `rotation_interval` - The key rotation interval. It's valid when rotation is enabled.

* `rotation_number` - The total number of key rotations. It's valid when rotation is enabled.
