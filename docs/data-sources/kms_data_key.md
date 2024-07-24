---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_kms_data_key"
description: ""
---

# hcs_kms_data_key

Use this data source to get the plaintext and the ciphertext of an available HuaweiCloudStack KMS DEK (data encryption key).

## Example Usage

```hcl
resource "hcs_kms_key" "key1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
}

data "hcs_kms_data_key" "kms_datakey1" {
  key_id         = hcs_kms_key.key1.id
  datakey_length = "512"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the keys. If omitted, the provider-level region will be
  used.

* `key_id` - (Required, String) The globally unique identifier for the key. Changing this gets the new data encryption
  key.

* `encryption_context` - (Optional, String) The value of this parameter must be a series of
  "key:value" pairs used to record resource context information. The value of this parameter must not contain sensitive
  information and must be within 8192 characters in length. Example: {"Key1":"Value1","Key2":"Value2"}

* `datakey_length` - (Required, String) Number of bits in the length of a DEK (data encryption keys). The maximum number
  is 512. Changing this gets the new data encryption key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `plain_text` - The plaintext of a DEK is expressed in hexadecimal format, and two characters indicate one byte.

* `cipher_text` - The ciphertext of a DEK is expressed in hexadecimal format, and two characters indicate one byte.
