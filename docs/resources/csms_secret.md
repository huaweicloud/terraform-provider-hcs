---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_csms_secret"
description: |
  Manages CSMS(Cloud Secret Management Service) secrets within HuaweiCloudStack.
---

# hcs_csms_secret

Manages CSMS(Cloud Secret Management Service) secrets within HuaweiCloudStack.

## Example Usage

### Encrypt Plaintext

```hcl
resource "hcs_csms_secret" "test1" {
  name        = "test_secret"
  secret_text = "this is a password"
}
```

### Encrypt JSON Data

```hcl
resource "hcs_csms_secret" "test2" {
  name        = "mysql_admin"
  secret_text = jsonencode({
    username = "admin"
    password = "123456"
  })
}
```

### Encrypt String Binary

```hcl
variable "secret_binary" {}

resource "hcs_csms_secret" "test3" {
  name          = "test_binary"
  secret_binary = var.secret_binary
}
```

### The secret associated event

```hcl
variable "name" {}
variable "secret_type" {}
variable "secret_text" {}

resource "huaweicloud_csms_event" "test" {
  ...
}

resource "hcs_csms_secret" "test" {
  name                = var.name
  secret_type         = var.secret_type
  secret_text         = var.secret_text
  event_subscriptions = [huaweicloud_csms_event.test.name]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CSMS secrets.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the secret name. The maximum length is 64 characters.
  Only digits, letters, underscores(_), hyphens(-) and dots(.) are allowed.

  Changing this parameter will create a new resource.

* `secret_text` - (Optional, String) Specifies the plaintext of a text secret. CSMS encrypts the plaintext and stores
  it in the initial version of the secret. The maximum size is 32 KB.

  Changing this parameter will create a new secret version.

* `kms_key_id` - (Optional, String) Specifies the ID of the KMS key used to encrypt secrets.
  If this parameter is not specified when creating the secret, the default master key **csms/default** will be used.
  The default key is automatically created by the CSMS.
  Use this data source
  [hcs_kms_key](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs/resources/kms_key)
  to get the KMS key.

* `description` - (Optional, String) Specifies the description of a secret.

* `tags` - (Optional, Map) Specifies the tags of a CSMS secrets, key/value pair format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is constructed from the secret ID and name, separated by a slash.

* `secret_id` - The secret ID in UUID format.

* `latest_version` - The latest version id.

* `status` - The CSMS secret status. Values can be: **ENABLED**, **DISABLED**, **PENDING_DELETE** and **FROZEN**.

* `create_time` - Time when the CSMS secrets created, in UTC format.

## Import

CSMS secret can be imported using the ID and the name of secret, separated by a slash, e.g.

```bash
terraform import hcs_csms_secret.test <id>/<name>
```
